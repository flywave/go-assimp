package assimp

// #include <stdlib.h>
// #cgo CFLAGS: -I ./lib
// #cgo linux CXXFLAGS: -I ./lib -std=c++14
// #cgo darwin CXXFLAGS: -I ./lib  -std=gnu++14
// #cgo darwin,arm CXXFLAGS: -I ./lib  -std=gnu++14
// #cgo darwin LDFLAGS: -L ./lib/darwin -lassimp -lzlibstatic
// #cgo darwin,arm LDFLAGS: -L ./lib/darwin_arm -lassimp -lzlibstatic
// #cgo linux LDFLAGS: -L ./lib/linux -Wl,--start-group -lpthread -ldl -lstdc++ -lm -lassimp -lzlibstatic  -Wl,--end-group
import "C"
