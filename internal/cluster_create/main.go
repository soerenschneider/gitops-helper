package cluster_create

import (
	"errors"
	"fmt"
	"gitops-helper/internal"
	"gitops-helper/internal/cluster_create/tui"
	"gitops-helper/pkg"
)

func CreateCluster() error {
	if !pkg.IsGitRepo() {
		return errors.New("current dir is not a git repository")
	}

	_, err := pkg.GetGithubRepositoryUrl()
	if err != nil {
		return fmt.Errorf("your git repository does not have an origin configured: %w", err)
	}

	components, err := pkg.AutodetectComponents(internal.ComponentsDir)
	if err != nil {
		return fmt.Errorf("could not automatically detect components in dir %q: %w", internal.ComponentsDir, err)
	}

	choices, err := tui.RunWizard(components)
	if err != nil {
		return err
	}

	if err := choices.Validate(); err != nil {
		return fmt.Errorf("incomplete choices: %w", err)
	}

	if err := build(choices); err != nil {
		return fmt.Errorf("could not build cluster %q: %w", choices.ClusterName, err)
	}

	return nil
}
