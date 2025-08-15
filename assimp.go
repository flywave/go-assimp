package assimp

/*
#cgo linux CFLAGS: -I ./lib/linux
#cgo darwin,amd64 CFLAGS: -I ./lib/darwin
#cgo darwin,arm64 CFLAGS: -I ./lib/darwin_arm

#include <stdlib.h>
#include <assimp/scene.h>

//Functions
struct aiScene* aiImportFile(const char* pFile, unsigned int pFlags);
void aiReleaseImport(const struct aiScene* pScene);
const char* aiGetErrorString();
unsigned int aiGetMaterialTextureCount(const struct aiMaterial* pMat, enum aiTextureType type);
*/
import "C"
import (
	"errors"
	"unsafe"

	"github.com/flywave/go3d/mat4"
	"github.com/flywave/go3d/vec3"
	"github.com/flywave/go3d/vec4"
)

type Node struct {
	Name           string
	Transformation *mat4.T
	Parent         *Node
	Children       []*Node
	MeshIndicies   []uint
	Metadata       map[string]Metadata
}

type EmbeddedTexture struct {
	cTex         *C.struct_aiTexture
	Width        uint
	Height       uint
	FormatHint   string
	Data         []byte
	IsCompressed bool
	Filename     string
}

type Metadata struct {
	Type  MetadataType
	Value interface{}
}

type MetadataEntry struct {
	Data []byte
}

type Scene struct {
	cScene *C.struct_aiScene
	Flags  SceneFlag

	RootNode  *Node
	Meshes    []*Mesh
	Materials []*Material
	Textures  []*EmbeddedTexture

	Animations []*Animation
	Lights     []*Light
	Cameras    []*Camera
}

func (s *Scene) releaseCResources() {
	C.aiReleaseImport(s.cScene)
}

func ImportFile(file string, postProcessFlags PostProcess) (s *Scene, release func(), err error) {

	cstr := C.CString(file)
	defer C.free(unsafe.Pointer(cstr))

	cs := C.aiImportFile(cstr, C.uint(postProcessFlags))
	if cs == nil {
		return nil, func() {}, getAiErr()
	}

	s = parseScene(cs)
	return s, func() { s.releaseCResources() }, nil
}

func getAiErr() error {
	return errors.New("asig error: " + C.GoString(C.aiGetErrorString()))
}

func parseScene(cs *C.struct_aiScene) *Scene {

	s := &Scene{cScene: cs}
	s.Flags = SceneFlag(cs.mFlags)
	s.RootNode = parseRootNode(cs.mRootNode)
	s.Meshes = parseMeshes(cs.mMeshes, uint(cs.mNumMeshes))
	s.Materials = parseMaterials(cs.mMaterials, uint(cs.mNumMaterials))
	s.Textures = parseTextures(cs.mTextures, uint(s.cScene.mNumTextures))

	s.Animations = parseAnimations(cs.mAnimations, uint(cs.mNumAnimations))
	s.Lights = parseLights(cs.mLights, uint(s.cScene.mNumLights))
	s.Cameras = parseCameras(cs.mCameras, uint(s.cScene.mNumCameras))

	return s
}

func parseRootNode(cNodesIn *C.struct_aiNode) *Node {

	rn := &Node{
		Name:           parseAiString(cNodesIn.mName),
		Transformation: parseMat4(&cNodesIn.mTransformation),
		Parent:         nil,
		MeshIndicies:   parseUInts(cNodesIn.mMeshes, uint(cNodesIn.mNumMeshes)),
		Metadata:       parseMetadata(cNodesIn.mMetaData),
	}

	rn.Children = parseNodes(cNodesIn.mChildren, rn, int(cNodesIn.mNumChildren))
	return rn
}

func parseNodes(cNodesIn **C.struct_aiNode, parent *Node, parentChildrenCount int) []*Node {

	if cNodesIn == nil || parentChildrenCount <= 0 {
		return []*Node{}
	}

	nodes := make([]*Node, parentChildrenCount)

	cNodes := unsafe.Slice((*C.struct_aiNode)(unsafe.Pointer(cNodesIn)), parentChildrenCount)

	for i := 0; i < len(nodes); i++ {

		n := cNodes[i]

		nodes[i] = &Node{
			Name:           parseAiString(n.mName),
			Transformation: parseMat4(&n.mTransformation),
			Parent:         parent,
			MeshIndicies:   parseUInts(n.mMeshes, uint(n.mNumMeshes)),
			Metadata:       parseMetadata(n.mMetaData),
		}

		nodes[i].Children = parseNodes(n.mChildren, nodes[i], parentChildrenCount)
	}

	return nodes
}

func parseMetadata(cMetaIn *C.struct_aiMetadata) map[string]Metadata {

	if cMetaIn == nil {
		return map[string]Metadata{}
	}

	meta := make(map[string]Metadata, cMetaIn.mNumProperties)

	cKeys := unsafe.Slice((*C.struct_aiString)(unsafe.Pointer(cMetaIn.mKeys)), int(cMetaIn.mNumProperties))
	cVals := unsafe.Slice((*C.struct_aiMetadataEntry)(unsafe.Pointer(cMetaIn.mValues)), int(cMetaIn.mNumProperties))

	for i := 0; i < int(cMetaIn.mNumProperties); i++ {
		meta[parseAiString(cKeys[i])] = parseMetadataEntry(cVals[i])
	}

	return meta
}

func parseMetadataEntry(cv C.struct_aiMetadataEntry) Metadata {

	m := Metadata{Type: MetadataType(cv.mType)}

	if cv.mData == nil {
		return m
	}

	switch m.Type {
	case MetadataTypeBool:
		m.Value = *(*bool)(cv.mData)
	case MetadataTypeFloat32:
		m.Value = *(*float32)(cv.mData)
	case MetadataTypeFloat64:
		m.Value = *(*float64)(cv.mData)
	case MetadataTypeInt32:
		m.Value = *(*int32)(cv.mData)
	case MetadataTypeUint64:
		m.Value = *(*uint64)(cv.mData)
	case MetadataTypeString:
		m.Value = parseAiString(*(*C.struct_aiString)(cv.mData))
	case MetadataTypeVec3:
		m.Value = parseVec3((*C.struct_aiVector3D)(cv.mData))
	}

	return m
}

func parseTextures(cTexIn **C.struct_aiTexture, count uint) []*EmbeddedTexture {

	if cTexIn == nil {
		return []*EmbeddedTexture{}
	}

	textures := make([]*EmbeddedTexture, count)

	cTex := unsafe.Slice((*C.struct_aiTexture)(unsafe.Pointer(cTexIn)), int(count))

	for i := 0; i < int(count); i++ {

		textures[i] = &EmbeddedTexture{
			cTex:         &cTex[i],
			Width:        uint(cTex[i].mWidth),
			Height:       uint(cTex[i].mHeight),
			FormatHint:   C.GoString(&cTex[i].achFormatHint[0]),
			Filename:     parseAiString(cTex[i].mFilename),
			Data:         parseTexels(cTex[i].pcData, uint(cTex[i].mWidth), uint(cTex[i].mHeight)),
			IsCompressed: cTex[i].mHeight == 0,
		}
	}

	return textures
}

func parseAnimations(cAnim **C.struct_aiAnimation, count uint) []*Animation {

	if cAnim == nil {
		return []*Animation{}
	}

	animations := make([]*Animation, count)

	cAmis := unsafe.Slice((*C.struct_aiAnimation)(unsafe.Pointer(cAnim)), int(count))

	for i := 0; i < int(count); i++ {

		animations[i] = &Animation{
			a: &cAmis[i],
		}
	}

	return animations
}

func parseLights(cLight **C.struct_aiLight, count uint) []*Light {

	if cLight == nil {
		return []*Light{}
	}

	lights := make([]*Light, count)

	cLights := unsafe.Slice((*C.struct_aiLight)(unsafe.Pointer(cLight)), int(count))

	for i := 0; i < int(count); i++ {

		lights[i] = &Light{
			l: &cLights[i],
		}
	}

	return lights
}

func parseCameras(cCamera **C.struct_aiCamera, count uint) []*Camera {
	if cCamera == nil {
		return []*Camera{}
	}

	cams := make([]*Camera, count)

	cCameras := unsafe.Slice((*C.struct_aiCamera)(unsafe.Pointer(cCamera)), int(count))

	for i := 0; i < int(count); i++ {

		cams[i] = &Camera{
			c: &cCameras[i],
		}
	}

	return cams
}

func parseTexels(cTexelsIn *C.struct_aiTexel, width, height uint) []byte {

	isCompressed := height == 0

	texelCount := width
	if !isCompressed {
		texelCount *= height
	}
	texelCount /= 4

	data := make([]byte, texelCount*4)

	cTexels := unsafe.Slice((*C.struct_aiTexel)(unsafe.Pointer(cTexelsIn)), int(texelCount))

	for i := 0; i < int(texelCount); i++ {

		index := i * 4
		data[index] = byte(cTexels[i].b)
		data[index+1] = byte(cTexels[i].g)
		data[index+2] = byte(cTexels[i].r)
		data[index+3] = byte(cTexels[i].a)
	}

	return data
}

func parseMeshes(cm **C.struct_aiMesh, count uint) []*Mesh {
	if cm == nil {
		return []*Mesh{}
	}

	meshes := make([]*Mesh, count)

	cmeshes := unsafe.Slice((*C.struct_aiMesh)(unsafe.Pointer(cm)), int(count))

	for i := 0; i < int(count); i++ {

		m := &Mesh{}

		cmesh := cmeshes[i]
		vertCount := uint(cmesh.mNumVertices)

		m.Vertices = parseVec3s(cmesh.mVertices, vertCount)
		m.Normals = parseVec3s(cmesh.mNormals, vertCount)
		m.Tangents = parseVec3s(cmesh.mTangents, vertCount)
		m.BitTangents = parseVec3s(cmesh.mBitangents, vertCount)

		m.ColorSets = parseColorSet(cmesh.mColors, vertCount)

		m.TexCoords = parseTexCoords(cmesh.mTextureCoords, vertCount)
		m.TexCoordChannelCount = [8]uint{}
		for j := 0; j < len(cmesh.mTextureCoords); j++ {

			if cmesh.mTextureCoords[j] == nil {
				continue
			}

			m.TexCoordChannelCount[j] = uint(cmeshes[j].mNumUVComponents[j])
		}

		cFaces := unsafe.Slice((*C.struct_aiFace)(unsafe.Pointer(cmesh.mFaces)), int(cmesh.mNumFaces))

		m.Faces = make([]Face, cmesh.mNumFaces)
		for j := 0; j < len(m.Faces); j++ {

			m.Faces[j] = Face{
				Indices: parseUInts(cFaces[j].mIndices, uint(cFaces[j].mNumIndices)),
			}
		}

		m.Bones = parseBones(cmesh.mBones, uint(cmesh.mNumBones))
		m.AnimMeshes = parseAnimMeshes(cmesh.mAnimMeshes, uint(cmesh.mNumAnimMeshes))
		m.AABB = AABB{
			Min: parseVec3(&cmesh.mAABB.mMin),
			Max: parseVec3(&cmesh.mAABB.mMax),
		}

		m.MorphMethod = MorphMethod(cmesh.mMethod)
		m.MaterialIndex = uint(cmesh.mMaterialIndex)
		m.Name = parseAiString(cmesh.mName)

		meshes[i] = m
	}

	return meshes
}

func parseVec3(cv *C.struct_aiVector3D) vec3.T {
	if cv == nil {
		return vec3.T{}
	}

	return vec3.T{
		float32(cv.x),
		float32(cv.y),
		float32(cv.z),
	}
}

func parseAnimMeshes(cam **C.struct_aiAnimMesh, count uint) []*AnimMesh {
	if cam == nil {
		return []*AnimMesh{}
	}

	animMeshes := make([]*AnimMesh, count)

	cAnimMeshes := unsafe.Slice((*C.struct_aiAnimMesh)(unsafe.Pointer(cam)), int(count))

	for i := 0; i < int(count); i++ {

		m := cAnimMeshes[i]
		animMeshes[i] = &AnimMesh{
			Name:        parseAiString(m.mName),
			Vertices:    parseVec3s(m.mVertices, uint(m.mNumVertices)),
			Normals:     parseVec3s(m.mNormals, uint(m.mNumVertices)),
			Tangents:    parseVec3s(m.mTangents, uint(m.mNumVertices)),
			BitTangents: parseVec3s(m.mBitangents, uint(m.mNumVertices)),
			Colors:      parseColorSet(m.mColors, uint(m.mNumVertices)),
			TexCoords:   parseTexCoords(m.mTextureCoords, uint(m.mNumVertices)),
			Weight:      float32(m.mWeight),
		}
	}

	return animMeshes
}

func parseTexCoords(ctc [MaxTexCoords]*C.struct_aiVector3D, vertCount uint) [MaxTexCoords][]vec3.T {
	texCoords := [MaxTexCoords][]vec3.T{}

	for j := 0; j < len(ctc); j++ {

		if ctc[j] == nil {
			continue
		}

		texCoords[j] = parseVec3s(ctc[j], vertCount)
	}

	return texCoords
}

func parseColorSet(cc [MaxColorSets]*C.struct_aiColor4D, vertCount uint) [MaxColorSets][]vec4.T {
	colorSet := [MaxColorSets][]vec4.T{}

	for j := 0; j < len(cc); j++ {

		if cc[j] == nil {
			continue
		}

		colorSet[j] = parseColors(cc[j], vertCount)
	}

	return colorSet
}

func parseBones(cbs **C.struct_aiBone, count uint) []*Bone {
	if cbs == nil {
		return []*Bone{}
	}

	bones := make([]*Bone, count)

	cbones := unsafe.Slice((*C.struct_aiBone)(unsafe.Pointer(cbs)), int(count))

	for i := 0; i < int(count); i++ {

		cBone := cbones[i]
		bones[i] = &Bone{
			Name:         parseAiString(cBone.mName),
			Weights:      parseVertexWeights(cBone.mWeights, uint(cBone.mNumWeights)),
			OffsetMatrix: *parseMat4(&cBone.mOffsetMatrix),
		}
	}

	return bones
}

func parseMat4(cm4 *C.struct_aiMatrix4x4) *mat4.T {
	if cm4 == nil {
		return &mat4.T{}
	}

	return &mat4.T{
		{float32(cm4.a1), float32(cm4.b1), float32(cm4.c1), float32(cm4.d1)},
		{float32(cm4.a2), float32(cm4.b2), float32(cm4.c2), float32(cm4.d2)},
		{float32(cm4.a3), float32(cm4.b3), float32(cm4.c3), float32(cm4.d3)},
		{float32(cm4.a4), float32(cm4.b4), float32(cm4.c4), float32(cm4.d4)},
	}
}

func parseVertexWeights(cWeights *C.struct_aiVertexWeight, count uint) []VertexWeight {
	if cWeights == nil {
		return []VertexWeight{}
	}

	vw := make([]VertexWeight, count)

	cvw := unsafe.Slice((*C.struct_aiVertexWeight)(unsafe.Pointer(cWeights)), int(count))

	for i := 0; i < int(count); i++ {

		vw[i] = VertexWeight{
			VertIndex: uint(cvw[i].mVertexId),
			Weight:    float32(cvw[i].mWeight),
		}
	}

	return vw
}

func parseAiString(aiString C.struct_aiString) string {
	return C.GoStringN(&aiString.data[0], C.int(aiString.length))
}

func parseUInts(cui *C.uint, count uint) []uint {
	if cui == nil {
		return []uint{}
	}

	uints := make([]uint, count)

	cUInts := unsafe.Slice((*C.uint)(unsafe.Pointer(cui)), int(count))

	for i := 0; i < len(cUInts); i++ {
		uints[i] = uint(cUInts[i])
	}

	return uints
}

func parseVec3s(cv *C.struct_aiVector3D, count uint) []vec3.T {
	if cv == nil {
		return []vec3.T{}
	}

	carr := unsafe.Slice((*C.struct_aiVector3D)(unsafe.Pointer(cv)), int(count))

	verts := make([]vec3.T, count)

	for i := 0; i < int(count); i++ {
		verts[i] = vec3.T{
			float32(carr[i].x),
			float32(carr[i].y),
			float32(carr[i].z),
		}
	}

	return verts
}

func parseColors(cv *C.struct_aiColor4D, count uint) []vec4.T {
	if cv == nil {
		return []vec4.T{}
	}

	carr := unsafe.Slice((*C.struct_aiColor4D)(unsafe.Pointer(cv)), int(count))

	verts := make([]vec4.T, count)

	for i := 0; i < int(count); i++ {
		verts[i] = vec4.T{
			float32(carr[i].r),
			float32(carr[i].g),
			float32(carr[i].b),
			float32(carr[i].a),
		}
	}

	return verts
}

func parseMaterials(cMatsIn **C.struct_aiMaterial, count uint) []*Material {
	mats := make([]*Material, count)

	cMats := unsafe.Slice((*C.struct_aiMaterial)(unsafe.Pointer(cMatsIn)), int(count))

	for i := 0; i < int(count); i++ {

		mats[i] = &Material{
			cMat:             &cMats[i],
			Properties:       parseMatProperties(cMats[i].mProperties, uint(cMats[i].mNumProperties)),
			AllocatedStorage: uint(cMats[i].mNumAllocated),
		}
	}

	return mats
}

func parseMatProperties(cMatPropsIn **C.struct_aiMaterialProperty, count uint) []*MaterialProperty {
	matProps := make([]*MaterialProperty, count)

	cMatProps := unsafe.Slice((*C.struct_aiMaterialProperty)(unsafe.Pointer(cMatPropsIn)), int(count))

	for i := 0; i < int(count); i++ {

		cmp := cMatProps[i]

		matProps[i] = &MaterialProperty{
			name:     parseAiString(cmp.mKey),
			Semantic: TextureType(cmp.mSemantic),
			Index:    uint(cmp.mIndex),
			TypeInfo: MatPropertyTypeInfo(cmp.mType),
			Data:     C.GoBytes(unsafe.Pointer(cmp.mData), C.int(cmp.mDataLength)),
		}
	}

	return matProps
}
