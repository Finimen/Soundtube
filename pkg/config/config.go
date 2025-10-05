package config

type Config struct {
	Environment Environment `mapstructure:"environment"`
	Redis       Redis       `mapstructure:"redis"`
	Repository  Repository  `mapstructure:"repository"`
	Server      Server      `mapstructure:"server"`
	Traycing    Traycing    `mapstructure:"traycing"`
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

type Repository struct {
}

type Server struct {
	Port         string `mapstructure:"port"`
	CookieSecure bool   `mapstructure:"cookie_secure"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
	IdleTimeout  int    `mapstructure:"idle_timeout"`
}

func LoadConfig() (*Config, error) {
	return nil, nil
}
