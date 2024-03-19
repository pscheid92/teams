package main

import (
	"fmt"
	"github.com/pscheid/teams/internal"
	"github.com/spf13/cobra"
	"log"
	"slices"
)

func buildKeysCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "keys",
		Short: "Manage your local keys",
	}
	command.AddCommand(
		buildGenerateKeyCmd(),
		buildListKeysCmd(),
		buildShowPublicKeyCmd(),
		buildShowPrivateKeyCmd(),
	)
	return command
}

func buildGenerateKeyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "generate",
		Short: "Generate a new public and private key pais",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]

			keys, err := internal.GenerateKeys()
			if err != nil {
				log.Fatalln(err)
			}

			template := `add this to your local configuration file

  %s:
    public: %s
    private: %s

`
			fmt.Printf(template, username, keys.Public, keys.Private)
		},
	}
}

func buildListKeysCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List you configured users",
		Run: func(cmd *cobra.Command, args []string) {
			app := cmd.Context().Value("app").(*AppContext)
			keysSet, err := app.BuildKeysSet()
			if err != nil {
				log.Fatalln(err)
			}

			users := keysSet.ListUsers()
			slices.Sort(users)
			for _, user := range users {
				fmt.Println(user)
			}
		},
	}
}

func buildShowPrivateKeyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "private",
		Short: "Show private key of user",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]
			app := cmd.Context().Value("app").(*AppContext)

			keysSet, err := app.BuildKeysSet()
			if err != nil {
				log.Fatalln(err)
			}

			key, ok := keysSet.GetKeys(username)
			if !ok {
				log.Fatalln("unknown username")
			}

			fmt.Println(key.Private)
		},
	}
}

func buildShowPublicKeyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "public",
		Short: "Show public key of user",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]
			app := cmd.Context().Value("app").(*AppContext)

			keysSet, err := app.BuildKeysSet()
			if err != nil {
				log.Fatalln(err)
			}

			key, ok := keysSet.GetKeys(username)
			if !ok {
				log.Fatalln("unknown username")
			}

			fmt.Println(key.Public)
		},
	}
}
