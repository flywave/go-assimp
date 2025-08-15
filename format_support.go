package assimp

/*
#cgo linux CFLAGS: -I ./lib/linux
#cgo darwin,amd64 CFLAGS: -I ./lib/darwin
#cgo darwin,arm64 CFLAGS: -I ./lib/darwin_arm

#include <assimp/cimport.h>
#include <assimp/importerdesc.h>
#include <stdlib.h>
*/
import "C"
import (
	"strings"
	"unsafe"
)

// ImporterDesc represents the metadata about a particular importer
// This is a Go wrapper for C struct aiImporterDesc
type ImporterDesc struct {
	Name           string // Full name of the importer (i.e. Blender3D importer)
	Author         string // Original author (left blank if unknown or whole assimp team)
	Maintainer     string // Current maintainer, left blank if the author maintains
	Comments       string // Implementation comments, i.e. unimplemented features
	Flags          uint   // Flags indicating some characteristics
	MinMajor       uint   // Minimum format version that can be loaded (major)
	MinMinor       uint   // Minimum format version that can be loaded (minor)
	MaxMajor       uint   // Maximum format version that can be loaded (major)
	MaxMinor       uint   // Maximum format version that can be loaded (minor)
	FileExtensions string // List of file extensions this importer can handle
}

// GetSupportedImporterCount returns the number of import file formats available
// in the current Assimp build. Use GetImportFormatDescription to retrieve
// information about specific import formats.
func GetSupportedImporterCount() int {
	return int(C.aiGetImportFormatCount())
}

// GetImportFormatDescription returns a description of the nth import file format.
// Use GetSupportedImporterCount() to learn how many import formats are supported.
// Index must be in range 0 to GetSupportedImporterCount()-1.
// Returns nil if index is out of range.
func GetImportFormatDescription(index int) *ImporterDesc {
	if index < 0 || index >= GetSupportedImporterCount() {
		return nil
	}

	cDesc := C.aiGetImportFormatDescription(C.size_t(index))
	if cDesc == nil {
		return nil
	}

	return &ImporterDesc{
		Name:           C.GoString(cDesc.mName),
		Author:         C.GoString(cDesc.mAuthor),
		Maintainer:     C.GoString(cDesc.mMaintainer),
		Comments:       C.GoString(cDesc.mComments),
		Flags:          uint(cDesc.mFlags),
		MinMajor:       uint(cDesc.mMinMajor),
		MinMinor:       uint(cDesc.mMinMinor),
		MaxMajor:       uint(cDesc.mMaxMajor),
		MaxMinor:       uint(cDesc.mMaxMinor),
		FileExtensions: C.GoString(cDesc.mFileExtensions),
	}
}

// IsFormatSupported checks if a given file extension is supported by Assimp.
// The extension parameter can include or exclude the leading dot (e.g., "obj" or ".obj").
// Returns true if the format is supported, false otherwise.
func IsFormatSupported(extension string) bool {
	// Remove leading dot if present
	ext := strings.ToLower(strings.TrimPrefix(extension, "."))

	// Create C string for extension
	cExt := C.CString(ext)
	defer C.free(unsafe.Pointer(cExt))

	// Check if there's an importer for this extension
	cDesc := C.aiGetImporterDesc(cExt)
	return cDesc != nil
}

// GetAllSupportedFormats returns a slice of all supported file extensions.
// Each extension is returned without the leading dot and in lowercase.
func GetAllSupportedFormats() []string {
	count := GetSupportedImporterCount()
	formats := make([]string, 0)

	for i := 0; i < count; i++ {
		desc := GetImportFormatDescription(i)
		if desc != nil && desc.FileExtensions != "" {
			// Split extensions by space or comma
			texts := strings.FieldsFunc(desc.FileExtensions, func(r rune) bool {
				return r == ' ' || r == ',' || r == ';'
			})

			for _, ext := range texts {
				if ext = strings.TrimSpace(ext); ext != "" {
					formats = append(formats, strings.ToLower(ext))
				}
			}
		}
	}

	return formats
}

// GetImporterInfo returns detailed information about the importer that handles the given extension.
// The extension parameter can include or exclude the leading dot.
// Returns nil if the format is not supported.
func GetImporterInfo(extension string) *ImporterDesc {
	// Remove leading dot if present
	ext := strings.ToLower(strings.TrimPrefix(extension, "."))

	// Create C string for extension
	cExt := C.CString(ext)
	defer C.free(unsafe.Pointer(cExt))

	// Get importer description for this extension
	cDesc := C.aiGetImporterDesc(cExt)
	if cDesc == nil {
		return nil
	}

	return &ImporterDesc{
		Name:           C.GoString(cDesc.mName),
		Author:         C.GoString(cDesc.mAuthor),
		Maintainer:     C.GoString(cDesc.mMaintainer),
		Comments:       C.GoString(cDesc.mComments),
		Flags:          uint(cDesc.mFlags),
		MinMajor:       uint(cDesc.mMinMajor),
		MinMinor:       uint(cDesc.mMinMinor),
		MaxMajor:       uint(cDesc.mMaxMajor),
		MaxMinor:       uint(cDesc.mMaxMinor),
		FileExtensions: C.GoString(cDesc.mFileExtensions),
	}
}

// PrintSupportedFormats prints all supported formats to stdout for debugging purposes.
func PrintSupportedFormats() {
	count := GetSupportedImporterCount()
	println("Assimp supports the following formats:")
	println("====================================")

	for i := 0; i < count; i++ {
		desc := GetImportFormatDescription(i)
		if desc != nil {
			println("Importer:", desc.Name)
			println("  Extensions:", desc.FileExtensions)
			println("  Author:", desc.Author)
			println("  Comments:", desc.Comments)
			println()
		}
	}
}
