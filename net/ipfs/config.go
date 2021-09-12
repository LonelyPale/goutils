package ipfs

type Config struct {
	URI string `value:"${ipfs.uri:=localhost:5001}"`
}

func DefaultConfig() *Config {
	return &Config{
		URI: "localhost:5001",
	}
}
