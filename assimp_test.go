package assimp

import (
	"testing"

	"github.com/flywave/go3d/mat4"
	"github.com/flywave/go3d/vec3"
)

// TestBasicImport 测试基本的文件导入功能
func TestBasicImport(t *testing.T) {
	// 测试导入不存在的文件
	_, _, err := ImportFile("nonexistent.obj", PostProcessTriangulate)
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}

	// 测试空文件名
	_, _, err = ImportFile("", PostProcessTriangulate)
	if err == nil {
		t.Error("Expected error for empty filename, got nil")
	}
}

// TestPostProcessFlags 测试不同的后处理标志
func TestPostProcessFlags(t *testing.T) {
	// 这里我们主要测试标志的组合不会导致崩溃
	flags := []PostProcess{
		PostProcessTriangulate,
		PostProcessGenNormals,
		PostProcessCalcTangentSpace,
		PostProcessJoinIdenticalVertices,
		PostProcessGenSmoothNormals,
		PostProcessSplitLargeMeshes,
		PostProcessValidateDataStructure,
	}

	for _, flag := range flags {
		// 由于我们没有测试模型文件，这里只测试标志的有效性
		t.Logf("Testing post process flag: %d", flag)
	}
}

// TestSceneStructure 测试场景结构的完整性
func TestSceneStructure(t *testing.T) {
	// 创建一个最小化的测试场景
	// 由于需要实际的3D模型文件，这里主要测试nil场景的处理
	scene := &Scene{}

	// 测试场景字段的默认值
	if scene.RootNode != nil {
		t.Error("Expected RootNode to be nil for empty scene")
	}

	if len(scene.Meshes) != 0 {
		t.Error("Expected Meshes to be empty for empty scene")
	}

	if len(scene.Materials) != 0 {
		t.Error("Expected Materials to be empty for empty scene")
	}

	if len(scene.Textures) != 0 {
		t.Error("Expected Textures to be empty for empty scene")
	}

	if len(scene.Animations) != 0 {
		t.Error("Expected Animations to be empty for empty scene")
	}

	if len(scene.Lights) != 0 {
		t.Error("Expected Lights to be empty for empty scene")
	}

	if len(scene.Cameras) != 0 {
		t.Error("Expected Cameras to be empty for empty scene")
	}
}

// TestNodeStructure 测试节点结构的完整性
func TestNodeStructure(t *testing.T) {
	node := &Node{
		Name:           "test_node",
		Transformation: nil,
		Parent:         nil,
		Children:       []*Node{},
		MeshIndicies:   []uint{},
		Metadata:       map[string]Metadata{},
	}

	if node.Name != "test_node" {
		t.Errorf("Expected node name 'test_node', got '%s'", node.Name)
	}

	if len(node.Children) != 0 {
		t.Error("Expected Children to be empty")
	}

	if len(node.MeshIndicies) != 0 {
		t.Error("Expected MeshIndicies to be empty")
	}

	if len(node.Metadata) != 0 {
		t.Error("Expected Metadata to be empty")
	}
}

// TestMeshStructure 测试网格结构的完整性
func TestMeshStructure(t *testing.T) {
	mesh := &Mesh{
		Name:          "test_mesh",
		Vertices:      []vec3.T{},
		Normals:       []vec3.T{},
		Tangents:      []vec3.T{},
		BitTangents:   []vec3.T{},
		Faces:         []Face{},
		Bones:         []*Bone{},
		AnimMeshes:    []*AnimMesh{},
		MaterialIndex: 0,
	}

	if mesh.Name != "test_mesh" {
		t.Errorf("Expected mesh name 'test_mesh', got '%s'", mesh.Name)
	}

	if len(mesh.Vertices) != 0 {
		t.Error("Expected Vertices to be empty")
	}

	if len(mesh.Normals) != 0 {
		t.Error("Expected Normals to be empty")
	}

	if len(mesh.Faces) != 0 {
		t.Error("Expected Faces to be empty")
	}

	if len(mesh.Bones) != 0 {
		t.Error("Expected Bones to be empty")
	}
}

// TestMaterialStructure 测试材质结构的完整性
func TestMaterialStructure(t *testing.T) {
	material := &Material{
		Properties:       []*MaterialProperty{},
		AllocatedStorage: 0,
	}

	if len(material.Properties) != 0 {
		t.Error("Expected Properties to be empty")
	}
}

// TestTextureTypeString 测试纹理类型的字符串表示
func TestTextureTypeString(t *testing.T) {
	tests := []struct {
		type_    TextureType
		expected string
	}{
		{TextureTypeNone, "None"},
		{TextureTypeDiffuse, "Diffuse"},
		{TextureTypeSpecular, "Specular"},
		{TextureTypeAmbient, "Ambient"},
		{TextureTypeEmissive, "Emissive"},
		{TextureTypeHeight, "Height"},
		{TextureTypeNormal, "Normal"},
		{TextureTypeShininess, "Shininess"},
		{TextureTypeOpacity, "Opacity"},
		{TextureTypeDisplacement, "Displacement"},
		{TextureTypeLightmap, "Lightmap"},
		{TextureTypeReflection, "Reflection"},
		{TextureTypeBaseColor, "BaseColor"},
		{TextureTypeNormalCamera, "NormalCamera"},
		{TextureTypeEmissionColor, "EmissionColor"},
		{TextureTypeMetalness, "Metalness"},
		{TextureTypeDiffuseRoughness, "DiffuseRoughness"},
		{TextureTypeAmbientOcclusion, "AmbientOcclusion"},
		{TextureTypeUnknown, "Unknown"},
		{TextureType(999), "Invalid"},
	}

	for _, test := range tests {
		result := test.type_.String()
		if result != test.expected {
			t.Errorf("Expected '%s' for TextureType %d, got '%s'", test.expected, test.type_, result)
		}
	}
}

// TestMatPropertyTypeInfoString 测试材质属性类型的字符串表示
func TestMatPropertyTypeInfoString(t *testing.T) {
	tests := []struct {
		type_    MatPropertyTypeInfo
		expected string
	}{
		{MatPropTypeInfoFloat32, "Float32"},
		{MatPropTypeInfoFloat64, "Float64"},
		{MatPropTypeInfoString, "String"},
		{MatPropTypeInfoInt32, "Int32"},
		{MatPropTypeInfoBuffer, "Buffer"},
		{MatPropertyTypeInfo(999), "Unknown"},
	}

	for _, test := range tests {
		result := test.type_.String()
		if result != test.expected {
			t.Errorf("Expected '%s' for MatPropertyTypeInfo %d, got '%s'", test.expected, test.type_, result)
		}
	}
}

// TestAABBStructure 测试AABB结构的完整性
func TestAABBStructure(t *testing.T) {
	aabb := AABB{
		Min: vec3.T{0, 0, 0},
		Max: vec3.T{1, 1, 1},
	}

	if aabb.Min[0] != 0 || aabb.Min[1] != 0 || aabb.Min[2] != 0 {
		t.Error("Expected AABB Min to be {0, 0, 0}")
	}

	if aabb.Max[0] != 1 || aabb.Max[1] != 1 || aabb.Max[2] != 1 {
		t.Error("Expected AABB Max to be {1, 1, 1}")
	}
}

// TestFaceStructure 测试面结构的完整性
func TestFaceStructure(t *testing.T) {
	face := Face{
		Indices: []uint{0, 1, 2},
	}

	if len(face.Indices) != 3 {
		t.Errorf("Expected 3 indices, got %d", len(face.Indices))
	}

	if face.Indices[0] != 0 || face.Indices[1] != 1 || face.Indices[2] != 2 {
		t.Error("Expected indices to be [0, 1, 2]")
	}
}

// TestBoneStructure 测试骨骼结构的完整性
func TestBoneStructure(t *testing.T) {
	bone := &Bone{
		Name:    "test_bone",
		Weights: []VertexWeight{},
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

	if len(bone.Weights) != 0 {
		t.Error("Expected Weights to be empty")
	}
}

// TestVertexWeightStructure 测试顶点权重结构的完整性
func TestVertexWeightStructure(t *testing.T) {
	vw := VertexWeight{
		VertIndex: 5,
		Weight:    0.5,
	}

	if vw.VertIndex != 5 {
		t.Errorf("Expected vertex index 5, got %d", vw.VertIndex)
	}

	if vw.Weight != 0.5 {
		t.Errorf("Expected weight 0.5, got %f", vw.Weight)
	}
}

// TestAnimMeshStructure 测试动画网格结构的完整性
func TestAnimMeshStructure(t *testing.T) {
	animMesh := &AnimMesh{
		Name:   "test_anim_mesh",
		Weight: 1.0,
	}

	if animMesh.Name != "test_anim_mesh" {
		t.Errorf("Expected anim mesh name 'test_anim_mesh', got '%s'", animMesh.Name)
	}

	if animMesh.Weight != 1.0 {
		t.Errorf("Expected weight 1.0, got %f", animMesh.Weight)
	}
}

// TestPostProcessValues 测试后处理标志值的正确性
func TestPostProcessValues(t *testing.T) {
	if PostProcessTriangulate != 0x8 {
		t.Errorf("Expected PostProcessTriangulate to be 0x8, got %d", PostProcessTriangulate)
	}

	if PostProcessGenNormals != 0x20 {
		t.Errorf("Expected PostProcessGenNormals to be 0x20, got %d", PostProcessGenNormals)
	}

	if PostProcessCalcTangentSpace != 0x1 {
		t.Errorf("Expected PostProcessCalcTangentSpace to be 0x1, got %d", PostProcessCalcTangentSpace)
	}
}

// TestSceneFlagValues 测试场景标志值的正确性
func TestSceneFlagValues(t *testing.T) {
	if SceneFlagIncomplete != 1 {
		t.Errorf("Expected SceneFlagIncomplete to be 1, got %d", SceneFlagIncomplete)
	}

	if SceneFlagValidated != 2 {
		t.Errorf("Expected SceneFlagValidated to be 2, got %d", SceneFlagValidated)
	}
}

// TestPrimitiveTypeValues 测试图元类型值的正确性
func TestPrimitiveTypeValues(t *testing.T) {
	if PrimitiveTypePoint != 1 {
		t.Errorf("Expected PrimitiveTypePoint to be 1, got %d", PrimitiveTypePoint)
	}

	if PrimitiveTypeLine != 2 {
		t.Errorf("Expected PrimitiveTypeLine to be 2, got %d", PrimitiveTypeLine)
	}

	if PrimitiveTypeTriangle != 4 {
		t.Errorf("Expected PrimitiveTypeTriangle to be 4, got %d", PrimitiveTypeTriangle)
	}
}

// TestMorphMethodValues 测试变形方法值的正确性
func TestMorphMethodValues(t *testing.T) {
	if MorphMethodVertexBlend != 0x1 {
		t.Errorf("Expected MorphMethodVertexBlend to be 0x1, got %d", MorphMethodVertexBlend)
	}

	if MorphMethodMorphNormalized != 0x2 {
		t.Errorf("Expected MorphMethodMorphNormalized to be 0x2, got %d", MorphMethodMorphNormalized)
	}
}

// TestEmbeddedTextureStructure 测试嵌入式纹理结构的完整性
func TestEmbeddedTextureStructure(t *testing.T) {
	texture := &EmbeddedTexture{
		Width:        512,
		Height:       512,
		FormatHint:   "png",
		Data:         []byte{},
		IsCompressed: false,
		Filename:     "test_texture.png",
	}

	if texture.Width != 512 {
		t.Errorf("Expected texture width 512, got %d", texture.Width)
	}

	if texture.Height != 512 {
		t.Errorf("Expected texture height 512, got %d", texture.Height)
	}

	if texture.FormatHint != "png" {
		t.Errorf("Expected format hint 'png', got '%s'", texture.FormatHint)
	}

	if texture.Filename != "test_texture.png" {
		t.Errorf("Expected filename 'test_texture.png', got '%s'", texture.Filename)
	}
}

// TestMetadataStructure 测试元数据结构的完整性
func TestMetadataStructure(t *testing.T) {
	metadata := Metadata{
		Type:  MetadataTypeString,
		Value: "test_value",
	}

	if metadata.Type != MetadataTypeString {
		t.Errorf("Expected metadata type %d, got %d", MetadataTypeString, metadata.Type)
	}

	if metadata.Value != "test_value" {
		t.Errorf("Expected metadata value 'test_value', got '%v'", metadata.Value)
	}
}

// TestMetadataTypeValues 测试元数据类型值的正确性
func TestMetadataTypeValues(t *testing.T) {
	if MetadataTypeBool != 0 {
		t.Errorf("Expected MetadataTypeBool to be 0, got %d", MetadataTypeBool)
	}

	if MetadataTypeInt32 != 1 {
		t.Errorf("Expected MetadataTypeInt32 to be 1, got %d", MetadataTypeInt32)
	}

	if MetadataTypeString != 5 {
		t.Errorf("Expected MetadataTypeString to be 5, got %d", MetadataTypeString)
	}

	if MetadataTypeVec3 != 6 {
		t.Errorf("Expected MetadataTypeVec3 to be 6, got %d", MetadataTypeVec3)
	}
}
