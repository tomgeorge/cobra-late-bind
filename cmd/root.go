/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v55/github"
	"github.com/jeremywohl/flatten"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type App struct {
	RepoLister RepoLister
	Config     *Config
}

type RepoLister interface {
	RepositoryNames(ctx context.Context, owner string) ([]string, error)
}

type GhRepoLister struct {
	Cli *github.Client
}

func (a *App) ListRepositories(ctx context.Context, opts ListRepoOpts) ([]string, error) {
	fmt.Println("App.ListRepositories()")
	return a.RepoLister.RepositoryNames(ctx, opts.owner)
}

func (a *App) BindServices() {
	fmt.Println("Binding services")
	if a.RepoLister == nil {
		a.RepoLister = &GhRepoLister{Cli: github.NewClient(nil).WithAuthToken(a.Config.Data.Github.Token)}
		fmt.Println("Set repolister")
	}
}
func (rl *GhRepoLister) RepositoryNames(ctx context.Context, owner string) ([]string, error) {
	repositories, _, err := rl.Cli.Repositories.List(ctx, owner, nil)
	if err != nil {
		return []string{}, err
	}
	names := make([]string, len(repositories))
	for i, repo := range repositories {
		names[i] = repo.GetName()
	}
	return names, nil
}

type Config struct {
	Data Data `mapstructure:"data"`
	*viper.Viper
}

type Data struct {
	Github Github `mapstructure:"github"`
}

type Github struct {
	Token string `mapstructure:"token"`
}

func NewRootCommand(a *App) *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:   "cobra-late-bind",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			a.Config.BindFlags(cmd)
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("Root: ", a.Config.Data.Github.Token)
		},
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
	}

	rootCmd.PersistentFlags().StringVar(&a.Config.Data.Github.Token, "github-token", "bar", "github token")
	rootCmd.AddCommand(NewFooCommand(a))
	return rootCmd
}

func (c *Config) BindFlags(cmd *cobra.Command) {
	err := c.BindPFlags(cmd.Flags())
	if err != nil {
		panic(fmt.Sprintf("Binding pflags: %v", err))
	}
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := f.Name
		if f.Changed {
			c.Viper.Set(configName, f.Value)
		}
		if !f.Changed && c.Viper.IsSet(configName) {
			val := c.Viper.Get(configName)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}

func DefaultConfigData() *Data {
	return &Data{
		Github: Github{
			Token: "",
		},
	}
}

func NewConfig() *Config {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	v.AutomaticEnv()
	c := &Config{}
	c.Viper = v
	c.BindEnv()
	if err := v.Unmarshal(&c.Data); err != nil {
		panic(fmt.Sprintf("Unmarshalling config: %v", err))
	}
	return c
}

func (c *Config) Load() {
	data := &Config{}
	c.Viper.Unmarshal(&data.Data)
	c.Data = data.Data
}

// https://github.com/spf13/viper/issues/761
func (c *Config) BindEnv() error {
	structure := map[string]interface{}{}
	if err := mapstructure.Decode(c.Data, &structure); err != nil {
		return err
	}
	defaults, err := flatten.Flatten(structure, "", flatten.DotStyle)
	if err != nil {
		return err
	}
	for k := range defaults {
		if err := c.Viper.BindEnv(k); err != nil {
			return err
		}
	}
	return nil
}
