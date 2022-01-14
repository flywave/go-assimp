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
