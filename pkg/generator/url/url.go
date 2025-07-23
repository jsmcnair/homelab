// Package url is responsible for pulling manifests from a url and outputting
// them to a directory. Components of the URL may contain go style template
// placeholders that will be replace by the values of keys in the Tokens map.
package url

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"text/template"

	"github.com/turong-dev/homelab/v2/pkg/utils"
)

type UrlGenerator struct {
	BaseDirectory string
	URL           url.URL
	Destination   string
	Tokens        map[string]string
}

// Generate implements the generator interface.
func (ug UrlGenerator) Generate() error {
	dest := path.Join(ug.BaseDirectory, ug.Destination)
	
	if err := utils.EnsureDestinationPath(dest); err != nil {
		return fmt.Errorf("Unable to write to destination %s: %w", dest, err)
	}
	
	return ug.DownloadManifests(dest)
}

// TemplateURL renders the URL template with the values from the Tokens map.
func (ug UrlGenerator) TemplateURL() (*url.URL, error) {
	tmpl, err := template.New("url").Parse(ug.URL.String())
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
	
	return nil
}