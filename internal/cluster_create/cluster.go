package cluster_create

import (
	"cmp"
	"fmt"
	"gitops-helper/internal"
	"gitops-helper/internal/cluster_create/templates"
	"gitops-helper/internal/cluster_create/tui"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const defaultClusterName = "default"

func build(choices tui.UserData) error {
	clusterName := cmp.Or(choices.ClusterName, defaultClusterName)
	//if !NewConfirmation(clusterName) {
	//	return nil
	//}

	clusterDir, err := createCluster(clusterName)
	if err != nil {
		return err
	}

	resources := getFoldersForComponents(choices.Components, choices.GitOpsTool)
	if err := createKustomizeResource(clusterDir, choices.GitOpsTool, resources, nil); err != nil {
		return err
	}

	switch choices.GitOpsTool {
	case internal.ArgoCD:
		return templates.WriteArgoCDApplication(clusterDir)
	case internal.FluxCD:
		return templates.WriteFluxApplication(clusterDir)
	}

	return fmt.Errorf("unknown gitops tool %q", choices.GitOpsTool)
}

func getFoldersForComponents(components []string, tool string) []string {
	ret := make([]string, 0, len(components))
	for _, c := range components {
		componentName := strings.TrimSpace(strings.Split(c, " ")[0])
		ret = append(ret, filepath.Join("../../", internal.ComponentsDir, tool, componentName))
	}
	return ret
}

func createKustomizeResource(dir string, gitopsTool string, resources, components []string) error {
	kustomization := Kustomization{
		APIVersion: "kustomize.config.k8s.io/v1beta1",
		Kind:       "Kustomization",
		Resources:  resources,
		Components: components,
	}

	yamlData, err := yaml.Marshal(kustomization)
	if err != nil {
		fmt.Printf("Error marshaling Kustomization: %v\n", err)
		return err
	}

	err = os.WriteFile(filepath.Join(dir, "kustomization.yaml"), yamlData, 0644)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return err
	}

	return nil
}

func createCluster(clusterName string) (string, error) {
	if clusterName == "" {
		return "", fmt.Errorf("clusterName is required")
	}

	dir := filepath.Join(internal.ClustersDir, clusterName)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create cluster directory: %w", err)
	}

	fmt.Printf("Cluster directory created at: %s\n", dir)
	return dir, nil
}
