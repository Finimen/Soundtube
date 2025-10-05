package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Environment         Environment         `mapstructure:"environment"`
	Redis               Redis               `mapstructure:"redis"`
	Database            Database            `mapstructure:"database"`
	DatabaseConnections DatabaseConnections `mapstructure:"database_connections"`
	Server              Server              `mapstructure:"server"`
	Traycing            Traycing            `mapstructure:"traycing"`
	Token               Token               `mapstructure:"token"`
	Email               Email               `mapstructure:"email"`
}

type Environment struct {
	Current string `mapstructure:"current"`
}

type Traycing struct {
	Enabled     bool   `mapstructure:"enabled"`
	ServiceName string `mapstructure:"service_name"`
	Endpoint    string `mapstructure:"endpoint"`
}

type Redis struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type DatabaseConnections struct {
	MaxOpenConns    int `mapstructure:"max_open_conns"`
	MaxIdleConns    int `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int `mapstructure:"max_life_time"`
	ConnMaxIdleTime int `mapstructure:"max_idle_time"`
}

type Server struct {
	Port         string `mapstructure:"port"`
	CookieSecure bool   `mapstructure:"cookie_secure"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
	IdleTimeout  int    `mapstructure:"idle_timeout"`
}

type Token struct {
	JwtKey string `mapstructure:"jwt_key"`
	Exp    int    `mapstructure:"exp"`
}

type Email struct {
	From string `mapstructure:"from"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("dev")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	viper.AutomaticEnv()

	viper.SetDefault("environment.current", "development")
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("ratelimit.maxrequests", 100)
	viper.SetDefault("ratelimit.window", time.Minute)

	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	if len(config.Token.JwtKey) < 16 {
		panic("invalid jwt key")
	}

	return nil, nil
}
