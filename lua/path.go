package lua

import "path"

type Path struct {
	Split         interface{}
	Match         interface{}
	Base          interface{}
	Dir           interface{}
	ErrBadPattern interface{}
	Join          interface{}
	Clean         interface{}
	IsAbs         interface{}
	Ext           interface{}
}

func NewPath() *Path {
	return &Path{
		Ext:           path.Ext,
		Split:         path.Split,
		Match:         path.Match,
		Base:          path.Base,
		IsAbs:         path.IsAbs,
		Clean:         path.Clean,
		ErrBadPattern: path.ErrBadPattern,
		Dir:           path.Dir,
		Join:          path.Join,
	}
}
