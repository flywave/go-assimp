package assimp

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/flywave/go-mst"
	"github.com/flywave/go3d/mat4"
	"github.com/flywave/go3d/vec3"
)

// TestAssimpToMSTConverter 测试assimp到mst的转换
func TestAssimpToMSTConverter(t *testing.T) {
	// 创建一个简单的测试场景
	scene := createTestScene()

	// 转换为MST格式
	mstMesh := AssimpToMSTConverter(scene)

	if mstMesh == nil {
		t.Fatal("转换结果为nil")
	}

	// 验证转换结果
	if len(mstMesh.Materials) == 0 {
		t.Error("没有转换材质")
	}

	if len(mstMesh.Nodes) == 0 {
		t.Error("没有转换网格节点")
	}

	// 验证材质类型
	for _, material := range mstMesh.Materials {
		if material == nil {
			t.Error("材质为nil")
		}
	}

	// 验证网格节点
	for _, node := range mstMesh.Nodes {
		if len(node.Vertices) == 0 {
			t.Error("网格节点没有顶点")
		}

		if len(node.FaceGroup) == 0 {
			t.Log("网格节点没有面")
		}
	}
}

// TestImportFileToMST 测试文件导入转换
func TestImportFileToMST(t *testing.T) {
	// 由于我们没有实际的3D文件，这里测试错误处理
	_, _, err := ImportFileToMST("nonexistent.obj", PostProcessTriangulate)
	if err == nil {
		t.Error("应该返回文件不存在的错误")
	}
}

// TestMaterialConversion 测试材质转换
func TestMaterialConversion(t *testing.T) {
	// 创建测试材质
	material := &Material{
		Properties: []*MaterialProperty{
			{
				name:     "$clr.diffuse",
				Data:     []byte{255, 0, 0},
				TypeInfo: MatPropTypeInfoBuffer,
			},
		},
	}

	converted := convertMaterial(material)
	if converted == nil {
		t.Fatal("材质转换结果为nil")
	}

	// 验证颜色
	color := converted.GetColor()
	if color[0] != 255 || color[1] != 0 || color[2] != 0 {
		t.Errorf("颜色转换错误，期望[255,0,0]，得到%v", color)
	}
}

// TestMeshConversion 测试网格转换
func TestMeshConversion(t *testing.T) {
	// 创建测试网格
	aiMesh := &Mesh{
		Vertices: []vec3.T{{0, 0, 0}, {1, 0, 0}, {0, 1, 0}},
		Normals:  []vec3.T{{0, 0, 1}, {0, 0, 1}, {0, 0, 1}},
		Faces: []Face{
			{Indices: []uint{0, 1, 2}},
		},
		MaterialIndex: 0,
	}

	converted := convertMesh(aiMesh)
	if converted == nil {
		t.Fatal("网格转换结果为nil")
	}

	// 验证顶点
	if len(converted.Vertices) != 3 {
		t.Errorf("顶点数量错误，期望3，得到%d", len(converted.Vertices))
	}

	// 验证面
	if len(converted.FaceGroup) != 1 {
		t.Errorf("面组数量错误，期望1，得到%d", len(converted.FaceGroup))
	}

	if len(converted.FaceGroup[0].Faces) != 1 {
		t.Errorf("面数量错误，期望1，得到%d", len(converted.FaceGroup[0].Faces))
	}
}

// TestComplexSceneConversion 测试复杂场景转换
func TestComplexSceneConversion(t *testing.T) {
	// 创建复杂测试场景
	scene := &Scene{
		RootNode: &Node{
			Name: "Root",
			Transformation: &mat4.T{
				{1, 0, 0, 0},
				{0, 1, 0, 0},
				{0, 0, 1, 0},
				{0, 0, 0, 1},
			},
			MeshIndicies: []uint{0},
		},
		Meshes: []*Mesh{
			{
				Vertices: []vec3.T{
					{0, 0, 0}, {1, 0, 0}, {0, 1, 0},
					{1, 1, 0}, {0, 0, 1}, {1, 0, 1},
				},
				Normals: []vec3.T{
					{0, 0, 1}, {0, 0, 1}, {0, 0, 1},
					{0, 0, 1}, {0, 0, 1}, {0, 0, 1},
				},
				Faces: []Face{
					{Indices: []uint{0, 1, 2}},
					{Indices: []uint{1, 3, 2}},
					{Indices: []uint{4, 5, 0}},
				},
				MaterialIndex: 0,
			},
		},
		Materials: []*Material{
			{
				Properties: []*MaterialProperty{
					{name: "$clr.diffuse", Data: []byte{255, 128, 64}},
				},
			},
		},
	}

	converted := AssimpToMSTConverter(scene)
	if converted == nil {
		t.Fatal("复杂场景转换失败")
	}

	// 验证场景结构
	if len(converted.Materials) != 1 {
		t.Errorf("材质数量错误，期望1，得到%d", len(converted.Materials))
	}

	if len(converted.Nodes) != 1 {
		t.Errorf("网格节点数量错误，期望1，得到%d", len(converted.Nodes))
	}

	if len(converted.Nodes[0].Vertices) != 6 {
		t.Errorf("顶点数量错误，期望6，得到%d", len(converted.Nodes[0].Vertices))
	}

	if len(converted.Nodes[0].FaceGroup) != 1 {
		t.Errorf("面组数量错误，期望1，得到%d", len(converted.Nodes[0].FaceGroup))
	}

	if len(converted.Nodes[0].FaceGroup[0].Faces) != 3 {
		t.Errorf("面数量错误，期望3，得到%d", len(converted.Nodes[0].FaceGroup[0].Faces))
	}
}

// TestEmptySceneConversion 测试空场景转换
func TestEmptySceneConversion(t *testing.T) {
	scene := &Scene{}

	converted := AssimpToMSTConverter(scene)
	if converted == nil {
		t.Fatal("空场景转换失败")
	}

	if len(converted.Materials) != 0 {
		t.Errorf("空场景不应有材质")
	}

	if len(converted.Nodes) != 0 {
		t.Errorf("空场景不应有网格节点")
	}
}

// TestRoundTripConversion 测试往返转换
func TestRoundTripConversion(t *testing.T) {
	// 创建测试场景
	scene := createTestScene()

	// 转换为MST
	mstMesh := AssimpToMSTConverter(scene)

	// 验证转换后的数据
	if mstMesh == nil {
		t.Fatal("转换失败")
	}

	// 验证版本
	if mstMesh.Version != mst.V4 {
		t.Errorf("版本错误，期望V4，得到%d", mstMesh.Version)
	}

	// 验证材质
	if len(mstMesh.Materials) == 0 {
		t.Error("没有材质")
	}

	// 验证网格节点
	if len(mstMesh.Nodes) == 0 {
		t.Error("没有网格节点")
	}

	// 验证材质颜色
	material := mstMesh.Materials[0]
	color := material.GetColor()
	if color[0] != 255 || color[1] != 128 || color[2] != 64 {
		t.Errorf("材质颜色错误，期望[255,128,64]，得到%v", color)
	}
}

// createTestScene 创建测试场景
func createTestScene() *Scene {
	return &Scene{
		RootNode: &Node{
			Name: "TestRoot",
			Transformation: &mat4.T{
				{1, 0, 0, 0},
				{0, 1, 0, 0},
				{0, 0, 1, 0},
				{0, 0, 0, 1},
			},
			MeshIndicies: []uint{0},
		},
		Meshes: []*Mesh{
			{
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
			},
		},
		Materials: []*Material{
			{
				Properties: []*MaterialProperty{
					{name: "$clr.diffuse", Data: []byte{255, 128, 64}},
				},
			},
		},
	}
}

// TestConverterIntegration 测试转换器集成
func TestConverterIntegration(t *testing.T) {
	// 创建临时目录用于测试
	tempDir := t.TempDir()

	// 测试文件路径
	testFile := filepath.Join(tempDir, "test_conversion.mst")

	// 创建测试场景
	scene := createTestScene()

	// 转换为MST
	mstMesh := AssimpToMSTConverter(scene)
	if mstMesh == nil {
		t.Fatal("转换失败")
	}

	// 写入文件
	file, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("创建文件失败: %v", err)
	}
	defer file.Close()

	// 序列化到文件
	var buf bytes.Buffer
	mst.MeshMarshal(&buf, mstMesh)

	if _, cerr := file.Write(buf.Bytes()); cerr != nil {
		t.Fatalf("写入文件失败: %v", cerr)
	}

	// 验证文件存在
	if _, cerr := os.Stat(testFile); os.IsNotExist(cerr) {
		t.Fatal("文件不存在")
	}

	// 从文件读取
	readFile, err := os.Open(testFile)
	if err != nil {
		t.Fatalf("打开文件失败: %v", err)
	}
	defer readFile.Close()

	// 反序列化
	readMesh := mst.MeshUnMarshal(readFile)
	if readMesh == nil {
		t.Fatal("读取网格失败")
	}

	// 验证读取的数据
	if len(readMesh.Materials) != len(mstMesh.Materials) {
		t.Errorf("材质数量不匹配")
	}

	if len(readMesh.Nodes) != len(mstMesh.Nodes) {
		t.Errorf("网格节点数量不匹配")
	}
}
