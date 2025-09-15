package assimp

import (
	"errors"
	"hash/fnv"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"

	mat4d "github.com/flywave/go3d/float64/mat4"

	"github.com/chai2010/tiff"
	"github.com/flywave/go-mst"
	"github.com/flywave/go3d/mat4"
	"github.com/flywave/go3d/vec2"
	"github.com/flywave/go3d/vec3"
	"golang.org/x/image/bmp"
)

// AssimpToMSTConverter 将assimp场景转换为mst网格
func AssimpToMSTConverter(scene *Scene) *mst.Mesh {
	if scene == nil {
		return nil
	}

	mesh := mst.NewMesh()

	// 转换材质
	for _, material := range scene.Materials {
		mstMaterial := convertMaterial(material)
		mesh.Materials = append(mesh.Materials, mstMaterial)
	}

	// 转换基础网格
	baseMesh := &mst.BaseMesh{
		Materials: mesh.Materials,
		Nodes:     make([]*mst.MeshNode, 0),
	}

	// 转换网格节点
	for _, aiMesh := range scene.Meshes {
		node := convertMesh(aiMesh)
		if node != nil {
			baseMesh.Nodes = append(baseMesh.Nodes, node)
		}
	}

	mesh.BaseMesh = *baseMesh

	// 处理场景节点层次结构和实例化
	if scene.RootNode != nil {
		processNodeHierarchy(scene.RootNode, scene, mesh)
	}

	return mesh
}

// convertMaterial 转换assimp材质到mst材质
func convertMaterial(material *Material) mst.MeshMaterial {
	if material == nil {
		return &mst.BaseMaterial{Color: [3]byte{128, 128, 128}}
	}

	// 材质属性映射
	color := [3]byte{128, 128, 128}
	diffuse := [3]byte{128, 128, 128}
	specular := [3]byte{0, 0, 0}
	ambient := [3]byte{0, 0, 0}
	emissive := [3]byte{0, 0, 0}
	transparency := float32(1.0)
	metallic := float32(0.0)
	roughness := float32(0.5)
	shininess := float32(32.0)

	// 提取材质属性
	for _, prop := range material.Properties {
		niceName := prop.GetNiceName()
		switch niceName {
		case "COLOR_DIFFUSE":
			if len(prop.Data) >= 3 {
				color = [3]byte{prop.Data[0], prop.Data[1], prop.Data[2]}
				diffuse = color
			}
		case "COLOR_SPECULAR":
			if len(prop.Data) >= 3 {
				specular = [3]byte{prop.Data[0], prop.Data[1], prop.Data[2]}
			}
		case "COLOR_AMBIENT":
			if len(prop.Data) >= 3 {
				ambient = [3]byte{prop.Data[0], prop.Data[1], prop.Data[2]}
			}
		case "COLOR_EMISSIVE":
			if len(prop.Data) >= 3 {
				emissive = [3]byte{prop.Data[0], prop.Data[1], prop.Data[2]}
			}
		case "OPACITY":
			if len(prop.Data) >= 4 {
				transparency = math.Float32frombits(binaryToUint32(prop.Data))
			}
		case "SHININESS":
			if len(prop.Data) >= 4 {
				shininess = math.Float32frombits(binaryToUint32(prop.Data))
			}
		case "METALLIC":
			if len(prop.Data) >= 4 {
				metallic = math.Float32frombits(binaryToUint32(prop.Data))
			}
		case "ROUGHNESS":
			if len(prop.Data) >= 4 {
				roughness = math.Float32frombits(binaryToUint32(prop.Data))
			}
		}
	}

	// 检查纹理
	hasTexture := false
	textureMaterial := &mst.TextureMaterial{
		BaseMaterial: mst.BaseMaterial{Color: color, Transparency: transparency},
	}

	// 处理漫反射纹理
	diffuseTexture := getTextureFromMaterial(material, TextureTypeDiffuse)
	if diffuseTexture != nil {
		textureMaterial.Texture = diffuseTexture
		hasTexture = true
	}

	// 处理法线纹理
	normalTexture := getTextureFromMaterial(material, TextureTypeNormal)
	if normalTexture != nil {
		textureMaterial.Normal = normalTexture
		hasTexture = true
	}

	// 如果有纹理，优先返回TextureMaterial
	if hasTexture {
		return textureMaterial
	}

	// 根据材质类型创建合适的MST材质
	if hasPBRProperties(material) {
		return &mst.PbrMaterial{
			TextureMaterial: *textureMaterial,
			Emissive:        emissive,
			Metallic:        metallic,
			Roughness:       roughness,
			Reflectance:     0.5,
		}
	} else if hasSpecularProperties(material) {
		return &mst.PhongMaterial{
			LambertMaterial: mst.LambertMaterial{
				TextureMaterial: *textureMaterial,
				Ambient:         ambient,
				Diffuse:         diffuse,
				Emissive:        emissive,
			},
			Specular:    specular,
			Shininess:   shininess,
			Specularity: 1.0,
		}
	} else {
		return &mst.LambertMaterial{
			TextureMaterial: *textureMaterial,
			Ambient:         ambient,
			Diffuse:         diffuse,
			Emissive:        emissive,
		}
	}
}

// hasPBRProperties 检查材质是否有PBR属性
func hasPBRProperties(material *Material) bool {
	for _, prop := range material.Properties {
		if prop.GetNiceName() == "METALLIC" || prop.GetNiceName() == "ROUGHNESS" {
			return true
		}
	}
	return false
}

// hasSpecularProperties 检查材质是否有镜面反射属性
func hasSpecularProperties(material *Material) bool {
	for _, prop := range material.Properties {
		if prop.GetNiceName() == "COLOR_SPECULAR" {
			return true
		}
	}
	return false
}

// getTextureFromMaterial 从assimp材质中提取纹理并加载实际纹理数据
func getTextureFromMaterial(material *Material, textureType TextureType) *mst.Texture {
	for _, prop := range material.Properties {
		if prop.Semantic == textureType && prop.GetNiceName() == "TEXTURE_BASE" {
			if len(prop.Data) == 0 {
				continue
			}

			texturePath := string(prop.Data)
			if texturePath == "" {
				continue
			}

			// 清理纹理路径
			texturePath = strings.TrimSpace(texturePath)
			texturePath = strings.Trim(texturePath, "\x00")

			// 检查文件是否存在
			if _, err := os.Stat(texturePath); os.IsNotExist(err) {
				// 尝试相对路径
				filename := filepath.Base(texturePath)

				// 尝试不同路径组合
				possiblePaths := []string{
					filename,
					filepath.Join(".", filename),
					filepath.Join("..", filename),
					filepath.Join("textures", filename),
					filepath.Join("Textures", filename),
					filepath.Join("..", "textures", filename),
					filepath.Join("..", "Textures", filename),
				}

				found := false
				for _, tryPath := range possiblePaths {
					if _, err := os.Stat(tryPath); err == nil {
						texturePath = tryPath
						found = true
						break
					}
				}

				if !found {
					continue
				}
			}

			// 使用convertTex逻辑加载纹理
			return loadTextureFromPath(texturePath)
		}
	}
	return nil
}

// loadTextureFromPath 从文件路径加载纹理数据
func loadTextureFromPath(path string) *mst.Texture {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	// 获取文件类型
	_, ft, err := image.DecodeConfig(f)
	if err != nil {
		f.Seek(0, 0)
		// 根据文件扩展名判断类型
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".jpg", ".jpeg":
			ft = "jpeg"
		case ".png":
			ft = "png"
		case ".gif":
			ft = "gif"
		case ".bmp":
			ft = "bmp"
		case ".tiff", ".tif":
			ft = "tiff"
		default:
			return nil
		}
	}

	// 重置文件指针
	f.Seek(0, 0)

	// 解码图像
	img, err := readImage(f, ft)
	if err != nil {
		return nil
	}

	// 获取图像尺寸
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 转换像素数据
	buf := []byte{}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			cl := img.At(x, y)
			r, g, b, a := color.RGBAModel.Convert(cl).RGBA()
			buf = append(buf, byte(r&0xff), byte(g&0xff), byte(b&0xff), byte(a&0xff))
		}
	}

	// 创建纹理对象
	texture := &mst.Texture{
		Id:         int32(hashString(path)),
		Name:       filepath.Base(path),
		Size:       [2]uint64{uint64(width), uint64(height)},
		Format:     mst.TEXTURE_FORMAT_RGBA,
		Type:       mst.TEXTURE_PIXEL_TYPE_UBYTE,
		Compressed: mst.TEXTURE_COMPRESSED_ZLIB,
		Data:       mst.CompressImage(buf),
	}

	return texture
}

// readImage 读取不同格式的图像文件
func readImage(rd io.Reader, ft string) (image.Image, error) {
	switch ft {
	case "jpeg", "jpg":
		return jpeg.Decode(rd)
	case "png":
		return png.Decode(rd)
	case "gif":
		return gif.Decode(rd)
	case "bmp":
		return bmp.Decode(rd)
	case "tiff", "tif":
		return tiff.Decode(rd)
	default:
		return nil, errors.New("unsupported image format")
	}
}

// hashString 生成字符串的哈希值用于纹理ID
func hashString(s string) int32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int32(h.Sum32())
}

// convertMesh 转换assimp网格到mst网格节点
func convertMesh(aiMesh *Mesh) *mst.MeshNode {
	if aiMesh == nil {
		return nil
	}

	node := &mst.MeshNode{
		Vertices:  make([]vec3.T, 0, len(aiMesh.Vertices)),
		Normals:   make([]vec3.T, 0, len(aiMesh.Normals)),
		Colors:    make([][3]byte, 0),
		TexCoords: make([]vec2.T, 0),
	}

	// 转换顶点
	for _, v := range aiMesh.Vertices {
		node.Vertices = append(node.Vertices, vec3.T{v[0], v[1], v[2]})
	}

	// 转换法线
	for _, n := range aiMesh.Normals {
		node.Normals = append(node.Normals, vec3.T{n[0], n[1], n[2]})
	}

	// 转换颜色
	if len(aiMesh.ColorSets) > 0 && len(aiMesh.ColorSets[0]) > 0 {
		node.Colors = make([][3]byte, len(aiMesh.ColorSets[0]))
		for i, c := range aiMesh.ColorSets[0] {
			node.Colors[i] = [3]byte{
				byte(math.Max(0, math.Min(255, float64(c[0]*255)))),
				byte(math.Max(0, math.Min(255, float64(c[1]*255)))),
				byte(math.Max(0, math.Min(255, float64(c[2]*255)))),
			}
		}
	}

	// 转换纹理坐标
	if len(aiMesh.TexCoords) > 0 && len(aiMesh.TexCoords[0]) > 0 {
		node.TexCoords = make([]vec2.T, len(aiMesh.TexCoords[0]))
		for i, tc := range aiMesh.TexCoords[0] {
			node.TexCoords[i] = vec2.T{tc[0], tc[1]}
		}
	}

	// 转换面
	if len(aiMesh.Faces) > 0 {
		faceGroup := &mst.MeshTriangle{
			Batchid: int32(aiMesh.MaterialIndex),
			Faces:   make([]*mst.Face, 0, len(aiMesh.Faces)),
		}

		for _, face := range aiMesh.Faces {
			if len(face.Indices) >= 3 {
				// 处理三角面
				for i := 0; i < len(face.Indices)-2; i++ {
					mstFace := &mst.Face{
						Vertex: [3]uint32{
							uint32(face.Indices[0]),
							uint32(face.Indices[i+1]),
							uint32(face.Indices[i+2]),
						},
					}
					faceGroup.Faces = append(faceGroup.Faces, mstFace)
				}
			}
		}

		if len(faceGroup.Faces) > 0 {
			node.FaceGroup = []*mst.MeshTriangle{faceGroup}
		}
	}

	return node
}

// processNodeHierarchy 处理场景节点层次结构
func processNodeHierarchy(node *Node, scene *Scene, mesh *mst.Mesh) {
	if node == nil {
		return
	}

	// 转换变换矩阵
	transform := convertMatrix(*node.Transformation)

	// 处理当前节点的网格实例
	for _, meshIndex := range node.MeshIndicies {
		if int(meshIndex) < len(mesh.Nodes) {
			// 创建实例化网格
			instance := &mst.InstanceMesh{
				Transfors: []*mat4d.T{transform},
				Mesh:      &mst.BaseMesh{Nodes: []*mst.MeshNode{mesh.Nodes[meshIndex]}},
			}
			mesh.Instances = append(mesh.Instances, instance)
		}
	}

	// 递归处理子节点
	for _, child := range node.Children {
		processNodeHierarchy(child, scene, mesh)
	}
}

// convertMatrix 将assimp矩阵转换为go3d矩阵
func convertMatrix(aiMat mat4.T) *mat4d.T {
	return &mat4d.T{
		{float64(aiMat[0][0]), float64(aiMat[0][1]), float64(aiMat[0][2]), float64(aiMat[0][3])},
		{float64(aiMat[1][0]), float64(aiMat[1][1]), float64(aiMat[1][2]), float64(aiMat[1][3])},
		{float64(aiMat[2][0]), float64(aiMat[2][1]), float64(aiMat[2][2]), float64(aiMat[2][3])},
		{float64(aiMat[3][0]), float64(aiMat[3][1]), float64(aiMat[3][2]), float64(aiMat[3][3])},
	}
}

// binaryToUint32 将字节数组转换为uint32
func binaryToUint32(data []byte) uint32 {
	if len(data) < 4 {
		return 0
	}
	return uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16 | uint32(data[3])<<24
}

// ImportFileToMST 从文件导入并转换为MST格式
func ImportFileToMST(file string, postProcessFlags PostProcess) (*mst.Mesh, func(), error) {
	scene, release, err := ImportFile(file, postProcessFlags)
	if err != nil {
		return nil, nil, err
	}

	mstMesh := AssimpToMSTConverter(scene)
	return mstMesh, release, nil
}
