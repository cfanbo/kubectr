package app

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/cfanbo/kubectr/internal/version"
)

func NewCmd() *cobra.Command {
	o := newCtrOptions()

	var rootCmd = &cobra.Command{
		Use:   "ctr <pod-name>", // This is prefixed by kubectl in the custom usage template
		Short: "display all containers in the pod",
		Long: `display all containers in the pod.

You can invoke ctr through kubectl: "kubectl ctr <pod-name>"`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if ver, err := cmd.Flags().GetBool("version"); err != nil {
				return err
			} else if ver {
				fmt.Println(version.FullVersion())
				return nil
			}

			if len(args) == 0 {
				_ = cmd.Help()
				return nil
			}

			if err := o.Complete(cmd, args); err != nil {
				return err
			}
			o.Run()
			return nil
		},
	}

	rootCmd.Flags().BoolVarP(&o.version, "version", "v", false, "Displays the current version number")
	o.configFlags.AddFlags(rootCmd.Flags())

	return rootCmd
}

func newCtrOptions() *CtrOptions {
	return &CtrOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
	}
}
