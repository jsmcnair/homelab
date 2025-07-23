package generator

type Generator interface {
	Generate() error
}

type GeneratorConfig struct {
	generators []Generator
	directory  string
}
