package assimp

/*
#cgo CFLAGS: -I ./lib

#include <stdlib.h>
#include <assimp/scene.h>
*/
import "C"
import (
	"reflect"
	"unsafe"

	"github.com/flywave/go3d/quaternion"
	"github.com/flywave/go3d/vec3"
)

type VectorKey C.struct_aiVectorKey

func (this VectorKey) Time() float64 {
	return float64(this.mTime)
}

func (this VectorKey) Value() vec3.T {
	return vec3.T{float32(this.mValue.x), float32(this.mValue.y), float32(this.mValue.z)}
}

type QuatKey C.struct_aiQuatKey

func (this QuatKey) Time() float64 {
	return float64(this.mTime)
}

func (this QuatKey) Value() quaternion.T {
	return quaternion.T{float32(this.mValue.x), float32(this.mValue.y), float32(this.mValue.z), float32(this.mValue.w)}
}

type MeshKey C.struct_aiMeshKey

func (this MeshKey) Time() float64 {
	return float64(this.mTime)
}

func (this MeshKey) Value() int {
	return int(this.mValue)
}

type AnimBehaviour C.enum_aiAnimBehaviour

const (
	AnimBehaviour_Default  AnimBehaviour = C.aiAnimBehaviour_DEFAULT
	AnimBehaviour_Constant AnimBehaviour = C.aiAnimBehaviour_CONSTANT
	AnimBehaviour_Linear   AnimBehaviour = C.aiAnimBehaviour_LINEAR
	AnimBehaviour_Repeat   AnimBehaviour = C.aiAnimBehaviour_REPEAT
)

type NodeAnim C.struct_aiNodeAnim

func (this *NodeAnim) Name() string {
	return C.GoString(&this.mNodeName.data[0])
}

func (this *NodeAnim) NumPositionKeys() int {
	return int(this.mNumPositionKeys)
}

func (this *NodeAnim) PositionKeys() []VectorKey {
	if this.mNumPositionKeys > 0 && this.mPositionKeys != nil {
		var result []VectorKey
		header := (*reflect.SliceHeader)(unsafe.Pointer(&result))
		header.Cap = int(this.mNumPositionKeys)
		header.Len = int(this.mNumPositionKeys)
		header.Data = uintptr(unsafe.Pointer(this.mPositionKeys))
		return result
	} else {
		return nil
	}
}

func (this *NodeAnim) NumRotationKeys() int {
	return int(this.mNumRotationKeys)
}

func (this *NodeAnim) RotationKeys() []QuatKey {
	if this.mNumRotationKeys > 0 && this.mRotationKeys != nil {
		var result []QuatKey
		header := (*reflect.SliceHeader)(unsafe.Pointer(&result))
		header.Cap = int(this.mNumRotationKeys)
		header.Len = int(this.mNumRotationKeys)
		header.Data = uintptr(unsafe.Pointer(this.mRotationKeys))
		return result
	} else {
		return nil
	}
}

func (this *NodeAnim) NumScalingKeys() int {
	return int(this.mNumScalingKeys)
}

func (this *NodeAnim) ScalingKeys() []VectorKey {
	if this.mNumScalingKeys > 0 && this.mScalingKeys != nil {
		var result []VectorKey
		header := (*reflect.SliceHeader)(unsafe.Pointer(&result))
		header.Cap = int(this.mNumScalingKeys)
		header.Len = int(this.mNumScalingKeys)
		header.Data = uintptr(unsafe.Pointer(this.mScalingKeys))
		return result
	} else {
		return nil
	}
}

func (this *NodeAnim) PreState() AnimBehaviour {
	return AnimBehaviour(this.mPreState)
}

func (this *NodeAnim) PostState() AnimBehaviour {
	return AnimBehaviour(this.mPostState)
}

type MeshAnim C.struct_aiMeshAnim

func (this *MeshAnim) Name() string {
	return C.GoString(&this.mName.data[0])
}

func (this *MeshAnim) NumKeys() int {
	return int(this.mNumKeys)
}

func (this *MeshAnim) Keys() []MeshKey {
	if this.mNumKeys > 0 && this.mKeys != nil {
		var result []MeshKey
		header := (*reflect.SliceHeader)(unsafe.Pointer(&result))
		header.Cap = int(this.mNumKeys)
		header.Len = int(this.mNumKeys)
		header.Data = uintptr(unsafe.Pointer(this.mKeys))
		return result
	} else {
		return nil
	}
}

type Animation struct {
	a *C.struct_aiAnimation
}

func (this *Animation) Name() string {
	return C.GoString(&this.a.mName.data[0])
}

func (this *Animation) Duration() float64 {
	return float64(this.a.mDuration)
}

func (this *Animation) TicksPerSecond() float64 {
	return float64(this.a.mTicksPerSecond)
}

func (this *Animation) NumChannels() int {
	return int(this.a.mNumChannels)
}

func (this *Animation) Channels() []*NodeAnim {
	if this.a.mNumChannels > 0 && this.a.mChannels != nil {
		var result []*NodeAnim
		header := (*reflect.SliceHeader)(unsafe.Pointer(&result))
		header.Cap = int(this.a.mNumChannels)
		header.Len = int(this.a.mNumChannels)
		header.Data = uintptr(unsafe.Pointer(this.a.mChannels))
		return result
	} else {
		return nil
	}
}

func (this *Animation) NumMeshChannels() int {
	return int(this.a.mNumMeshChannels)
}

func (this *Animation) MeshChannels() []*MeshAnim {
	if this.a.mNumMeshChannels > 0 && this.a.mMeshChannels != nil {
		var result []*MeshAnim
		header := (*reflect.SliceHeader)(unsafe.Pointer(&result))
		header.Cap = int(this.a.mNumMeshChannels)
		header.Len = int(this.a.mNumMeshChannels)
		header.Data = uintptr(unsafe.Pointer(this.a.mMeshChannels))
		return result
	} else {
		return nil
	}
}
