package assimp

import (
	"testing"

	"github.com/flywave/go3d/mat4"
	"github.com/flywave/go3d/vec3"
)

// TestSceneCreation 测试场景创建
func TestSceneCreation(t *testing.T) {
	scene := &Scene{
		Flags:      SceneFlagValidated,
		RootNode:   &Node{Name: "root"},
		Meshes:     []*Mesh{{Name: "mesh1"}},
		Materials:  []*Material{{Properties: []*MaterialProperty{{name: "mat1"}}}},
		Textures:   []*EmbeddedTexture{{Filename: "texture1.png"}},
		Animations: []*Animation{{}},
		Lights:     []*Light{{}},
		Cameras:    []*Camera{{}},
	}

	if scene.Flags != SceneFlagValidated {
		t.Errorf("Expected SceneFlagValidated, got %d", scene.Flags)
	}

	if scene.RootNode == nil {
		t.Error("Expected non-nil root node")
	}

	if len(scene.Meshes) != 1 {
		t.Errorf("Expected 1 mesh, got %d", len(scene.Meshes))
	}

	if len(scene.Materials) != 1 {
		t.Errorf("Expected 1 material, got %d", len(scene.Materials))
	}

	if len(scene.Textures) != 1 {
		t.Errorf("Expected 1 texture, got %d", len(scene.Textures))
	}

	if len(scene.Animations) != 1 {
		t.Errorf("Expected 1 animation, got %d", len(scene.Animations))
	}

	if len(scene.Lights) != 1 {
		t.Errorf("Expected 1 light, got %d", len(scene.Lights))
	}

	if len(scene.Cameras) != 1 {
		t.Errorf("Expected 1 camera, got %d", len(scene.Cameras))
	}
}

// TestNodeCreation 测试节点创建
func TestNodeCreation(t *testing.T) {
	node := &Node{
		Name:           "test_node",
		Transformation: &mat4.T{{1, 0, 0, 0}, {0, 1, 0, 0}, {0, 0, 1, 0}, {0, 0, 0, 1}},
		Parent:         nil,
		Children: []*Node{
			{Name: "child1"},
			{Name: "child2"},
		},
		MeshIndicies: []uint{0, 1, 2},
		Metadata: map[string]Metadata{
			"test_key": {
				Type:  MetadataTypeString,
				Value: "test_value",
			},
		},
	}

	if node.Name != "test_node" {
		t.Errorf("Expected node name 'test_node', got '%s'", node.Name)
	}

	if node.Transformation == nil {
		t.Error("Expected non-nil transformation")
	}

	if len(node.Children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(node.Children))
	}

	if len(node.MeshIndicies) != 3 {
		t.Errorf("Expected 3 mesh indices, got %d", len(node.MeshIndicies))
	}

	if len(node.Metadata) != 1 {
		t.Errorf("Expected 1 metadata entry, got %d", len(node.Metadata))
	}
}

// TestNodeHierarchy 测试节点层次结构
func TestNodeHierarchy(t *testing.T) {
	root := &Node{Name: "root"}
	child1 := &Node{Name: "child1", Parent: root}
	child2 := &Node{Name: "child2", Parent: root}
	grandchild := &Node{Name: "grandchild", Parent: child1}

	root.Children = []*Node{child1, child2}
	child1.Children = []*Node{grandchild}

	if root.Parent != nil {
		t.Error("Expected root node to have no parent")
	}

	if child1.Parent != root {
		t.Error("Expected child1 parent to be root")
	}

	if child2.Parent != root {
		t.Error("Expected child2 parent to be root")
	}

	if grandchild.Parent != child1 {
		t.Error("Expected grandchild parent to be child1")
	}

	if len(root.Children) != 2 {
		t.Errorf("Expected root to have 2 children, got %d", len(root.Children))
	}

	if len(child1.Children) != 1 {
		t.Errorf("Expected child1 to have 1 child, got %d", len(child1.Children))
	}
}

// TestEmptyScene 测试空场景
func TestEmptyScene(t *testing.T) {
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

// TestEmptyNode 测试空节点
func TestEmptyNode(t *testing.T) {
	node := &Node{}

	if node.Name != "" {
		t.Errorf("Expected empty name, got '%s'", node.Name)
	}

	if node.Transformation != nil {
		t.Error("Expected nil transformation")
	}

	if node.Parent != nil {
		t.Error("Expected nil parent")
	}

	if len(node.Children) != 0 {
		t.Error("Expected empty children")
	}

	if len(node.MeshIndicies) != 0 {
		t.Error("Expected empty mesh indices")
	}

	if len(node.Metadata) != 0 {
		t.Error("Expected empty metadata")
	}
}

// TestSceneFlagOperations 测试场景标志操作
func TestSceneFlagOperations(t *testing.T) {
	flags := []SceneFlag{
		SceneFlagIncomplete,
		SceneFlagValidated,
		SceneFlagValidationWarning,
		SceneFlagNonVerboseFormat,
		SceneFlagTerrain,
		SceneFlagAllowShared,
	}

	for _, flag := range flags {
		if flag <= 0 {
			t.Errorf("Expected positive flag value, got %d", flag)
		}
	}
}

// TestNodeMeshIndicies 测试节点网格索引
func TestNodeMeshIndicies(t *testing.T) {
	node := &Node{
		MeshIndicies: []uint{0, 1, 2, 3, 4},
	}

	if len(node.MeshIndicies) != 5 {
		t.Errorf("Expected 5 mesh indices, got %d", len(node.MeshIndicies))
	}

	for i, index := range node.MeshIndicies {
		if index != uint(i) {
			t.Errorf("Expected index %d, got %d", i, index)
		}
	}
}

// TestNodeMetadata 测试节点元数据
func TestNodeMetadata(t *testing.T) {
	node := &Node{
		Metadata: map[string]Metadata{
			"author": {
				Type:  MetadataTypeString,
				Value: "test_author",
			},
			"version": {
				Type:  MetadataTypeInt32,
				Value: int32(1),
			},
			"scale": {
				Type:  MetadataTypeFloat32,
				Value: float32(1.0),
			},
			"visible": {
				Type:  MetadataTypeBool,
				Value: true,
			},
		},
	}

	if len(node.Metadata) != 4 {
		t.Errorf("Expected 4 metadata entries, got %d", len(node.Metadata))
	}

	if node.Metadata["author"].Type != MetadataTypeString {
		t.Errorf("Expected string type for 'author' metadata")
	}

	if node.Metadata["version"].Type != MetadataTypeInt32 {
		t.Errorf("Expected int32 type for 'version' metadata")
	}

	if node.Metadata["scale"].Type != MetadataTypeFloat32 {
		t.Errorf("Expected float32 type for 'scale' metadata")
	}

	if node.Metadata["visible"].Type != MetadataTypeBool {
		t.Errorf("Expected bool type for 'visible' metadata")
	}
}

// TestNodeTransformation 测试节点变换
func TestNodeTransformation(t *testing.T) {
	identity := mat4.T{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}

	node := &Node{
		Transformation: &identity,
	}

	if node.Transformation == nil {
		t.Error("Expected non-nil transformation")
	}

	if (*node.Transformation)[0][0] != 1 {
		t.Error("Expected identity matrix")
	}
}

// TestComplexScene 测试复杂场景
func TestComplexScene(t *testing.T) {
	scene := &Scene{
		Flags: SceneFlagValidated | SceneFlagTerrain,
		RootNode: &Node{
			Name: "root",
			Children: []*Node{
				{
					Name:         "character",
					MeshIndicies: []uint{0, 1},
					Children: []*Node{
						{
							Name:         "head",
							MeshIndicies: []uint{2},
						},
						{
							Name:         "body",
							MeshIndicies: []uint{3},
						},
					},
				},
				{
					Name:         "environment",
					MeshIndicies: []uint{4, 5, 6},
				},
			},
		},
		Meshes: []*Mesh{
			{Name: "character_mesh1"},
			{Name: "character_mesh2"},
			{Name: "head_mesh"},
			{Name: "body_mesh"},
			{Name: "ground"},
			{Name: "building"},
			{Name: "tree"},
		},
		Materials: []*Material{
			{Properties: []*MaterialProperty{{name: "skin_mat"}}},
			{Properties: []*MaterialProperty{{name: "clothes_mat"}}},
			{Properties: []*MaterialProperty{{name: "hair_mat"}}},
			{Properties: []*MaterialProperty{{name: "environment_mat"}}},
		},
	}

	if len(scene.Meshes) != 7 {
		t.Errorf("Expected 7 meshes, got %d", len(scene.Meshes))
	}

	if len(scene.Materials) != 4 {
		t.Errorf("Expected 4 materials, got %d", len(scene.Materials))
	}

	// 验证节点层次结构
	if scene.RootNode == nil {
		t.Fatal("Expected non-nil root node")
	}

	if len(scene.RootNode.Children) != 2 {
		t.Errorf("Expected root to have 2 children, got %d", len(scene.RootNode.Children))
	}

	character := scene.RootNode.Children[0]
	if character.Name != "character" {
		t.Errorf("Expected first child name 'character', got '%s'", character.Name)
	}

	if len(character.Children) != 2 {
		t.Errorf("Expected character to have 2 children, got %d", len(character.Children))
	}
}

// TestSceneBounds 测试场景边界
func TestSceneBounds(t *testing.T) {
	scene := &Scene{
		Meshes: []*Mesh{
			{
				Vertices: []vec3.T{
					{-1, -1, -1},
					{1, 1, 1},
				},
				AABB: AABB{
					Min: vec3.T{-1, -1, -1},
					Max: vec3.T{1, 1, 1},
				},
			},
		},
	}

	if len(scene.Meshes) != 1 {
		t.Errorf("Expected 1 mesh, got %d", len(scene.Meshes))
	}

	mesh := scene.Meshes[0]
	if mesh.AABB.Min[0] != -1 {
		t.Errorf("Expected min x -1, got %f", mesh.AABB.Min[0])
	}

	if mesh.AABB.Max[0] != 1 {
		t.Errorf("Expected max x 1, got %f", mesh.AABB.Max[0])
	}
}
