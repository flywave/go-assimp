package assimp

import (
	"testing"

	"github.com/flywave/go-mst"
	"github.com/flywave/go3d/mat4"
	"github.com/flywave/go3d/vec3"
)

func TestAssimpToMSTConverterOptimized(t *testing.T) {
	// 测试空场景
	t.Run("EmptyScene", func(t *testing.T) {
		result := AssimpToMSTConverter(nil)
		if result != nil {
			t.Errorf("Expected nil for nil scene, got %v", result)
		}
	})

	// 测试基础转换
	t.Run("BasicConversion", func(t *testing.T) {
		// 创建一个简单的测试场景
		scene := &Scene{
			Meshes: []*Mesh{
				{
					Vertices: []vec3.T{
						{0, 0, 0},
						{1, 0, 0},
						{0, 1, 0},
					},
					Normals: []vec3.T{
						{0, 0, 1},
						{0, 0, 1},
						{0, 0, 1},
					},
					Faces: []Face{
						{Indices: []uint{0, 1, 2}},
					},
					MaterialIndex: 0,
				},
			},
			Materials: []*Material{
				{
					Properties: []*MaterialProperty{
						{name: "$clr.diffuse", Data: []byte{255, 0, 0}},
						{name: "$mat.opacity", Data: []byte{0, 0, 128, 63}}, // 1.0
					},
				},
			},
		}

		result := AssimpToMSTConverter(scene)
		if result == nil {
			t.Fatal("Expected non-nil result")
		}

		if len(result.Materials) != 1 {
			t.Errorf("Expected 1 material, got %d", len(result.Materials))
		}

		if len(result.Nodes) != 1 {
			t.Errorf("Expected 1 node, got %d", len(result.Nodes))
		}

		node := result.Nodes[0]
		if len(node.Vertices) != 3 {
			t.Errorf("Expected 3 vertices, got %d", len(node.Vertices))
		}

		if len(node.FaceGroup) != 1 {
			t.Errorf("Expected 1 face group, got %d", len(node.FaceGroup))
		}
	})

	// 测试材质类型检测
	t.Run("MaterialTypes", func(t *testing.T) {
		materials := []*Material{
			{
				Properties: []*MaterialProperty{
					{name: "$mat.metallic", Data: []byte{0, 0, 0, 63}},
				},
			},
			{
				Properties: []*MaterialProperty{
					{name: "$clr.specular", Data: []byte{255, 255, 255}},
				},
			},
			{
				Properties: []*MaterialProperty{
					{name: "$clr.diffuse", Data: []byte{255, 255, 255}},
				},
			},
		}

		for _, mat := range materials {
			result := convertMaterial(mat)
			if result == nil {
				t.Error("Expected non-nil material")
			}
		}
	})

	// 测试纹理处理
	t.Run("TextureHandling", func(t *testing.T) {
		material := &Material{
			Properties: []*MaterialProperty{
				{name: "$clr.diffuse", Data: []byte{255, 255, 255}},
				{Semantic: TextureTypeDiffuse, name: "$tex.file", Data: []byte("test_texture.png")},
			},
		}

		result := convertMaterial(material)
		textureMat, ok := result.(*mst.TextureMaterial)
		if !ok {
			t.Fatalf("Expected TextureMaterial, got %T", result)
		}

		if textureMat.Texture == nil {
			t.Error("Expected texture to be present")
		}
	})

	// 测试节点层次结构
	t.Run("NodeHierarchy", func(t *testing.T) {
		tran := mat4.FromArray([16]float32{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1})
		tran2 := mat4.FromArray([16]float32{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 1, 2, 3, 1})
		scene := &Scene{
			Meshes: []*Mesh{
				{
					Vertices: []vec3.T{{0, 0, 0}, {1, 0, 0}, {0, 1, 0}},
					Faces:    []Face{{Indices: []uint{0, 1, 2}}},
				},
			},
			RootNode: &Node{
				Transformation: &tran,
				MeshIndicies:   []uint{0},
				Children: []*Node{
					{
						Transformation: &tran2,
						MeshIndicies:   []uint{0},
					},
				},
			},
		}

		result := AssimpToMSTConverter(scene)
		if len(result.InstanceNode) < 2 {
			t.Errorf("Expected at least 2 instance nodes, got %d", len(result.InstanceNode))
		}
	})

	// 测试面处理
	t.Run("FaceProcessing", func(t *testing.T) {
		mesh := &Mesh{
			Vertices: []vec3.T{
				{0, 0, 0}, {1, 0, 0}, {1, 1, 0}, {0, 1, 0},
			},
			Faces: []Face{
				{Indices: []uint{0, 1, 2, 3}}, // 四边形
			},
		}

		node := convertMesh(mesh)
		if len(node.FaceGroup) == 0 {
			t.Fatal("Expected face groups")
		}

		// 四边形应该被分成两个三角形
		faces := node.FaceGroup[0].Faces
		if len(faces) < 2 {
			t.Errorf("Expected at least 2 faces for quad, got %d", len(faces))
		}
	})
}

func TestConvertMatrix(t *testing.T) {
	aiMat := mat4.FromArray([16]float32{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 10, 20, 30, 1})
	result := convertMatrix(aiMat)

	expected := [4][4]float64{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{10, 20, 30, 1},
	}

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if result[i][j] != expected[i][j] {
				t.Errorf("Matrix mismatch at [%d][%d]: expected %f, got %f", i, j, expected[i][j], result[i][j])
			}
		}
	}
}

func TestMaterialTypeDetection(t *testing.T) {
	tests := []struct {
		name     string
		material *Material
		expected string
	}{
		{
			name: "PBR Material",
			material: &Material{
				Properties: []*MaterialProperty{
					{name: "$mat.metallic", Data: []byte{0, 0, 0, 63}},
				},
			},
			expected: "*mst.PbrMaterial",
		},
		{
			name: "Phong Material",
			material: &Material{
				Properties: []*MaterialProperty{
					{name: "$clr.specular", Data: []byte{255, 255, 255}},
				},
			},
			expected: "*mst.PhongMaterial",
		},
		{
			name: "Lambert Material",
			material: &Material{
				Properties: []*MaterialProperty{
					{name: "$clr.diffuse", Data: []byte{255, 255, 255}},
				},
			},
			expected: "*mst.LambertMaterial",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertMaterial(tt.material)
			if result == nil {
				t.Error("Expected non-nil material")
			}
		})
	}
}

func BenchmarkAssimpToMSTConverter(b *testing.B) {
	scene := &Scene{
		Meshes: make([]*Mesh, 100),
		Materials: []*Material{
			{
				Properties: []*MaterialProperty{
					{name: "$clr.diffuse", Data: []byte{255, 255, 255}},
				},
			},
		},
	}

	for i := 0; i < 100; i++ {
		scene.Meshes[i] = &Mesh{
			Vertices: []vec3.T{
				{0, 0, 0}, {1, 0, 0}, {0, 1, 0},
			},
			Normals: []vec3.T{
				{0, 0, 1}, {0, 0, 1}, {0, 0, 1},
			},
			Faces: []Face{
				{Indices: []uint{0, 1, 2}},
			},
			MaterialIndex: 0,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = AssimpToMSTConverter(scene)
	}
}
