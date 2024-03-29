module github.com/lonelypale/goutils

go 1.13

require (
	github.com/BurntSushi/toml v1.2.1
	github.com/allegro/bigcache v1.2.1
	github.com/bwmarrin/snowflake v0.3.1-0.20190412223032-c09e69ae5993
	github.com/eclipse/paho.mqtt.golang v1.3.2
	github.com/fatih/color v1.13.0
	github.com/gin-contrib/sessions v0.0.3
	github.com/gin-contrib/static v0.0.0-20200916080430-d45d9a37d28e
	github.com/gin-gonic/gin v1.9.1
	github.com/go-playground/locales v0.14.1
	github.com/go-playground/universal-translator v0.18.1
	github.com/go-playground/validator/v10 v10.14.0
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-spring/spring-boot v1.0.5
	github.com/go-spring/spring-core v1.0.5
	github.com/go-spring/spring-gin v1.0.5
	github.com/go-spring/spring-logger v1.0.5
	github.com/go-spring/spring-utils v1.0.5
	github.com/go-spring/spring-web v1.0.5
	github.com/go-spring/starter-web v1.0.5
	github.com/golang-jwt/jwt/v4 v4.5.0
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.5.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/ipfs-cluster/ipfs-cluster v1.0.6
	github.com/ipfs/go-ipfs-api v0.3.0
	github.com/jpillora/backoff v1.0.0
	github.com/json-iterator/go v1.1.12
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/minio/sio v0.3.1
	github.com/mitchellh/go-ps v1.0.0
	github.com/multiformats/go-multiaddr v0.8.0
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/panjf2000/ants/v2 v2.8.1
	github.com/pkg/errors v0.9.1
	github.com/shirou/gopsutil/v3 v3.23.3
	github.com/sirupsen/logrus v1.8.1
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/stretchr/testify v1.8.3
	github.com/syndtr/goleveldb v1.0.0
	github.com/vmihailenco/msgpack/v5 v5.3.5
	go.mongodb.org/mongo-driver v1.11.3
	golang.org/x/crypto v0.17.0
	golang.org/x/image v0.6.0 // indirect
	google.golang.org/grpc v1.45.0
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.29.1
