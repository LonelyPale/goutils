module github.com/LonelyPale/goutils

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/allegro/bigcache v1.2.1
	github.com/bwmarrin/snowflake v0.3.1-0.20190412223032-c09e69ae5993
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.3.0
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/google/uuid v1.1.1
	github.com/json-iterator/go v1.1.9
	github.com/mitchellh/go-ps v1.0.0
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.2
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/vmihailenco/msgpack/v5 v5.0.0-alpha.2
	go.mongodb.org/mongo-driver v1.3.2
	golang.org/x/image v0.0.0-20200927104501-e162460cd6b5 // indirect
	google.golang.org/grpc v1.29.1
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.29.1
