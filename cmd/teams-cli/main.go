package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/pscheid/teams/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func main() {
	var configFile string

	rootCmd := &cobra.Command{
		Use:   "teams",
		Short: "Make Zalando Postgres Operator Teams API easy to use",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			app, err := NewAppContext(cmd, configFile)
			cobra.CheckErr(err)

			ctx := context.WithValue(cmd.Context(), "app", app)
			cmd.SetContext(ctx)
			return nil
		},
	}

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "path to configuration file")
	rootCmd.PersistentFlags().StringP("server", "s", "", "url of target server")

	rootCmd.AddCommand(
		buildKeysCmd(),
		buildLoginCmd(),
		buildVerifyCmd(),
		buildListTeamCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type AppContext struct {
	config *viper.Viper
}

func NewAppContext(cmd *cobra.Command, configFile string) (*AppContext, error) {
	config := viper.NewWithOptions(viper.KeyDelimiter("::"))
	config.SetConfigType("yaml")

	if configFile != "" {
		config.SetConfigFile(configFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		config.AddConfigPath(home)
		config.AddConfigPath(".")
		config.SetConfigName(".teams")
	}

	_ = config.BindPFlag("server", cmd.Flags().Lookup("server"))

	err := config.ReadInConfig()
	if errors.Is(err, viper.ConfigFileNotFoundError{}) {
		return nil, errors.New("no config file found")
	}
	if err != nil {
		return nil, err
	}

	app := &AppContext{config: config}
	return app, nil
}

func (app *AppContext) BuildClient() (*internal.Client, error) {
	server := app.config.GetString("server")
	if server == "" {
		return nil, errors.New("building client: no server specified")
	}
	return internal.NewClient(server), nil
}

func (app *AppContext) BuildKeysSet() (*internal.KeysSet, error) {
	keysSet := internal.KeysSet{}
	if err := app.config.Unmarshal(&keysSet); err != nil {
		err = fmt.Errorf("building keys set: %w", err)
		return nil, err
	}
	return &keysSet, nil
}
