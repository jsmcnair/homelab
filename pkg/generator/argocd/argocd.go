// Package argocd is responsible for generating an argocd application that
// can work with the drectory structure of an application in the applications
// directory.It can support Helm, Kustomize and plain manifests.
package argocd

import (
	"fmt"
	"os"
	"path"

	argocdv1alpha1 "github.com/argoproj/argo-cd/v3/pkg/apis/application/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const (
	ARGOCD_DEFAULT_CLUSTER = "in-cluster"
	ARGOCD_NAMESPACE       = "argocd"
	ARGOCD_DEFAULT_REPO    = "git@github.com:jsmcnair/homelab"
)

type ArgoCDGenerator struct {
	BaseDirectory             string                             `yaml:"-"`
	Name                      string                             `yaml:"name"`
	Namespace                 string                             `yaml:"namespace"`
	EnableExtraManifestSource bool                               `yaml:"enableExtraManifestSource"`
	PrivilegedNamespace       bool                               `yaml:"privilegedNamespace"`
	Sources                   []argocdv1alpha1.ApplicationSource `yaml:"sources"`
	IgnoreDifferences         argocdv1alpha1.IgnoreDifferences   `yaml:"ignoreDifferences"`
}

func (ag ArgoCDGenerator) Generate() error {
	
	if ag.Name == "" {
		return fmt.Errorf("name is required")
	}
	
	app := ag.ConfigureApplication()
	outFile := path.Join(ag.BaseDirectory, "application.yaml")

	bytes, err := yaml.Marshal(app)
	if err != nil {
		return fmt.Errorf("failed to marshal application: %w", err)
	}

	err = os.WriteFile(outFile, bytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (ag ArgoCDGenerator) ConfigureApplication() argocdv1alpha1.Application {
	appSpec := ag.GetDefaultArgoCDApplicationSpec()
	
	if ag.Namespace != "" {
		appSpec.Destination.Namespace = ag.Namespace
	}

	for _, source := range ag.Sources {
		appSpec.Sources = append(appSpec.Sources, source)
	}

	if ag.EnableExtraManifestSource {
		appSpec.Sources = append(appSpec.Sources, ag.GetExtraManifestSource())
	}

	if len(ag.IgnoreDifferences) > 0 {
		appSpec.IgnoreDifferences = ag.IgnoreDifferences
	}

	if ag.PrivilegedNamespace {

		if appSpec.SyncPolicy.ManagedNamespaceMetadata.Labels == nil {
			appSpec.SyncPolicy.ManagedNamespaceMetadata.Labels = make(map[string]string)
		}

		appSpec.SyncPolicy.ManagedNamespaceMetadata.Labels["pod-security.kubernetes.io/audit"] = "privileged"
		appSpec.SyncPolicy.ManagedNamespaceMetadata.Labels["pod-security.kubernetes.io/enforce"] = "privileged"
		appSpec.SyncPolicy.ManagedNamespaceMetadata.Labels["pod-security.kubernetes.io/warn"] = "privileged"
	}

	return argocdv1alpha1.Application{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "argoproj.io/v1alpha1",
			Kind:       "Application",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      ag.Name,
			Namespace: ARGOCD_NAMESPACE,
		},
		Spec: *appSpec,
	}

}

func (ag ArgoCDGenerator) GetDefaultArgoCDApplicationSpec() *argocdv1alpha1.ApplicationSpec {
	return &argocdv1alpha1.ApplicationSpec{
		Destination: argocdv1alpha1.ApplicationDestination{
			Name:      "in-cluster",
			Namespace: "default",
		},
		Project: "default",
		SyncPolicy: &argocdv1alpha1.SyncPolicy{
			Automated: &argocdv1alpha1.SyncPolicyAutomated{
				Prune:    true,
				SelfHeal: true,
			},
			SyncOptions:              []string{"CreateNamespace=true"},
			ManagedNamespaceMetadata: &argocdv1alpha1.ManagedNamespaceMetadata{},
		},
	}
}

func (ag ArgoCDGenerator) GetExtraManifestSource() argocdv1alpha1.ApplicationSource {
	return argocdv1alpha1.ApplicationSource{
		RepoURL: ARGOCD_DEFAULT_REPO,
		Path:    fmt.Sprintf("applications/%s/extra", ag.Name),
		Directory: &argocdv1alpha1.ApplicationSourceDirectory{
			Recurse: true,
		},
	}
}
