package assimp

import (
	"fmt"
	"testing"
)

func TestGetSupportedImporterCount(t *testing.T) {
	count := GetSupportedImporterCount()
	if count <= 0 {
		t.Errorf("Expected positive number of supported importers, got %d", count)
	}
	t.Logf("Assimp supports %d import formats", count)
}

func TestGetImportFormatDescription(t *testing.T) {
	count := GetSupportedImporterCount()
	if count == 0 {
		t.Skip("No importers available")
	}

	// Test valid indices
	for i := 0; i < count && i < 3; i++ {
		desc := GetImportFormatDescription(i)
		if desc == nil {
			t.Errorf("Expected non-nil description for index %d", i)
			continue
		}
		t.Logf("Importer %d: %s supports extensions: %s", i, desc.Name, desc.FileExtensions)
	}

	// Test invalid indices
	invalidDesc := GetImportFormatDescription(-1)
	if invalidDesc != nil {
		t.Error("Expected nil description for negative index")
	}

	invalidDesc = GetImportFormatDescription(count)
	if invalidDesc != nil {
		t.Error("Expected nil description for index >= count")
	}
}

func TestIsFormatSupported(t *testing.T) {
	tests := []struct {
		extension string
		expected  bool
	}{
		{"obj", true},
		{"OBJ", true},
		{"fbx", true},
		{"FBX", true},
		{"stl", true},
		{"STL", true},
		{"gltf", true},
		{"glb", true}, // Note: GLB might be handled by GLTF importer
		{"ply", true},
		{"3ds", true},
		{"xyz", false}, // unlikely to be supported
		{"unsupported", false},
		{"", false}, // empty string should return false
	}

	for _, test := range tests {
		supported := IsFormatSupported(test.extension)
		if supported != test.expected {
			if test.expected {
				t.Errorf("Expected format %s to be supported, but it isn't", test.extension)
			} else {
				t.Errorf("Expected format %s to be unsupported, but it is", test.extension)
			}
		} else {
			t.Logf("Format %s support check: %v (expected: %v)", test.extension, supported, test.expected)
		}
	}

	// Test with leading dot
	withDot := IsFormatSupported(".obj")
	if !withDot {
		t.Error("Expected .obj (with leading dot) to be supported")
	}
}

func TestGetAllSupportedFormats(t *testing.T) {
	formats := GetAllSupportedFormats()
	if len(formats) == 0 {
		t.Error("Expected at least one supported format")
	}

	t.Logf("Found %d supported file extensions:", len(formats))
	for i, ext := range formats {
		if i >= 10 { // Limit output for readability
			t.Log("...")
			break
		}
		t.Logf("  %s", ext)
	}

	// Check that common formats are included
	commonFormats := []string{"obj", "fbx", "stl", "gltf", "3ds"}
	for _, format := range commonFormats {
		found := false
		for _, supported := range formats {
			if supported == format {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected %s to be in supported formats", format)
		}
	}
}

func TestGetImporterInfo(t *testing.T) {
	// Test with a known supported format
	info := GetImporterInfo("obj")
	if info == nil {
		t.Error("Expected non-nil info for obj format")
	} else {
		t.Logf("OBJ Importer Info:")
		t.Logf("  Name: %s", info.Name)
		t.Logf("  Extensions: %s", info.FileExtensions)
		t.Logf("  Author: %s", info.Author)
	}

	// Test with unsupported format
	info = GetImporterInfo("unsupported")
	if info != nil {
		t.Error("Expected nil info for unsupported format")
	}

	// Test with leading dot
	info = GetImporterInfo(".fbx")
	if info == nil {
		t.Error("Expected non-nil info for .fbx format")
	}
}

func ExampleIsFormatSupported() {
	fmt.Println("Checking format support:")

	formats := []string{"obj", "fbx", "stl", "xyz"}
	for _, format := range formats {
		supported := IsFormatSupported(format)
		fmt.Printf("%s: %v\n", format, supported)
	}

	// Output:
	// Checking format support:
	// obj: true
	// fbx: true
	// stl: true
	// xyz: false
}

func ExampleGetAllSupportedFormats() {
	formats := GetAllSupportedFormats()
	fmt.Printf("Total supported formats: %d\n", len(formats))

	// Show first 5 formats
	for i := 0; i < 5 && i < len(formats); i++ {
		fmt.Printf("  %s\n", formats[i])
	}

	// Output will vary based on Assimp version
}

func ExampleGetImportFormatDescription() {
	count := GetSupportedImporterCount()
	if count > 0 {
		desc := GetImportFormatDescription(0)
		if desc != nil {
			fmt.Printf("First importer: %s\n", desc.Name)
			fmt.Printf("Extensions: %s\n", desc.FileExtensions)
		}
	}

	// Output will vary based on Assimp version
}
