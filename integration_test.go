package assimp

import (
	"testing"

	"github.com/flywave/go3d/vec3"
)

// TestIntegrationBasic 测试基本的集成流程
func TestIntegrationBasic(t *testing.T) {
	// 测试场景创建
	scene := &Scene{
		Flags: SceneFlagValidated,
		RootNode: &Node{
			Name: "root",
			Children: []*Node{
				{
					Name:         "mesh_node",
					MeshIndicies: []uint{0},
				},
			},
		},
		Meshes: []*Mesh{
			{
				Name:          "test_mesh",
				Vertices:      make([]vec3.T, 100),
				Normals:       make([]vec3.T, 100),
				Faces:         make([]Face, 50),
				MaterialIndex: 0,
				AABB: AABB{
					Min: vec3.T{-1, -1, -1},
					Max: vec3.T{1, 1, 1},
				},
			},
		},
		Materials: []*Material{
			{
				Properties: []*MaterialProperty{
					{
						name:     "diffuse_color",
						Semantic: TextureTypeDiffuse,
						Index:    0,
						TypeInfo: MatPropTypeInfoFloat32,
					},
				},
			},
		},
	}

	// 验证场景结构
	if scene.RootNode == nil {
		t.Fatal("Expected non-nil root node")
	}

	if len(scene.Meshes) != 1 {
		t.Errorf("Expected 1 mesh, got %d", len(scene.Meshes))
	}

	if len(scene.Materials) != 1 {
		t.Errorf("Expected 1 material, got %d", len(scene.Materials))
	}

	mesh := scene.Meshes[0]
	if len(mesh.Vertices) != 100 {
		t.Errorf("Expected 100 vertices, got %d", len(mesh.Vertices))
	}

	if len(mesh.Normals) != 100 {
		t.Errorf("Expected 100 normals, got %d", len(mesh.Normals))
	}

	if len(mesh.Faces) != 50 {
		t.Errorf("Expected 50 faces, got %d", len(mesh.Faces))
	}

	if mesh.Name != "test_mesh" {
		t.Errorf("Expected mesh name 'test_mesh', got '%s'", mesh.Name)
	}

	// 验证节点关系
	if scene.RootNode.Name != "root" {
		t.Errorf("Expected root node name 'root', got '%s'", scene.RootNode.Name)
	}

	if len(scene.RootNode.Children) != 1 {
		t.Errorf("Expected root to have 1 child, got %d", len(scene.RootNode.Children))
	}

	meshNode := scene.RootNode.Children[0]
	if meshNode.Name != "mesh_node" {
		t.Errorf("Expected child node name 'mesh_node', got '%s'", meshNode.Name)
	}

	if len(meshNode.MeshIndicies) != 1 {
		t.Errorf("Expected 1 mesh index, got %d", len(meshNode.MeshIndicies))
	}

	if meshNode.MeshIndicies[0] != 0 {
		t.Errorf("Expected mesh index 0, got %d", meshNode.MeshIndicies[0])
	}

	// 验证材质
	material := scene.Materials[0]
	if len(material.Properties) != 1 {
		t.Errorf("Expected 1 material property, got %d", len(material.Properties))
	}

	prop := material.Properties[0]
	if prop.name != "diffuse_color" {
		t.Errorf("Expected property name 'diffuse_color', got '%s'", prop.name)
	}

	if prop.Semantic != TextureTypeDiffuse {
		t.Errorf("Expected semantic TextureTypeDiffuse, got %d", prop.Semantic)
	}

	// 验证边界框
	if mesh.AABB.Min[0] != -1 || mesh.AABB.Max[0] != 1 {
		t.Error("Expected AABB x range [-1, 1]")
	}
}

// TestIntegrationComplexScene 测试复杂场景集成
func TestIntegrationComplexScene(t *testing.T) {
	scene := &Scene{
		Flags: SceneFlagValidated,
		RootNode: &Node{
			Name: "world",
			Children: []*Node{
				{
					Name:         "character",
					MeshIndicies: []uint{0, 1, 2},
					Children: []*Node{
						{
							Name:         "head",
							MeshIndicies: []uint{0},
							Children: []*Node{
								{
									Name:         "eye_left",
									MeshIndicies: []uint{1},
								},
								{
									Name:         "eye_right",
									MeshIndicies: []uint{2},
								},
							},
						},
						{
							Name:         "body",
							MeshIndicies: []uint{3},
						},
					},
				},
				{
					Name:         "environment",
					MeshIndicies: []uint{4, 5, 6, 7},
					Children: []*Node{
						{
							Name:         "building",
							MeshIndicies: []uint{4, 5},
						},
						{
							Name:         "terrain",
							MeshIndicies: []uint{6},
						},
					},
				},
			},
		},
		Meshes: []*Mesh{
			{Name: "head_mesh", Vertices: make([]vec3.T, 500), Faces: make([]Face, 300)},
			{Name: "eye_left_mesh", Vertices: make([]vec3.T, 50), Faces: make([]Face, 30)},
			{Name: "eye_right_mesh", Vertices: make([]vec3.T, 50), Faces: make([]Face, 30)},
			{Name: "body_mesh", Vertices: make([]vec3.T, 1000), Faces: make([]Face, 600)},
			{Name: "building_mesh", Vertices: make([]vec3.T, 2000), Faces: make([]Face, 1200)},
			{Name: "building_detail_mesh", Vertices: make([]vec3.T, 500), Faces: make([]Face, 300)},
			{Name: "terrain_mesh", Vertices: make([]vec3.T, 10000), Faces: make([]Face, 5000)},
			{Name: "prop_mesh", Vertices: make([]vec3.T, 100), Faces: make([]Face, 50)},
		},
		Materials: []*Material{
			{Properties: []*MaterialProperty{{name: "skin_material"}}},
			{Properties: []*MaterialProperty{{name: "eye_material"}}},
			{Properties: []*MaterialProperty{{name: "body_material"}}},
			{Properties: []*MaterialProperty{{name: "building_material"}}},
			{Properties: []*MaterialProperty{{name: "terrain_material"}}},
			{Properties: []*MaterialProperty{{name: "prop_material"}}},
		},
		Textures: []*EmbeddedTexture{
			{Filename: "skin_texture.png"},
			{Filename: "eye_texture.png"},
			{Filename: "body_texture.png"},
			{Filename: "building_texture.png"},
			{Filename: "terrain_texture.png"},
		},
		Animations: []*Animation{
			{}, {}, // 2个动画
		},
		Lights: []*Light{
			{}, {}, {}, // 3个光源
		},
		Cameras: []*Camera{
			{}, // 1个相机
		},
	}

	// 验证总数量
	if len(scene.Meshes) != 8 {
		t.Errorf("Expected 8 meshes, got %d", len(scene.Meshes))
	}

	if len(scene.Materials) != 6 {
		t.Errorf("Expected 6 materials, got %d", len(scene.Materials))
	}

	if len(scene.Textures) != 5 {
		t.Errorf("Expected 5 textures, got %d", len(scene.Textures))
	}

	if len(scene.Animations) != 2 {
		t.Errorf("Expected 2 animations, got %d", len(scene.Animations))
	}

	if len(scene.Lights) != 3 {
		t.Errorf("Expected 3 lights, got %d", len(scene.Lights))
	}

	if len(scene.Cameras) != 1 {
		t.Errorf("Expected 1 camera, got %d", len(scene.Cameras))
	}

	// 验证节点层次结构
	if scene.RootNode == nil {
		t.Fatal("Expected non-nil root node")
	}

	if len(scene.RootNode.Children) != 2 {
		t.Errorf("Expected root to have 2 children, got %d", len(scene.RootNode.Children))
	}

	// 验证character节点
	character := scene.RootNode.Children[0]
	if character.Name != "character" {
		t.Errorf("Expected first child name 'character', got '%s'", character.Name)
	}

	if len(character.Children) != 2 {
		t.Errorf("Expected character to have 2 children, got %d", len(character.Children))
	}

	// 验证head节点
	head := character.Children[0]
	if head.Name != "head" {
		t.Errorf("Expected head node name 'head', got '%s'", head.Name)
	}

	if len(head.Children) != 2 {
		t.Errorf("Expected head to have 2 children, got %d", len(head.Children))
	}

	// 验证environment节点
	environment := scene.RootNode.Children[1]
	if environment.Name != "environment" {
		t.Errorf("Expected second child name 'environment', got '%s'", environment.Name)
	}

	// 验证网格分配
	if len(character.MeshIndicies) != 3 {
		t.Errorf("Expected character to have 3 mesh indices, got %d", len(character.MeshIndicies))
	}

	if len(environment.MeshIndicies) != 4 {
		t.Errorf("Expected environment to have 4 mesh indices, got %d", len(environment.MeshIndicies))
	}

	// 验证网格顶点数量
	expectedVertices := []int{500, 50, 50, 1000, 2000, 500, 10000, 100}
	for i, expected := range expectedVertices {
		if len(scene.Meshes[i].Vertices) != expected {
			t.Errorf("Expected mesh %d to have %d vertices, got %d",
				i, expected, len(scene.Meshes[i].Vertices))
		}
	}

	// 验证网格面数量
	expectedFaces := []int{300, 30, 30, 600, 1200, 300, 5000, 50}
	for i, expected := range expectedFaces {
		if len(scene.Meshes[i].Faces) != expected {
			t.Errorf("Expected mesh %d to have %d faces, got %d",
				i, expected, len(scene.Meshes[i].Faces))
		}
	}
}

// TestIntegrationEmptyScene 测试空场景集成
func TestIntegrationEmptyScene(t *testing.T) {
	scene := &Scene{}

	if scene.RootNode != nil {
		t.Error("Expected nil root node for empty scene")
	}

	if len(scene.Meshes) != 0 {
		t.Error("Expected empty meshes")
	}

	if len(scene.Materials) != 0 {
		t.Error("Expected empty materials")
	}

	if len(scene.Textures) != 0 {
		t.Error("Expected empty textures")
	}

	if len(scene.Animations) != 0 {
		t.Error("Expected empty animations")
	}

	if len(scene.Lights) != 0 {
		t.Error("Expected empty lights")
	}

	if len(scene.Cameras) != 0 {
		t.Error("Expected empty cameras")
	}
}

// TestIntegrationSingleMeshScene 测试单网格场景集成
func TestIntegrationSingleMeshScene(t *testing.T) {
	scene := &Scene{
		RootNode: &Node{
			Name:         "root",
			MeshIndicies: []uint{0},
		},
		Meshes: []*Mesh{
			{
				Name:     "single_mesh",
				Vertices: make([]vec3.T, 1000),
				Normals:  make([]vec3.T, 1000),
				Faces:    make([]Face, 500),
				AABB: AABB{
					Min: vec3.T{-10, -10, -10},
					Max: vec3.T{10, 10, 10},
				},
			},
		},
		Materials: []*Material{
			{
				Properties: []*MaterialProperty{
					{name: "material_property"},
				},
			},
		},
	}

	if len(scene.Meshes) != 1 {
		t.Errorf("Expected 1 mesh, got %d", len(scene.Meshes))
	}

	if len(scene.Materials) != 1 {
		t.Errorf("Expected 1 material, got %d", len(scene.Materials))
	}

	if scene.RootNode.MeshIndicies[0] != 0 {
		t.Errorf("Expected mesh index 0, got %d", scene.RootNode.MeshIndicies[0])
	}

	mesh := scene.Meshes[0]
	if len(mesh.Vertices) != 1000 {
		t.Errorf("Expected 1000 vertices, got %d", len(mesh.Vertices))
	}

	if len(mesh.Faces) != 500 {
		t.Errorf("Expected 500 faces, got %d", len(mesh.Faces))
	}
}
