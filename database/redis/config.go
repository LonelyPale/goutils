package redis

type Config struct {
	Endpoint string `value:"${redis.endpoint:=localhost:6379}"`
	Password string `value:"${redis.password:=}"`
	Database int    `value:"${redis.database:=0}"`
}

func DefaultConfig() *Config {
	return &Config{
		Endpoint: "localhost:6379",
		Password: "",
		Database: 0,
	}
}
