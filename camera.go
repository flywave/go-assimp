package assimp

/*
#cgo CFLAGS: -I ./lib

#include <stdlib.h>
#include <assimp/scene.h>
*/
import "C"
import "github.com/flywave/go3d/vec3"

type Camera struct {
	c *C.struct_aiCamera
}

func (cam *Camera) Name() string {
	return C.GoStringN(&cam.c.mName.data[0], C.int(cam.c.mName.length))
}

func (cam *Camera) Position() vec3.T {
	return vec3.T{float32(cam.c.mPosition.x), float32(cam.c.mPosition.y), float32(cam.c.mPosition.z)}
}

func (cam *Camera) Up() vec3.T {
	return vec3.T{float32(cam.c.mUp.x), float32(cam.c.mUp.y), float32(cam.c.mUp.z)}
}

func (cam *Camera) LookAt() vec3.T {
	return vec3.T{float32(cam.c.mLookAt.x), float32(cam.c.mLookAt.y), float32(cam.c.mLookAt.z)}
}

func (cam *Camera) HorizontalFov() float32 {
	return float32(cam.c.mHorizontalFOV)
}

func (cam *Camera) ClipPlaneNear() float32 {
	return float32(cam.c.mClipPlaneNear)
}

func (cam *Camera) ClipPlaneFar() float32 {
	return float32(cam.c.mClipPlaneFar)
}

func (cam *Camera) Aspect() float32 {
	return float32(cam.c.mAspect)
}
