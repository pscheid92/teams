package main

import (
	"github.com/pscheid/teams/internal"
	"github.com/spf13/viper"
	"log"
	"strings"
)

func main() {
	config := loadConfig()

	monitor := buildDataMonitor(config)
	defer monitor.Close()

	jwt := buildJwtHelper(config)
	repository := internal.NewYAMLFileDataRepository(monitor)

	server := internal.NewServer(jwt, repository)
	server.InitRoutes()

	if err := server.Start(":8080"); err != nil {
		log.Fatalln(err)
	}
}

func loadConfig() *viper.Viper {
	config := viper.New()
	config.AddConfigPath(".")
	config.SetConfigName("server")
	config.SetConfigType("yaml")

	_ = config.BindEnv("jwt.secret", "JWT_SECRET")
	_ = config.BindEnv("data.path", "DATA_PATH")

	config.SetEnvPrefix("TEAMS")
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.AutomaticEnv()

	if err := config.ReadInConfig(); err != nil {
		log.Fatalln("error reading configuration file")
	}

	return config
}

func buildJwtHelper(config *viper.Viper) *internal.JwtHelper {
	sub := config.Sub("jwt")

	secret := sub.GetString("secret")
	if secret == "" {
		log.Fatalln("missing jwt secret")
	}

	jwtConfig := internal.JwtHelperConfig{
		Issuer:   sub.GetString("issuer"),
		Audience: sub.GetString("audience"),
		Leeway:   sub.GetDuration("leeway"),
		Secret:   []byte(secret),
	}
	return internal.NewJwtHelper(jwtConfig)
}

func buildDataMonitor(config *viper.Viper) *internal.DataMonitor {
	path := config.GetString("data.path")
	if path == "" {
		log.Fatalln("missing data path")
	}

	monitor, err := internal.NewDataMonitor(path)
	if err != nil {
		panic(err)
	}
	return monitor
}
