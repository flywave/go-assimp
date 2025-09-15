module github.com/flywave/go-assimp

go 1.24

toolchain go1.24.4

require (
	github.com/chai2010/tiff v0.0.0-20211005095045-4ec2aa243943
	github.com/flywave/go-mst v0.0.0-20250814104510-37f0a6660bc0
	github.com/flywave/go3d v0.0.0-20250619003741-cab1a6ea6de6
	golang.org/x/image v0.26.0
)

require github.com/flywave/gltf v0.20.4-0.20250828104044-ebb99e75f3cc // indirect

replace github.com/flywave/go-mst => ../go-mst
