package main

import (
	"fmt"
	"github.com/pscheid/teams/internal"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

func main() {
	rootCmd := buildRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func buildRootCmd() *cobra.Command {
	command := &cobra.Command{
		Use:              "teams",
		Short:            "Manage your local teams accounts.",
		TraverseChildren: true,
	}

	command.PersistentFlags().String("server", "https://patrickscheid.de/teams/", "...usage...")

	directory := getAppDir()
	repository := internal.NewFSKeyRepository(directory)

	server, _ := command.PersistentFlags().GetString("server")
	client := internal.NewClient(server)

	command.AddCommand(
		buildKeysCmd(repository),
		buildLoginCmd(repository, client),
		buildVerifyCmd(client),
		buildListTeamCmd(client),
	)

	return command
}

func getAppDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	appDir := filepath.Join(dir, ".teams")
	if err := os.MkdirAll(appDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	return appDir
}
