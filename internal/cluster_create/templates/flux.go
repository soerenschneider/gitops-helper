package templates

import (
	"gitops-helper/pkg"
	"log/slog"
	"os"
	"path"
	"text/template"
)

const fluxRepoTemplate = `apiVersion: source.toolkit.fluxcd.io/v1
kind: GitRepository
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
spec:
  interval: {{ .Interval }}
  url: {{ .RepoURL }}
  ref:
    branch: {{ .Branch }}
---
apiVersion: kustomize.toolkit.fluxcd.io/v1
kind: Kustomization
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
spec:
  interval: {{ .Interval }}
  path: {{ .Path }}
  prune: true
  wait: false
  sourceRef:
    kind: {{ .Name }}
    name: {{ .Namespace }}
`

type FluxRepoConfig struct {
	Name       string
	Namespace  string
	Interval   string
	RepoURL    string
	Branch     string
	Path       string
	SecretName string
}

func WriteFluxApplication(dir string) error {
	repoUrl, err := pkg.GetGithubRepositoryUrl()
	if err != nil {
		return err
	}

	data := FluxRepoConfig{
		Name:      path.Base(dir),
		Namespace: "flux-system",
		Interval:  "1m0s",
		RepoURL:   repoUrl,
		Branch:    "main",
		Path:      dir,
	}

	return writeFluxResources(path.Join(dir, "fluxcd.yaml"), data)
}

func writeFluxResources(filePath string, data FluxRepoConfig) error {
	tmpl, err := template.New("fluxResources").Parse(fluxRepoTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			slog.Warn("could not close file", "file", filePath, "err", err)
		}
	}()

	if err := tmpl.Execute(file, data); err != nil {
		return err
	}
	return nil
}
