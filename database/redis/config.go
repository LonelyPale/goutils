package redis

type Config struct {
	Endpoint string `value:"${redis.endpoint:=localhost:6379}"`
	Password string `value:"${redis.password:=123456}"`
	Database int    `value:"${redis.database:=0}"`
}

func DefaultConfig() *Config {
	return &Config{
		Endpoint: "localhost:6379",
		Password: "123456",
		Database: 0,
	}
}
