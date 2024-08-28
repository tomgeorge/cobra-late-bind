/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"cobra-late-bind/cmd"
	"fmt"
	"os"

	"github.com/google/go-github/v55/github"
)

func main() {
	config := cmd.NewConfig()
	cli := github.NewClient(nil).WithAuthToken(config.Data.Github.Token)
	a := cmd.App{
		Config:     config,
		RepoLister: &cmd.GhRepoLister{Cli: cli},
	}
	if err := cmd.NewRootCommand(a).Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
