package world

import (
	"path/filepath"
)

type Resources struct {
	World   string
	workDir string
}

func NewResources(world string) *Resources {
	r := &Resources{}
	r.World = filepath.Join("world", world, "data")
	r.workDir, _ = filepath.Abs(".")
	return r
}

func (r *Resources) GetPath(path string) string {
	return filepath.Join(r.workDir, r.World, path)
}
