package main

import (
	"errors"
	"fmt"
	"github.com/pscheid/teams/internal"
	"github.com/spf13/cobra"
	"log"
	"time"
)

func buildLoginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Obtain an OAuth token",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]
			app := cmd.Context().Value("app").(*AppContext)

			keysSet, err := app.BuildKeysSet()
			if err != nil {
				log.Fatalln(err)
			}

			client, err := app.BuildClient()
			if err != nil {
				log.Fatalln(err)
			}

			key, err := keysSet.GetPrivateKey(username)
			if errors.Is(err, internal.ErrKeysNotFound) {
				log.Fatalln("unknown username")
			}
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

func buildVerifyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "verify",
		Short: "Verify if a provided OAuth token is still valid",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			accessToken := args[0]
			app := cmd.Context().Value("app").(*AppContext)

			client, err := app.BuildClient()
			if err != nil {
				log.Fatalln(err)
			}

			response, err := client.Verify(accessToken)
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Printf("verified for %s\n", response.Username)
		},
	}
}

func buildListTeamCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List a team's members",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			teamID := args[0]
			app := cmd.Context().Value("app").(*AppContext)

			client, err := app.BuildClient()
			if err != nil {
				log.Fatalln(err)
			}

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
