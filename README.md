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

```

## Go Test

```shell
go test -v sys/sysinfo/sysinfo_test.go -test.run Test
go test -v -test.run Test
```

## License

This project is under MIT License. See the [LICENSE](LICENSE) file for the full license text.
