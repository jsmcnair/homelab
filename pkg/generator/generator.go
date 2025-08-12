package generator

import (
	"net/url"
	
	genclean "github.com/turong-dev/homelab/v2/pkg/generator/clean"
	genurl "github.com/turong-dev/homelab/v2/pkg/generator/url"
)

func TestGenerators() {

	baseDirectory := "/we/are/here"
	
	cleangen := genclean.CleanGenerator{
		BaseDirectory: baseDirectory,
		Globs:         []string{"a/b", "b/c"},
	}

	urlgen := genurl.UrlGenerator{
		BaseDirectory: baseDirectory,
		Destination: "manifests",
	}
	if rem, err := url.Parse("https://foo.bar/manifests"); err != nil {
		panic("Could not parse URL")
	} else {
		urlgen.URL = rem.String()
	}

	gc := GeneratorConfig{}
	gc.generators = append(gc.generators, cleangen, urlgen)

	for _, gen := range gc.generators {
		if err := gen.Generate(); err != nil {
			panic("Error running generator")
		}
	}
}
