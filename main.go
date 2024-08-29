/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"cobra-late-bind/cmd"
	"fmt"
	"os"
)

func main() {
	config := cmd.NewConfig()
	a := &cmd.App{
		Config:     config,
		RepoLister: nil,
	}
	if err := cmd.NewRootCommand(a).Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
