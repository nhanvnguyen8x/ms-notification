package configs

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Server    *ServerConfig
	Supabase  *SupabaseConfig
	SendPulse *SendPulseConfig
	Jwt       *JwtConfig
}

type SupabaseConfig struct {
	ApiUrl string
	ApiKey string
}

type ServerConfig struct {
	Address string
}

type JwtConfig struct {
	SecretKey string
}

type SendPulseConfig struct {
	ClientID     string
	ClientSecret string
	AuthUrl      string
	TemplateUrl  string
	EmailUrl     string
}

func NewApplicationConfig() *AppConfig {
	config := &AppConfig{}

	viper.SetConfigName("application.yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		logrus.Errorf("Read config file error: %s", err.Error())
		panic(err)
	}

	config.Supabase = &SupabaseConfig{
		ApiUrl: viper.GetString("supabase.api-url"),
		ApiKey: viper.GetString("supabase.api-key"),
	}
	config.Server = &ServerConfig{
		Address: ":" + viper.GetString("server.address"),
	}
	config.Jwt = &JwtConfig{
		SecretKey: viper.GetString("jwt.secret-key"),
	}
	config.SendPulse = &SendPulseConfig{
		ClientID:     viper.GetString("sendpulse.client-id"),
		ClientSecret: viper.GetString("sendpulse.client-secret"),
		AuthUrl:      viper.GetString("sendpulse.auth-url"),
		TemplateUrl:  viper.GetString("sendpulse.template-url"),
		EmailUrl:     viper.GetString("sendpulse.email-url"),
	}
	return config
}
