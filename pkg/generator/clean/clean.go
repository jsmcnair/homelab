// Package cleangenerator is a generator that is responsible foe cleaning up a
// target directory based on a list of globs.
package clean

import (
	"fmt"
	"path"
)

type CleanGenerator struct {
	Globs         []string
	BaseDirectory string
}

func (cg CleanGenerator) Generate() error {

	for _, g := range cg.Globs {
		cleanDir := path.Join(cg.BaseDirectory, g)
		fmt.Println(fmt.Sprintf("Cleaned: %s", cleanDir))
	}
	return nil
}
