package assimp

type aiReturn int32

const (
	aiReturnSuccess     = 0x0
	aiReturnFailure     = -0x1
	aiReturnOutofMemory = -0x3
)

type SceneFlag int32

const (
	SceneFlagIncomplete        SceneFlag = 1 << 0
	SceneFlagValidated         SceneFlag = 1 << 1
	SceneFlagValidationWarning SceneFlag = 1 << 2
	SceneFlagNonVerboseFormat  SceneFlag = 1 << 3
	SceneFlagTerrain           SceneFlag = 1 << 4
	SceneFlagAllowShared       SceneFlag = 1 << 5
)

type PrimitiveType int32

const (
	PrimitiveTypePoint    = 1 << 0
	PrimitiveTypeLine     = 1 << 1
	PrimitiveTypeTriangle = 1 << 2
	PrimitiveTypePolygon  = 1 << 3
)

type MorphMethod int32

const (
	MorphMethodVertexBlend     = 0x1
	MorphMethodMorphNormalized = 0x2
	MorphMethodMorphRelative   = 0x3
)

type PostProcess int64

const (
	PostProcessCalcTangentSpace         PostProcess = 0x1
	PostProcessJoinIdenticalVertices    PostProcess = 0x2
	PostProcessMakeLeftHanded           PostProcess = 0x4
	PostProcessTriangulate              PostProcess = 0x8
	PostProcessRemoveComponent          PostProcess = 0x10
	PostProcessGenNormals               PostProcess = 0x20
	PostProcessGenSmoothNormals         PostProcess = 0x40
	PostProcessSplitLargeMeshes         PostProcess = 0x80
	PostProcessPreTransformVertices     PostProcess = 0x100
	PostProcessLimitBoneWeights         PostProcess = 0x200
	PostProcessValidateDataStructure    PostProcess = 0x400
	PostProcessImproveCacheLocality     PostProcess = 0x800
	PostProcessRemoveRedundantMaterials PostProcess = 0x1000
	PostProcessFixInfacingNormals       PostProcess = 0x2000
	PostProcessSortByPType              PostProcess = 0x8000
	PostProcessFindDegenerates          PostProcess = 0x10000
	PostProcessFindInvalidData          PostProcess = 0x20000
	PostProcessGenUVCoords              PostProcess = 0x40000
	PostProcessTransformUVCoords        PostProcess = 0x80000
	PostProcessFindInstances            PostProcess = 0x100000
	PostProcessOptimizeMeshes           PostProcess = 0x200000
	PostProcessOptimizeGraph            PostProcess = 0x400000
	PostProcessFlipUVs                  PostProcess = 0x800000
	PostProcessFlipWindingOrder         PostProcess = 0x1000000
	PostProcessSplitByBoneCount         PostProcess = 0x2000000
	PostProcessDebone                   PostProcess = 0x4000000
	PostProcessGlobalScale              PostProcess = 0x8000000
	PostProcessEmbedTextures            PostProcess = 0x10000000
	PostProcessForceGenNormals          PostProcess = 0x20000000
	PostProcessDropNormals              PostProcess = 0x40000000
	PostProcessGenBoundingBoxes         PostProcess = 0x80000000
)

type TextureType int32

const (
	TextureTypeNone         TextureType = 0
	TextureTypeDiffuse      TextureType = 1
	TextureTypeSpecular     TextureType = 2
	TextureTypeAmbient      TextureType = 3
	TextureTypeEmissive     TextureType = 4
	TextureTypeHeight       TextureType = 5
	TextureTypeNormal       TextureType = 6
	TextureTypeShininess    TextureType = 7
	TextureTypeOpacity      TextureType = 8
	TextureTypeDisplacement TextureType = 9
	TextureTypeLightmap     TextureType = 10
	TextureTypeReflection   TextureType = 11
	TextureTypeUnknown      TextureType = 18
)

const (
	TextureTypeBaseColor        TextureType = 12
	TextureTypeNormalCamera     TextureType = 13
	TextureTypeEmissionColor    TextureType = 14
	TextureTypeMetalness        TextureType = 15
	TextureTypeDiffuseRoughness TextureType = 16
	TextureTypeAmbientOcclusion TextureType = 17
)

func (tp TextureType) String() string {

	switch tp {
	case TextureTypeNone:
		return "None"
	case TextureTypeDiffuse:
		return "Diffuse"
	case TextureTypeSpecular:
		return "Specular"
	case TextureTypeAmbient:
		return "Ambient"
	case TextureTypeAmbientOcclusion:
		return "AmbientOcclusion"
	case TextureTypeBaseColor:
		return "BaseColor"
	case TextureTypeDiffuseRoughness:
		return "DiffuseRoughness"
	case TextureTypeDisplacement:
		return "Displacement"
	case TextureTypeEmissionColor:
		return "EmissionColor"
	case TextureTypeEmissive:
		return "Emissive"
	case TextureTypeHeight:
		return "Height"
	case TextureTypeLightmap:
		return "Lightmap"
	case TextureTypeMetalness:
		return "Metalness"
	case TextureTypeNormal:
		return "Normal"
	case TextureTypeNormalCamera:
		return "NormalCamera"
	case TextureTypeOpacity:
		return "Opacity"
	case TextureTypeReflection:
		return "Reflection"
	case TextureTypeShininess:
		return "Shininess"
	case TextureTypeUnknown:
		return "Unknown"
	default:
		return "Invalid"
	}
}

type MatPropertyTypeInfo int32

const (
	MatPropTypeInfoFloat32 MatPropertyTypeInfo = iota + 1
	MatPropTypeInfoFloat64
	MatPropTypeInfoString
	MatPropTypeInfoInt32

	MatPropTypeInfoBuffer
)

func (mpti MatPropertyTypeInfo) String() string {

	switch mpti {
	case MatPropTypeInfoFloat32:
		return "Float32"
	case MatPropTypeInfoFloat64:
		return "Float64"
	case MatPropTypeInfoString:
		return "String"
	case MatPropTypeInfoInt32:
		return "Int32"
	case MatPropTypeInfoBuffer:
		return "Buffer"
	default:
		return "Unknown"
	}
}

type MetadataType int32

const (
	MetadataTypeBool    MetadataType = 0
	MetadataTypeInt32   MetadataType = 1
	MetadataTypeUint64  MetadataType = 2
	MetadataTypeFloat32 MetadataType = 3
	MetadataTypeFloat64 MetadataType = 4
	MetadataTypeString  MetadataType = 5
	MetadataTypeVec3    MetadataType = 6
	MetadataTypeMAX     MetadataType = 7
)
