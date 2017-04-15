package main

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/miekg/dns"

	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

func buildLua(h *PatternHandler, rlog *bytes.Buffer, w dns.ResponseWriter, r *dns.Msg) *lua.LState {
	L := lua.NewState()
	setglobal := func(name string, value interface{}) {
		L.SetGlobal(name, luar.New(L, value))
	}
	setglobal("printf", fmt.Printf)
	setglobal("sprintf", fmt.Sprintf)
	setglobal("log", log)
	setglobal("rlog", rlog)
	return L
}

func (h *PatternHandler) ServeLua(rlog *bytes.Buffer, w dns.ResponseWriter, r *dns.Msg) error {
	L := buildLua(h, rlog, w, r)

	var luascript string
	if strings.HasPrefix(h.LuaScript, "/") {
		luascript = h.LuaScript
	} else {
		luascript = filepath.Join(h.RootDir, h.LuaScript)
	}
	err := L.DoFile(luascript)
	if err != nil {
		log.Warnf("lua_script %s: %s", luascript, err)
		return err
	}
	serveFunction := "ServeDns"
	args := []lua.LValue{
		luar.New(L, w),
		luar.New(L, r),
	}
	err = L.CallByParam(lua.P{
		Fn:      L.GetGlobal(serveFunction),
		NRet:    1,
		Protect: true,
	}, args...)
	if err != nil {
		log.Warnf("lua_script %s:%s: %s", h.LuaScript, serveFunction, err)
		return err
	}
	ret := L.Get(-1) // returned value
	L.Pop(1)         // remove received value

	log.Tracef("lua_script %s return %#v", h.LuaScript, ret)

	return nil
}
