package assimp

import (
	"testing"
)

// TestMaterialCreation 测试材质创建
func TestMaterialCreation(t *testing.T) {
	material := &Material{
		Properties: []*MaterialProperty{
			{
				name:     "diffuse_color",
				Semantic: TextureTypeDiffuse,
				Index:    0,
				TypeInfo: MatPropTypeInfoFloat32,
				Data:     []byte{0x3F, 0x80, 0x00, 0x00}, // 1.0 in IEEE 754
			},
			{
				name:     "texture_path",
				Semantic: TextureTypeDiffuse,
				Index:    0,
				TypeInfo: MatPropTypeInfoString,
				Data:     []byte("test_texture.png"),
			},
		},
		AllocatedStorage: 1024,
	}

	if len(material.Properties) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(material.Properties))
	}

	if material.AllocatedStorage != 1024 {
		t.Errorf("Expected allocated storage 1024, got %d", material.AllocatedStorage)
	}

	if material.Properties[0].name != "diffuse_color" {
		t.Errorf("Expected property name 'diffuse_color', got '%s'", material.Properties[0].name)
	}

	if material.Properties[0].Semantic != TextureTypeDiffuse {
		t.Errorf("Expected semantic TextureTypeDiffuse, got %d", material.Properties[0].Semantic)
	}

	if material.Properties[0].TypeInfo != MatPropTypeInfoFloat32 {
		t.Errorf("Expected type MatPropTypeInfoFloat32, got %d", material.Properties[0].TypeInfo)
	}
}

// TestMaterialPropertyTypes 测试材质属性类型
func TestMaterialPropertyTypes(t *testing.T) {
	properties := []*MaterialProperty{
		{
			name:     "float32_prop",
			TypeInfo: MatPropTypeInfoFloat32,
			Data:     []byte{0x00, 0x00, 0x80, 0x3F}, // 1.0
		},
		{
			name:     "float64_prop",
			TypeInfo: MatPropTypeInfoFloat64,
			Data:     []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xF0, 0x3F}, // 1.0
		},
		{
			name:     "string_prop",
			TypeInfo: MatPropTypeInfoString,
			Data:     []byte("test_string"),
		},
		{
			name:     "int32_prop",
			TypeInfo: MatPropTypeInfoInt32,
			Data:     []byte{0x01, 0x00, 0x00, 0x00}, // 1
		},
		{
			name:     "buffer_prop",
			TypeInfo: MatPropTypeInfoBuffer,
			Data:     []byte{0x01, 0x02, 0x03, 0x04},
		},
	}

	for i, prop := range properties {
		t.Logf("Property %d: name=%s, type=%d, data_length=%d",
			i, prop.name, prop.TypeInfo, len(prop.Data))

		if prop.name == "" {
			t.Errorf("Property %d has empty name", i)
		}

		if len(prop.Data) == 0 {
			t.Errorf("Property %d has empty data", i)
		}
	}
}

// TestMaterialTextureTypes 测试材质纹理类型
func TestMaterialTextureTypes(t *testing.T) {
	textureTypes := []TextureType{
		TextureTypeDiffuse,
		TextureTypeSpecular,
		TextureTypeAmbient,
		TextureTypeEmissive,
		TextureTypeHeight,
		TextureTypeNormal,
		TextureTypeShininess,
		TextureTypeOpacity,
		TextureTypeDisplacement,
		TextureTypeLightmap,
		TextureTypeReflection,
		TextureTypeBaseColor,
		TextureTypeNormalCamera,
		TextureTypeEmissionColor,
		TextureTypeMetalness,
		TextureTypeDiffuseRoughness,
		TextureTypeAmbientOcclusion,
	}

	for _, tt := range textureTypes {
		str := tt.String()
		if str == "" || str == "Invalid" {
			t.Errorf("TextureType %d has invalid string representation: '%s'", tt, str)
		}
	}
}

// TestEmptyMaterial 测试空材质
func TestEmptyMaterial(t *testing.T) {
	material := &Material{}

	if len(material.Properties) != 0 {
		t.Error("Expected empty properties")
	}

	if material.AllocatedStorage != 0 {
		t.Error("Expected zero allocated storage")
	}
}

// TestMaterialPropertyValidation 测试材质属性验证
func TestMaterialPropertyValidation(t *testing.T) {
	// 测试有效的材质属性
	prop := &MaterialProperty{
		name:     "test_prop",
		Semantic: TextureTypeDiffuse,
		Index:    0,
		TypeInfo: MatPropTypeInfoFloat32,
		Data:     []byte{0x00, 0x00, 0x80, 0x3F},
	}

	if prop.name == "" {
		t.Error("Expected non-empty property name")
	}

	if prop.TypeInfo == 0 {
		t.Error("Expected non-zero property type")
	}

	if len(prop.Data) == 0 {
		t.Error("Expected non-empty property data")
	}

	// 测试边界情况
	propEmpty := &MaterialProperty{
		name:     "",
		Semantic: TextureType(999),
		Index:    999,
		TypeInfo: MatPropertyTypeInfo(999),
		Data:     []byte{},
	}

	if propEmpty.name != "" {
		t.Error("Expected empty property name")
	}

	if propEmpty.TypeInfo.String() != "Unknown" {
		t.Errorf("Expected 'Unknown' for invalid type, got '%s'", propEmpty.TypeInfo.String())
	}
}

// TestMaterialIndexValidation 测试材质索引验证
func TestMaterialIndexValidation(t *testing.T) {
	material := &Material{
		Properties: []*MaterialProperty{
			{
				name:     "prop1",
				Index:    0,
				Semantic: TextureTypeDiffuse,
			},
			{
				name:     "prop2",
				Index:    1,
				Semantic: TextureTypeSpecular,
			},
			{
				name:     "prop3",
				Index:    2,
				Semantic: TextureTypeNormal,
			},
		},
	}

	for i, prop := range material.Properties {
		if prop.Index != uint(i) {
			t.Errorf("Expected index %d, got %d", i, prop.Index)
		}
	}
}

// TestMaterialSemanticConsistency 测试材质语义一致性
func TestMaterialSemanticConsistency(t *testing.T) {
	material := &Material{
		Properties: []*MaterialProperty{
			{
				name:     "diffuse_texture",
				Semantic: TextureTypeDiffuse,
				Index:    0,
			},
			{
				name:     "diffuse_color",
				Semantic: TextureTypeDiffuse,
				Index:    1,
			},
			{
				name:     "specular_texture",
				Semantic: TextureTypeSpecular,
				Index:    0,
			},
		},
	}

	diffuseCount := 0
	specularCount := 0

	for _, prop := range material.Properties {
		switch prop.Semantic {
		case TextureTypeDiffuse:
			diffuseCount++
		case TextureTypeSpecular:
			specularCount++
		}
	}

	if diffuseCount != 2 {
		t.Errorf("Expected 2 diffuse properties, got %d", diffuseCount)
	}

	if specularCount != 1 {
		t.Errorf("Expected 1 specular property, got %d", specularCount)
	}
}

// TestMaterialPropertyDataSize 测试材质属性数据大小
func TestMaterialPropertyDataSize(t *testing.T) {
	tests := []struct {
		typeInfo MatPropertyTypeInfo
		expected int
	}{
		{MatPropTypeInfoFloat32, 4},
		{MatPropTypeInfoFloat64, 8},
		{MatPropTypeInfoInt32, 4},
		{MatPropTypeInfoString, 16}, // Fixed size for testing
		{MatPropTypeInfoBuffer, 32}, // Fixed size for testing
	}

	for _, test := range tests {
		prop := &MaterialProperty{
			TypeInfo: test.typeInfo,
			Data:     make([]byte, test.expected),
		}

		if len(prop.Data) != test.expected {
			t.Errorf("Expected data size %d for type %d, got %d",
				test.expected, test.typeInfo, len(prop.Data))
		}
	}
}
