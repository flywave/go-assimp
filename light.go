package assimp

/*
#cgo CFLAGS: -I ./lib

#include <stdlib.h>
#include <assimp/scene.h>
*/
import "C"
import "github.com/flywave/go3d/vec3"

type LightSourceType C.enum_aiLightSourceType

const (
	LightSource_Undefined   LightSourceType = C.aiLightSource_UNDEFINED
	LightSource_Directional LightSourceType = C.aiLightSource_DIRECTIONAL
	LightSource_Point       LightSourceType = C.aiLightSource_POINT
	LightSource_Spot        LightSourceType = C.aiLightSource_SPOT
)

type Light struct {
	l *C.struct_aiLight
}

func (l *Light) Name() string {
	return C.GoStringN(&l.l.mName.data[0], C.int(l.l.mName.length))
}

func (l *Light) Type() LightSourceType {
	return LightSourceType(l.l.mType)
}

func (l *Light) Position() vec3.T {
	return vec3.T{float32(l.l.mPosition.x), float32(l.l.mPosition.y), float32(l.l.mPosition.z)}
}

func (l *Light) Direction() vec3.T {
	return vec3.T{float32(l.l.mDirection.x), float32(l.l.mDirection.y), float32(l.l.mDirection.z)}
}

func (l *Light) AttenuationConstant() float32 {
	return float32(l.l.mAttenuationConstant)
}

func (l *Light) AttenuationLinear() float32 {
	return float32(l.l.mAttenuationLinear)
}

func (l *Light) AttenuationQuadratic() float32 {
	return float32(l.l.mAttenuationQuadratic)
}

func (l *Light) ColorDiffuse() vec3.T {
	return vec3.T{float32(l.l.mColorDiffuse.r), float32(l.l.mColorDiffuse.g), float32(l.l.mColorDiffuse.b)}
}

func (l *Light) ColorSpecular() vec3.T {
	return vec3.T{float32(l.l.mColorSpecular.r), float32(l.l.mColorSpecular.g), float32(l.l.mColorSpecular.b)}
}

func (l *Light) ColorAmbient() vec3.T {
	return vec3.T{float32(l.l.mColorAmbient.r), float32(l.l.mColorAmbient.g), float32(l.l.mColorAmbient.b)}
}

func (l *Light) AngleInnerCone() float32 {
	return float32(l.l.mAngleInnerCone)
}

func (l *Light) AngleOuterCone() float32 {
	return float32(l.l.mAngleOuterCone)
}
