// Package url is responsible for pulling manifests from a url and outputting
// them to a directory. Components of the URL may contain go style template
// placeholders that will be replace by the values of keys in the Tokens map.
package url

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"text/template"

	"github.com/turong-dev/homelab/v2/pkg/utils"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type UrlGenerator struct {
	BaseDirectory string
	URL           string
	Destination   string
	Tokens        map[string]string
}

// Generate implements the generator interface.
func (ug UrlGenerator) Generate() error {

	destination := ug.Destination
	if ug.Destination == "" {
		destination = "manifests"
	}

	dest := path.Join(ug.BaseDirectory, destination)

	if err := utils.EnsureDestinationPath(dest); err != nil {
		return fmt.Errorf("Unable to write to destination %s: %w", dest, err)
	}

	return ug.DownloadManifests(dest)
}

// TemplateURL renders the URL template with the values from the Tokens map.
func (ug UrlGenerator) TemplateURL() (*url.URL, error) {
	tmpl, err := template.New("url").Parse(ug.URL)
	if err != nil {
		return nil, err
	}

	var newUrlBytes bytes.Buffer
	if err = tmpl.Execute(&newUrlBytes, ug.Tokens); err != nil {
		return nil, err
	}

	var newUrl *url.URL
	if newUrl, err = url.Parse(newUrlBytes.String()); err != nil {
		return nil, err
	}

	return newUrl, nil
}

// DownloadManifests downloads the manifests from the URL
func (ug UrlGenerator) DownloadManifests(dest string) error {

	url, err := ug.TemplateURL()
	if err != nil {
		return fmt.Errorf("failed to template url: %w", err)
	}

	resp, err := http.Get(url.String())
	if err != nil {
		return fmt.Errorf("failed to download manifests from %s: %w", url.String(), err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download manifests from %s: %s", url.String(), resp.Status)
	}

	decoder := yaml.NewDecoder(resp.Body)
	for {
		var parsed map[string]interface{}
		err := decoder.Decode(&parsed)
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("failed to decode manifest: %w", err)
			}
			break
		}
		
		var ok bool
		var kind string
		if kind, ok = parsed["kind"].(string); !ok {
			return fmt.Errorf("failed to parse kind")
		}
		
		metadata, ok := parsed["metadata"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("failed to parse metadata")
		}
		
		var name string
		if name, ok = metadata["name"].(string); !ok {
			return fmt.Errorf("failed to parse name")
		}
		
		var namespace string
		if namespace, ok = metadata["namespace"].(string); !ok {
			namespace = "_no_namespace"
		}

		namespaceDest := path.Join(dest, namespace)
		if err := os.MkdirAll(namespaceDest, 0644); err != nil {
			return fmt.Errorf("failed to create namespace directory: %w", err)
		}
		
		file, err := os.Create(path.Join(namespaceDest, fmt.Sprintf("%s_%s.yaml", kind, name)))
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}

		defer file.Close()

		if err := yaml.NewEncoder(file).Encode(parsed); err != nil {
			return fmt.Errorf("failed to encode manifest: %w", err)
		}
	}

	return nil
}
