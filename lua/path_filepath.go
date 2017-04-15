package lua

import "path/filepath"

type Filepath struct {
	EvalSymlinks  interface{}
	HasPrefix     interface{}
	Match         interface{}
	SplitList     interface{}
	SkipDir       interface{}
	Rel           interface{}
	ToSlash       interface{}
	VolumeName    interface{}
	Base          interface{}
	Dir           interface{}
	ErrBadPattern interface{}
	Join          interface{}
	Walk          interface{}
	Glob          interface{}
	FromSlash     interface{}
	Split         interface{}
	Abs           interface{}
	Clean         interface{}
	IsAbs         interface{}
	Ext           interface{}
}

func NewFilepath() *Filepath {
	return &Filepath{
		Dir:           filepath.Dir,
		Join:          filepath.Join,
		Glob:          filepath.Glob,
		ToSlash:       filepath.ToSlash,
		Rel:           filepath.Rel,
		Base:          filepath.Base,
		Clean:         filepath.Clean,
		Match:         filepath.Match,
		SkipDir:       filepath.SkipDir,
		Walk:          filepath.Walk,
		Ext:           filepath.Ext,
		SplitList:     filepath.SplitList,
		FromSlash:     filepath.FromSlash,
		ErrBadPattern: filepath.ErrBadPattern,
		IsAbs:         filepath.IsAbs,
		EvalSymlinks:  filepath.EvalSymlinks,
		HasPrefix:     filepath.HasPrefix,
		Split:         filepath.Split,
		VolumeName:    filepath.VolumeName,
		Abs:           filepath.Abs,
	}
}
