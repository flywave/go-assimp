package assimp

/*
#cgo CFLAGS: -I ./lib
#include <stdlib.h> //Needed for C.free

#include <assimp/scene.h>

//Functions
unsigned int aiGetMaterialTextureCount(const struct aiMaterial* pMat, enum aiTextureType type);

enum aiReturn aiGetMaterialTexture(
	const struct aiMaterial* mat,
    enum aiTextureType type,
    unsigned int  index,
    struct aiString* path,
    enum aiTextureMapping* mapping,
    unsigned int* uvindex,
    ai_real* blend,
    enum aiTextureOp* op,
    enum aiTextureMapMode* mapmode,
    unsigned int* flags);
*/
import "C"
import (
	"errors"
	"fmt"
)

type Material struct {
	cMat             *C.struct_aiMaterial
	Properties       []*MaterialProperty
	AllocatedStorage uint
}

type MaterialProperty struct {
	name     string
	Semantic TextureType
	Index    uint
	TypeInfo MatPropertyTypeInfo
	Data     []byte
}

func GetMaterialTextureCount(m *Material, texType TextureType) int {
	return int(C.aiGetMaterialTextureCount(m.cMat, uint32(texType)))
}

type GetMatTexInfo struct {
	Path string
}

func GetMaterialTexture(m *Material, texType TextureType, texIndex uint) (*GetMatTexInfo, error) {
	outCPath := &C.struct_aiString{}
	status := aiReturn(C.aiGetMaterialTexture(m.cMat, uint32(texType), C.uint(texIndex), outCPath, nil, nil, nil, nil, nil, nil))
	if status == aiReturnSuccess {
		return &GetMatTexInfo{
			Path: parseAiString(*outCPath),
		}, nil
	}

	if status == aiReturnFailure {
		return nil, errors.New("get texture failed: " + getAiErr().Error())
	}

	if status == aiReturnOutofMemory {
		return nil, errors.New("get texture failed: out of memory")
	}

	return nil, errors.New("get texture failed: unknown error with code " + fmt.Sprintf("%v", status))
}

// Helper to convert ASSIMP property names to better format
func NicePropName(prop string) string {
	switch prop {
	case "?mat.name":
		return "NAME"
	case "$mat.twosided":
		return "TWOSIDED"
	case "$mat.shadingm":
		return "SHADING_MODEL"
	case "$mat.wireframe":
		return "ENABLE_WIREFRAME"
	case "$mat.blend":
		return "BLEND_FUNC"
	case "$mat.opacity":
		return "OPACITY"
	case "$mat.bumpscaling":
		return "BUMPSCALING"
	case "$mat.shininess":
		return "SHININESS"
	case "$mat.reflectivity":
		return "REFLECTIVITY"
	case "$mat.shinpercent":
		return "SHININESS_STRENGTH"
	case "$mat.refracti":
		return "REFRACTI"
	case "$clr.diffuse":
		return "COLOR_DIFFUSE"
	case "$clr.ambient":
		return "COLOR_AMBIENT"
	case "$clr.specular":
		return "COLOR_SPECULAR"
	case "$clr.emissive":
		return "COLOR_EMISSIVE"
	case "$clr.transparent":
		return "COLOR_TRANSPARENT"
	case "$clr.reflective":
		return "COLOR_REFLECTIVE"
	case "?bg.global":
		return "GLOBAL_BACKGROUND_IMAGE"
	case "$tex.file":
		return "TEXTURE_BASE"
	case "$tex.mapping":
		return "MAPPING_BASE"
	case "$tex.flags":
		return "TEXFLAGS_BASE"
	case "$tex.uvwsrc":
		return "UVWSRC_BASE"
	case "$tex.mapmodev":
		return "MAPPINGMODE_V_BASE"
	case "$tex.mapaxis":
		return "TEXMAP_AXIS_BASE"
	case "$tex.blend":
		return "TEXBLEND_BASE"
	case "$tex.uvtrafo":
		return "UVTRANSFORM_BASE"
	case "$tex.op":
		return "TEXOP_BASE"
	case "$tex.mapmodeu":
		return "MAPPINGMODE_U_BASE"
	default:
		return "NONE"
	}
}

// GetNiceName returns the property name in a more readable format
func (mp *MaterialProperty) GetNiceName() string {
	return NicePropName(mp.name)
}
