package main

import (
	"fmt"
	"os"
	"path"

	gen "github.com/turong-dev/homelab/v2/pkg/generator"
)

func main() {

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
		return
	}

	configPath := path.Join(wd, "config.yaml")
	config := gen.Config{}
	
	if err = config.Read(configPath); err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}
	
	fmt.Println("Config loaded successfully!")
	for _, generator := range config.Generators {
		if err = generator.Generate(); err != nil {
			fmt.Println("Error running Generate: ", err.Error())
		}
	}
}
