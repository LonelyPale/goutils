package ipfs

type Config struct {
	URI         string `value:"${ipfs.uri:=localhost:5001}"`
	ClusterID   string `value:"${ipfs.cluster_id:=}"`
	ClusterAddr string `value:"${ipfs.cluster_addr:=}"`
}

func DefaultConfig() *Config {
	return &Config{
		URI: "localhost:5001",
	}
}
