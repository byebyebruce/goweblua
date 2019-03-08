package luajit

/*
#cgo CFLAGS: -I../c/LuaJIT-2.1.0-beta1/src -I../c/pbc -I../c/lpeg -I../c/pbc -I../c/cjson
#cgo LDFLAGS: -lluajit -lpbc -llpeg -lcjson -L../c/LuaJIT-2.1.0-beta1/src -L../c/pbc -L../c/lpeg -L../c/pbc -L../c/cjson
#cgo linux LDFLAGS: -lm -ldl

#include <lua.h>
#include <lauxlib.h>
#include <luajit.h>
#include <lualib.h>
#include <stddef.h>
#include <stdlib.h>

extern int			luaopen_protobuf_c(lua_State*);
extern int			luaopen_lpeg(lua_State*);
extern int			luaopen_cjson(lua_State*);
*/
import "C"
import "unsafe"

func (s *State) Openpbclibs() {
	C.luaopen_protobuf_c(s.l)
}

func (s *State) Openlpeglibs() {
	C.luaopen_lpeg(s.l)
}

func (s *State) Opencjsonlibs() {
	C.luaopen_cjson(s.l)
}

// Pushes the string str onto the stack.
func (s *State) Pushlstring(str string, len int) {
	cs := C.CString(str)
	defer C.free(unsafe.Pointer(cs))
	C.lua_pushlstring(s.l, cs, C.size_t(len))
}

func (s *State) Tolstring(index int) (string, int) {
	var c C.size_t = 0
	str := C.lua_tolstring(s.l, C.int(index), &c)
	if str == nil {
		return "", 0
	}
	return C.GoString(str), int(c)
}
