package main

import (
	"fmt"
	"github.com/pscheid/teams/internal"
	"github.com/spf13/cobra"
	"log"
)

func buildKeysCmd(repository internal.KeysRepository) *cobra.Command {
	command := &cobra.Command{Use: "keys"}

	command.AddCommand(
		buildListKeysCmd(repository),
		buildGenerateKeyCmd(repository),
		buildDeleteKeyCmd(repository),
		buildShowPublicKeyCmd(repository),
		buildShowPrivateKeyCmd(repository),
	)

	return command
}

func buildListKeysCmd(repository internal.KeysRepository) *cobra.Command {
	return &cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
			users, err := repository.ListUsers()
			if err != nil {
				log.Fatalln(err)
			}

			for _, user := range users {
				fmt.Println(user)
			}
		},
	}
}

func buildGenerateKeyCmd(repository internal.KeysRepository) *cobra.Command {
	return &cobra.Command{
		Use:  "generate",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]

			public, _, err := repository.GenerateKeys(username)
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Println(public)
		},
	}
}

func buildDeleteKeyCmd(repository internal.KeysRepository) *cobra.Command {
	return &cobra.Command{
		Use:  "delete",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]

			err := repository.DeleteKeys(username)
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Println("success")
		},
	}
}

func buildShowPrivateKeyCmd(repository internal.KeysRepository) *cobra.Command {
	return &cobra.Command{
		Use:  "private",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]

			key, err := repository.LoadEncodedPrivateKey(username)
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Println(key)
		},
	}
}

func buildShowPublicKeyCmd(repository internal.KeysRepository) *cobra.Command {
	return &cobra.Command{
		Use:  "public",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]

			key, err := repository.LoadEncodedPublicKey(username)
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Println(key)
		},
	}
}
