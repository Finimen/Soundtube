package config

type Config struct {
	Environment Environment `mapstructure:"environment"`
	Redis       Redis       `mapstructure:"redis"`
	Repository  Repository  `mapstructure:"repository"`
	Server      Server      `mapstructure:"server"`
	Traycing    Traycing    `mapstructure:"traycing"`
	Token       Token       `mapstructure:"token"`
	Email       Email       `mapstructure:"email"`
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
	Path string
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
	return nil, nil
}
