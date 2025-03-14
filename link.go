package assimp

// #include <stdlib.h>
// #cgo linux CXXFLAGS: -I ./lib/linux -std=c++14
// #cgo darwin,amd64 CXXFLAGS: -I ./lib/darwin  -std=gnu++14
// #cgo darwin,arm64 CXXFLAGS: -I ./lib/darwin_arm  -std=gnu++14
// #cgo darwin,amd64 LDFLAGS: -L ./lib/darwin -lassimp -lzlibstatic -lc++
// #cgo darwin,arm64 LDFLAGS: -L ./lib/darwin_arm -lassimp -lzlibstatic -lc++
// #cgo linux LDFLAGS: -L ./lib/linux -Wl,--start-group -lpthread -ldl -lstdc++ -lm -lassimp -lzlibstatic  -Wl,--end-group
import "C"
