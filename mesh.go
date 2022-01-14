package assimp

import (
	"github.com/flywave/go3d/mat4"
	"github.com/flywave/go3d/vec3"
	"github.com/flywave/go3d/vec4"
)

const (
	MaxColorSets = 8
	MaxTexCoords = 8
)

type Mesh struct {
	PrimitiveTypes PrimitiveType
	Vertices       []vec3.T
	Normals        []vec3.T
	Tangents       []vec3.T
	BitTangents    []vec3.T

	ColorSets [MaxColorSets][]vec4.T

	TexCoords            [MaxTexCoords][]vec3.T
	TexCoordChannelCount [MaxTexCoords]uint

	Faces       []Face
	Bones       []*Bone
	AnimMeshes  []*AnimMesh
	AABB        AABB
	MorphMethod MorphMethod

	MaterialIndex uint
	Name          string
}

type Face struct {
	Indices []uint
}

type AnimMesh struct {
	Name string

	Vertices    []vec3.T
	Normals     []vec3.T
	Tangents    []vec3.T
	BitTangents []vec3.T
	Colors      [MaxColorSets][]vec4.T
	TexCoords   [MaxTexCoords][]vec3.T

	Weight float32
}

type AABB struct {
	Min vec3.T
	Max vec3.T
}

type Bone struct {
	Name string

	Weights []VertexWeight

	OffsetMatrix mat4.T
}

type VertexWeight struct {
	VertIndex uint
	Weight    float32
}
