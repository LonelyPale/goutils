# goutils

Go Utils

常用功能放在goutils包
大的模块功能放在子包内

```shell
go mod tidy
go mod download

go get -u github.com/BurntSushi/toml
go get -u go.mongodb.org/mongo-driver

go list -m -json all
```

```shell
git tag -ln
git tag -a v0.0.7 -m "goutils-v0.0.7"
git push origin --tags
```

## Build

```shell
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix cgo -o target/sysinfod_linux-arm64 cmd/sysinfo/main.go
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix cgo -o target/sysinfod_linux-amd64 cmd/sysinfo/main.go

```

## Linux

```shell
cp target/sysinfod_linux-arm64 /usr/local/bin/sysinfod
cp target/sysinfod_linux-amd64 /usr/local/bin/sysinfod

cp sysinfo.service /lib/systemd/system
chmod +711 /usr/local/bin/sysinfod

systemctl daemon-reload
systemctl start sysinfo
systemctl status sysinfo

#失败时查看系统日志
tail -f /var/log/messages

ps -ef |grep sysinfod
curl localhost:9999/sys-info

```

## Docker

```shell
docker build -t sysinfo:latest -f docker/sysinfo.Dockerfile . \
  --network=host \
  --build-arg "HTTP_PROXY=http://127.0.0.1:7890" \
  --build-arg "HTTPS_PROXY=http://127.0.0.1:7890"

docker-compose -p sysinfo -f docker/sysinfo.yaml up -d
curl localhost:9999/sys-info

docker exec -it sysinfo sh
docker stop sysinfo && docker rm sysinfo
docker rmi sysinfo:latest
docker builder prune

docker system df

```

## Go Test

```shell
go test -v sys/sysinfo/sysinfo_test.go -test.run Test
go test -v -test.run Test
```

## License

This project is under MIT License. See the [LICENSE](LICENSE) file for the full license text.
