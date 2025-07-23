package generator

import (
	"encoding/json"
	"fmt"
	"os"

	genargocd "github.com/turong-dev/homelab/v2/pkg/generator/argocd"
	genclean "github.com/turong-dev/homelab/v2/pkg/generator/clean"
	genurl "github.com/turong-dev/homelab/v2/pkg/generator/url"
	"sigs.k8s.io/yaml"
)

// type Config stores the configuration for a given application. The
// RawGenerators field holds the raw bytes of the unmarshaled generators field
// in the config, prior to being unmarshaled to the correct type in the
// Generators field.
type Config struct {
	RawGenerators []json.RawMessage `json:"generators"`
	Generators    []Generator       `json:"-"`
}

// Read reads the configuration YAML from the given path and uses the K8S sig
// YAML library to decode it into the correct type. It uses the double decode
// approach to ensure that the YAML is properly decoded into the correct type,
// by first decoding the type field only and then using the unmarshal function
// for the corresponding type.
func (c *Config) Read(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("unable to read file: %v", err)
	}

	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		return fmt.Errorf("unable to parse yaml: %v", err)
	}

	for _, generator := range c.RawGenerators {
		typeStruct := struct {
			Type string `json:"type"`
		}{}

		err := unmarshal(generator, &typeStruct)
		if err != nil {
			return fmt.Errorf("unable to decode generator type: %v", err)
		}

		if typeStruct.Type == "" {
			return fmt.Errorf("generator type cannot be empty")
		}

		switch typeStruct.Type {
		case "argocd":
			var gen genargocd.ArgoCDGenerator
			err := unmarshal(generator, &gen)
			if err != nil {
				return fmt.Errorf("unable to cast generator to ArgoCDGenerator")
			}
			c.Generators = append(c.Generators, gen)
		case "clean":
			var gen genclean.CleanGenerator
			err := unmarshal(generator, &gen)
			if err != nil {
				return fmt.Errorf("unable to cast generator to CleanGenerator")
			}
			c.Generators = append(c.Generators, gen)
		case "url":
			var gen genurl.UrlGenerator
			err := unmarshal(generator, &gen)
			if err != nil {
				return fmt.Errorf("unable to cast generator to URLGenerator")
			}
			c.Generators = append(c.Generators, gen)
		}
	}
	return nil
}

func unmarshal[T any](raw []byte, t T) error {
	err := json.Unmarshal(raw, t)
	if err != nil {
		return fmt.Errorf("unable to parse json for %T", t)
	}

	return nil
}
