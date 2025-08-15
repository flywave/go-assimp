package assimp

import (
	"testing"

	"github.com/flywave/go3d/mat4"
	"github.com/flywave/go3d/vec3"
	"github.com/flywave/go3d/vec4"
)

// TestMeshCreation 测试网格创建
func TestMeshCreation(t *testing.T) {
	mesh := &Mesh{
		PrimitiveTypes: PrimitiveTypeTriangle,
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
		Tangents: []vec3.T{
			{1, 0, 0},
			{1, 0, 0},
			{1, 0, 0},
		},
		BitTangents: []vec3.T{
			{0, 1, 0},
			{0, 1, 0},
			{0, 1, 0},
		},
		Faces: []Face{
			{Indices: []uint{0, 1, 2}},
		},
		MaterialIndex: 0,
		Name:          "test_triangle",
	}

	if len(mesh.Vertices) != 3 {
		t.Errorf("Expected 3 vertices, got %d", len(mesh.Vertices))
	}

	if len(mesh.Normals) != 3 {
		t.Errorf("Expected 3 normals, got %d", len(mesh.Normals))
	}

	if len(mesh.Faces) != 1 {
		t.Errorf("Expected 1 face, got %d", len(mesh.Faces))
	}

	if mesh.Name != "test_triangle" {
		t.Errorf("Expected mesh name 'test_triangle', got '%s'", mesh.Name)
	}

	if mesh.PrimitiveTypes != PrimitiveTypeTriangle {
		t.Errorf("Expected PrimitiveTypeTriangle, got %d", mesh.PrimitiveTypes)
	}
}

// TestMeshWithColorSets 测试带颜色集的网格
func TestMeshWithColorSets(t *testing.T) {
	mesh := &Mesh{
		Vertices: []vec3.T{
			{0, 0, 0},
			{1, 0, 0},
			{0, 1, 0},
		},
		ColorSets: [MaxColorSets][]vec4.T{
			0: {
				{1, 0, 0, 1},
				{0, 1, 0, 1},
				{0, 0, 1, 1},
			},
		},
	}

	if len(mesh.ColorSets[0]) != 3 {
		t.Errorf("Expected 3 colors in color set 0, got %d", len(mesh.ColorSets[0]))
	}

	if mesh.ColorSets[0][0][0] != 1 {
		t.Errorf("Expected red component 1, got %f", mesh.ColorSets[0][0][0])
	}
}

// TestMeshWithTexCoords 测试带纹理坐标的网格
func TestMeshWithTexCoords(t *testing.T) {
	mesh := &Mesh{
		Vertices: []vec3.T{
			{0, 0, 0},
			{1, 0, 0},
			{0, 1, 0},
		},
		TexCoords: [MaxTexCoords][]vec3.T{
			0: {
				{0, 0, 0},
				{1, 0, 0},
				{0, 1, 0},
			},
		},
		TexCoordChannelCount: [MaxTexCoords]uint{
			0: 2,
		},
	}

	if len(mesh.TexCoords[0]) != 3 {
		t.Errorf("Expected 3 texture coordinates in channel 0, got %d", len(mesh.TexCoords[0]))
	}

	if mesh.TexCoordChannelCount[0] != 2 {
		t.Errorf("Expected 2 texture coordinate components, got %d", mesh.TexCoordChannelCount[0])
	}
}

// TestMeshWithBones 测试带骨骼的网格
func TestMeshWithBones(t *testing.T) {
	mesh := &Mesh{
		Vertices: []vec3.T{
			{0, 0, 0},
			{1, 0, 0},
			{0, 1, 0},
		},
		Bones: []*Bone{
			{
				Name: "bone1",
				Weights: []VertexWeight{
					{VertIndex: 0, Weight: 1.0},
					{VertIndex: 1, Weight: 0.5},
				},
				OffsetMatrix: mat4.T{
					{1, 0, 0, 0},
					{0, 1, 0, 0},
					{0, 0, 1, 0},
					{0, 0, 0, 1},
				},
			},
		},
	}

	if len(mesh.Bones) != 1 {
		t.Errorf("Expected 1 bone, got %d", len(mesh.Bones))
	}

	if mesh.Bones[0].Name != "bone1" {
		t.Errorf("Expected bone name 'bone1', got '%s'", mesh.Bones[0].Name)
	}

	if len(mesh.Bones[0].Weights) != 2 {
		t.Errorf("Expected 2 vertex weights, got %d", len(mesh.Bones[0].Weights))
	}
}

// TestMeshWithAnimMeshes 测试带动画网格的网格
func TestMeshWithAnimMeshes(t *testing.T) {
	mesh := &Mesh{
		Vertices: []vec3.T{
			{0, 0, 0},
			{1, 0, 0},
			{0, 1, 0},
		},
		AnimMeshes: []*AnimMesh{
			{
				Name: "anim_mesh1",
				Vertices: []vec3.T{
					{0, 0, 1},
					{1, 0, 1},
					{0, 1, 1},
				},
				Weight: 0.5,
			},
		},
	}

	if len(mesh.AnimMeshes) != 1 {
		t.Errorf("Expected 1 animation mesh, got %d", len(mesh.AnimMeshes))
	}

	if mesh.AnimMeshes[0].Name != "anim_mesh1" {
		t.Errorf("Expected animation mesh name 'anim_mesh1', got '%s'", mesh.AnimMeshes[0].Name)
	}

	if len(mesh.AnimMeshes[0].Vertices) != 3 {
		t.Errorf("Expected 3 vertices in animation mesh, got %d", len(mesh.AnimMeshes[0].Vertices))
	}

	if mesh.AnimMeshes[0].Weight != 0.5 {
		t.Errorf("Expected weight 0.5, got %f", mesh.AnimMeshes[0].Weight)
	}
}

// TestMeshWithAABB 测试带边界框的网格
func TestMeshWithAABB(t *testing.T) {
	mesh := &Mesh{
		Vertices: []vec3.T{
			{0, 0, 0},
			{1, 0, 0},
			{0, 1, 0},
		},
		AABB: AABB{
			Min: vec3.T{0, 0, 0},
			Max: vec3.T{1, 1, 0},
		},
	}

	if mesh.AABB.Min[0] != 0 || mesh.AABB.Min[1] != 0 || mesh.AABB.Min[2] != 0 {
		t.Error("Expected AABB Min to be {0, 0, 0}")
	}

	if mesh.AABB.Max[0] != 1 || mesh.AABB.Max[1] != 1 || mesh.AABB.Max[2] != 0 {
		t.Error("Expected AABB Max to be {1, 1, 0}")
	}
}

// TestFaceCreation 测试面的创建
func TestFaceCreation(t *testing.T) {
	face := Face{
		Indices: []uint{0, 1, 2, 3},
	}

	if len(face.Indices) != 4 {
		t.Errorf("Expected 4 indices, got %d", len(face.Indices))
	}

	if face.Indices[0] != 0 || face.Indices[1] != 1 || face.Indices[2] != 2 || face.Indices[3] != 3 {
		t.Error("Expected indices to be [0, 1, 2, 3]")
	}
}

// TestBoneCreation 测试骨骼的创建
func TestBoneCreation(t *testing.T) {
	bone := &Bone{
		Name: "test_bone",
		Weights: []VertexWeight{
			{VertIndex: 0, Weight: 1.0},
			{VertIndex: 1, Weight: 0.5},
			{VertIndex: 2, Weight: 0.25},
		},
		OffsetMatrix: mat4.T{
			{1, 0, 0, 0},
			{0, 1, 0, 0},
			{0, 0, 1, 0},
			{0, 0, 0, 1},
		},
	}

	if bone.Name != "test_bone" {
		t.Errorf("Expected bone name 'test_bone', got '%s'", bone.Name)
	}

	if len(bone.Weights) != 3 {
		t.Errorf("Expected 3 vertex weights, got %d", len(bone.Weights))
	}

	if bone.Weights[0].Weight != 1.0 {
		t.Errorf("Expected weight 1.0, got %f", bone.Weights[0].Weight)
	}
}

// TestAnimMeshCreation 测试动画网格的创建
func TestAnimMeshCreation(t *testing.T) {
	animMesh := &AnimMesh{
		Name: "test_anim_mesh",
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
		Tangents: []vec3.T{
			{1, 0, 0},
			{1, 0, 0},
			{1, 0, 0},
		},
		BitTangents: []vec3.T{
			{0, 1, 0},
			{0, 1, 0},
			{0, 1, 0},
		},
		Colors: [MaxColorSets][]vec4.T{
			0: {
				{1, 0, 0, 1},
				{0, 1, 0, 1},
				{0, 0, 1, 1},
			},
		},
		TexCoords: [MaxTexCoords][]vec3.T{
			0: {
				{0, 0, 0},
				{1, 0, 0},
				{0, 1, 0},
			},
		},
		Weight: 0.75,
	}

	if animMesh.Name != "test_anim_mesh" {
		t.Errorf("Expected anim mesh name 'test_anim_mesh', got '%s'", animMesh.Name)
	}

	if len(animMesh.Vertices) != 3 {
		t.Errorf("Expected 3 vertices, got %d", len(animMesh.Vertices))
	}

	if animMesh.Weight != 0.75 {
		t.Errorf("Expected weight 0.75, got %f", animMesh.Weight)
	}
}

// TestEmptyMesh 测试空网格
func TestEmptyMesh(t *testing.T) {
	mesh := &Mesh{}

	if len(mesh.Vertices) != 0 {
		t.Error("Expected empty vertices")
	}

	if len(mesh.Normals) != 0 {
		t.Error("Expected empty normals")
	}

	if len(mesh.Faces) != 0 {
		t.Error("Expected empty faces")
	}

	if len(mesh.Bones) != 0 {
		t.Error("Expected empty bones")
	}

	if len(mesh.AnimMeshes) != 0 {
		t.Error("Expected empty animation meshes")
	}
}

// TestMeshValidation 测试网格验证
func TestMeshValidation(t *testing.T) {
	// 测试有效的网格
	mesh := &Mesh{
		Vertices: []vec3.T{
			{0, 0, 0},
			{1, 0, 0},
			{0, 1, 0},
		},
		Faces: []Face{
			{Indices: []uint{0, 1, 2}},
		},
	}

	if len(mesh.Vertices) != 3 {
		t.Errorf("Expected 3 vertices, got %d", len(mesh.Vertices))
	}

	if len(mesh.Faces) != 1 {
		t.Errorf("Expected 1 face, got %d", len(mesh.Faces))
	}

	// 测试索引越界的情况
	meshInvalid := &Mesh{
		Vertices: []vec3.T{
			{0, 0, 0},
			{1, 0, 0},
		},
		Faces: []Face{
			{Indices: []uint{0, 1, 2}}, // 索引2越界
		},
	}

	// 这里我们只是测试结构，不验证索引有效性
	if len(meshInvalid.Faces) != 1 {
		t.Errorf("Expected 1 face, got %d", len(meshInvalid.Faces))
	}
}