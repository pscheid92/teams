package main

import (
	"fmt"
	"github.com/pscheid/teams/internal"
	"github.com/spf13/cobra"
	"log"
	"time"
)

func buildLoginCmd(repository internal.KeysRepository, client *internal.Client) *cobra.Command {
	return &cobra.Command{
		Use:  "login",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]

			server, err := cmd.Flags().GetString("server")
			if err != nil {
				log.Fatalln(err)
			}
			client := internal.NewClient(server)

			key, err := repository.LoadPrivateKey(username)
			if err != nil {
				log.Fatalln(err)
			}

			now := time.Now()
			challenge := internal.CreateChallenge(username, now, key)

			requestBody := internal.LoginRequest{
				Username:  username,
				Timestamp: now,
				Challenge: challenge,
			}

			response, err := client.Login(requestBody)
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Println(response.AccessToken)
		},
	}
}

func buildVerifyCmd(client *internal.Client) *cobra.Command {
	return &cobra.Command{
		Use:  "verify",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			accessToken := args[0]

			response, err := client.Verify(accessToken)
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Printf("verified for %s\n", response.Username)
		},
	}
}

func buildListTeamCmd(client *internal.Client) *cobra.Command {
	return &cobra.Command{
		Use:  "team",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			teamID := args[0]

			team, err := client.Team(teamID)
			if err != nil {
				log.Fatalln(err)
			}

			for _, username := range team.Members {
				fmt.Println(username)
			}
		},
	}
}
