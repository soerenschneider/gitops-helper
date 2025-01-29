package templates

import (
	"gitops-helper/pkg"
	"log/slog"
	"os"
	"path"
	"text/template"
)

const argoTemplate = `apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
spec:
  project: {{ .Project }}
  source:
    repoURL: {{ .RepoURL }}
    targetRevision: {{ .TargetRevision }}
    path: {{ .Path }}
  destination:
    server: {{ .Server }}
    namespace: {{ .DestinationNamespace }}
  syncPolicy:
    automated:
      prune: {{ .Prune }}
      selfHeal: {{ .SelfHeal }}
    syncOptions:
    {{- range .SyncOptions }}
    - {{ . }}
    {{- end }}
`

type ApplicationData struct {
	Name                 string
	Namespace            string
	Project              string
	RepoURL              string
	TargetRevision       string
	Path                 string
	Server               string
	DestinationNamespace string
	Prune                bool
	SelfHeal             bool
	SyncOptions          []string
}

func WriteArgoCDApplication(dir string) error {
	repoUrl, err := pkg.GetGithubRepositoryUrl()
	if err != nil {
		return err
	}

	data := ApplicationData{
		Name:                 path.Base(dir),
		Namespace:            "argocd",
		Project:              "default",
		RepoURL:              repoUrl,
		TargetRevision:       "main",
		Path:                 dir,
		Server:               "https://kubernetes.default.svc",
		DestinationNamespace: "default",
		Prune:                true,
		SelfHeal:             true,
		SyncOptions:          []string{"CreateNamespace=true"},
	}

	return writeArgoCDApplication(path.Join(dir, "application.yaml"), data)
}

func writeArgoCDApplication(filePath string, data ApplicationData) error {
	tmpl, err := template.New("argoCDApplication").Parse(argoTemplate)
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
