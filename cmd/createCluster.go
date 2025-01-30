package cmd

import (
	"gitops-helper/internal/cluster_create"

	"github.com/spf13/cobra"
)

// createClusterCmd represents the createCluster command
var createClusterCmd = &cobra.Command{
	Use:   "managed-cluster",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true // https://github.com/spf13/cobra/issues/340
		return cluster_create.CreateCluster()
	},
}

func init() {
	createCmd.AddCommand(createClusterCmd)
}
