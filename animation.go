package assimp

/*
#cgo CFLAGS: -I ./lib

#include <stdlib.h>
#include <assimp/scene.h>
*/
import "C"
import (
	"unsafe"

	"github.com/flywave/go3d/quaternion"
	"github.com/flywave/go3d/vec3"
)

type VectorKey C.struct_aiVectorKey

func (v VectorKey) Time() float64 {
	return float64(v.mTime)
}

func (v VectorKey) Value() vec3.T {
	return vec3.T{float32(v.mValue.x), float32(v.mValue.y), float32(v.mValue.z)}
}

type QuatKey C.struct_aiQuatKey

func (v QuatKey) Time() float64 {
	return float64(v.mTime)
}

func (v QuatKey) Value() quaternion.T {
	return quaternion.T{float32(v.mValue.x), float32(v.mValue.y), float32(v.mValue.z), float32(v.mValue.w)}
}

type MeshKey C.struct_aiMeshKey

func (v MeshKey) Time() float64 {
	return float64(v.mTime)
}

func (v MeshKey) Value() int {
	return int(v.mValue)
}

type AnimBehaviour C.enum_aiAnimBehaviour

const (
	AnimBehaviour_Default  AnimBehaviour = C.aiAnimBehaviour_DEFAULT
	AnimBehaviour_Constant AnimBehaviour = C.aiAnimBehaviour_CONSTANT
	AnimBehaviour_Linear   AnimBehaviour = C.aiAnimBehaviour_LINEAR
	AnimBehaviour_Repeat   AnimBehaviour = C.aiAnimBehaviour_REPEAT
)

type NodeAnim C.struct_aiNodeAnim

func (v *NodeAnim) Name() string {
	return C.GoString(&v.mNodeName.data[0])
}

func (v *NodeAnim) NumPositionKeys() int {
	return int(v.mNumPositionKeys)
}

func (v *NodeAnim) PositionKeys() []VectorKey {
	if v.mNumPositionKeys > 0 && v.mPositionKeys != nil {
		return unsafe.Slice((*VectorKey)(unsafe.Pointer(v.mPositionKeys)), int(v.mNumPositionKeys))
	} else {
		return nil
	}
}

func (v *NodeAnim) NumRotationKeys() int {
	return int(v.mNumRotationKeys)
}

func (v *NodeAnim) RotationKeys() []QuatKey {
	if v.mNumRotationKeys > 0 && v.mRotationKeys != nil {
		return unsafe.Slice((*QuatKey)(unsafe.Pointer(v.mRotationKeys)), int(v.mNumRotationKeys))
	} else {
		return nil
	}
}

func (v *NodeAnim) NumScalingKeys() int {
	return int(v.mNumScalingKeys)
}

func (v *NodeAnim) ScalingKeys() []VectorKey {
	if v.mNumScalingKeys > 0 && v.mScalingKeys != nil {
		return unsafe.Slice((*VectorKey)(unsafe.Pointer(v.mScalingKeys)), int(v.mNumScalingKeys))
	} else {
		return nil
	}
}

func (v *NodeAnim) PreState() AnimBehaviour {
	return AnimBehaviour(v.mPreState)
}

func (v *NodeAnim) PostState() AnimBehaviour {
	return AnimBehaviour(v.mPostState)
}

type MeshAnim C.struct_aiMeshAnim

func (v *MeshAnim) Name() string {
	return C.GoString(&v.mName.data[0])
}

func (v *MeshAnim) NumKeys() int {
	return int(v.mNumKeys)
}

func (v *MeshAnim) Keys() []MeshKey {
	if v.mNumKeys > 0 && v.mKeys != nil {
		return unsafe.Slice((*MeshKey)(unsafe.Pointer(v.mKeys)), int(v.mNumKeys))
	} else {
		return nil
	}
}

type Animation struct {
	a *C.struct_aiAnimation
}

func (v *Animation) Name() string {
	return C.GoString(&v.a.mName.data[0])
}

func (v *Animation) Duration() float64 {
	return float64(v.a.mDuration)
}

func (v *Animation) TicksPerSecond() float64 {
	return float64(v.a.mTicksPerSecond)
}

func (v *Animation) NumChannels() int {
	return int(v.a.mNumChannels)
}

func (v *Animation) Channels() []*NodeAnim {
	if v.a.mNumChannels > 0 && v.a.mChannels != nil {
		channels := unsafe.Slice(v.a.mChannels, int(v.a.mNumChannels))
		result := make([]*NodeAnim, len(channels))
		for i := range channels {
			result[i] = (*NodeAnim)(unsafe.Pointer(&channels[i]))
		}
		return result
	} else {
		return nil
	}
}

func (v *Animation) NumMeshChannels() int {
	return int(v.a.mNumMeshChannels)
}

func (v *Animation) MeshChannels() []*MeshAnim {
	if v.a.mNumMeshChannels > 0 && v.a.mMeshChannels != nil {
		channels := unsafe.Slice(v.a.mMeshChannels, int(v.a.mNumMeshChannels))
		result := make([]*MeshAnim, len(channels))
		for i := range channels {
			result[i] = (*MeshAnim)(unsafe.Pointer(&channels[i]))
		}
		return result
	} else {
		return nil
	}
}
