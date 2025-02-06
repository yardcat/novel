package world

import (
	"path/filepath"
)

type Resources struct {
	World   string
	workDir string
}

func (r *Resources) Init(world string) {
	r.World = filepath.Join("world", world, "data")
	r.workDir, _ = filepath.Abs(".")
}

func (r *Resources) GetPath(path string) string {
	return filepath.Join(r.workDir, r.World, path)
}
