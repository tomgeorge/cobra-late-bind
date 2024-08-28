/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type ListRepoOpts struct {
	owner string
}

func NewFooCommand(a App) *cobra.Command {
	// fooCmd represents the foo command
	opts := ListRepoOpts{}
	cmd := &cobra.Command{
		Use:   "foo",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Foo RunE ", a.Config.Data.Github.Token, a.RepoLister)
			repositories, err := a.ListRepositories(cmd.Context(), opts)
			cmd.Println(repositories)
			return err
		},
	}
	fmt.Println("NewFooCommand ", a.Config.Data.Github.Token, a.RepoLister)
	cmd.Flags().StringVar(&opts.owner, "owner", "", "owner to list repositories for")
	return cmd
}
