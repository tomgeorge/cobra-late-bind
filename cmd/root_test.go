package cmd

import (
	"os"
	"testing"
)

func Test_Uses_Flag_Default(t *testing.T) {
	app := App{
		Config: NewConfig(),
	}
	err := NewRootCommand(app).Execute()

	if err != nil {
		t.Fatalf("did not expect an error but got %v", err)
	}

	if app.Config.Data.Github.Token != "bar" {
		t.Fatalf("token is %s", app.Config.Data.Github.Token)
	}
	if app.Config.GetString("github.token") != "bar" {
		t.Fatalf("expected viper.GetString() to be bar but was %s", app.Config.GetString("github.token"))
	}
}

func TestFlag(t *testing.T) {
	err := os.Setenv("GITHUB_TOKEN", "foo")
	if err != nil {
		t.Fatalf("setting env var: %v", err)
	}
	defer os.Unsetenv("GITHUB_TOKEN")

	val, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		t.Fatalf("Could not look up GITHUB_TOKEN")
	}
	if val != "foo" {
		t.Fatalf("envvar is not foo")
	}
	c := NewConfig()
	app := App{
		Config: c,
	}
	cmd := NewRootCommand(app)
	cmd.Execute()
	fromViper := c.GetString("github.token")
	if fromViper != "foo" {
		t.Fatalf("viper.GetString(), wanted foo but was %s", fromViper)
	}
	if app.Config.Data.Github.Token != "foo" {
		t.Fatalf("Config.Data.Github.Tooken: wanted foo but got %s", app.Config.Data.Github.Token)
	}
}

func TestFlagPrecedence(t *testing.T) {
	err := os.Setenv("GITHUB_TOKEN", "foo")
	if err != nil {
		t.Fatalf("setting env var: %v", err)
	}
	defer os.Unsetenv("GITHUB_TOKEN")

	val, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		t.Fatalf("Could not look up GITHUB_TOKEN")
	}
	if val != "foo" {
		t.Fatalf("envvar is not foo")
	}
	c := NewConfig()
	app := App{
		Config: c,
	}

	cmd := NewRootCommand(app)
	cmd.SetArgs([]string{"--github-token", "baz"})
	cmd.Execute()
	// fromViper := c.v.GetString("github.token")
	// if fromViper != "baz" {
	// 	t.Fatalf("viper.GetString(), wanted baz but was %s", fromViper)
	// }
	if app.Config.Data.Github.Token != "baz" {
		t.Fatalf("Config.Data.Github.Tooken: wanted foo but got %s", app.Config.Data.Github.Token)
	}
}
