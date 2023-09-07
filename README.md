# Compliance HUB Service

## Local environment

**To start local development you need:**

* installed docker locally and postgres service (uou can start DB from Infra repository)
* copy content of `example.env` to new file `.env` in root directory. Modify it according to your needs
* start docker image `docker-compose.yml` with command:
  `docker-compose -p pwg-core -f docker-compose.yml up -d`

--

* to regenerate swagger docs run command:
  `swag init --parseDependency --parseInternal --parseDepth 1 --dir ./cmd,./internal/http/rest && swag fmt`
    * to install swaggo you need run command:
    * `go install github.com/swaggo/swag/cmd/swag@latest`
  
--
  
* to regenerate grpc docs from proto3 code run command:
  `protoc -I=./ --go_out=pkg/grpc/ --go-grpc_out=pkg/grpc/ internal/grpc/proto/*.proto --experimental_allow_proto3_optional`
  * to install protobuf compiler see https://grpc.io/docs/protoc-installation/ or for ArchLinux use `yay -S protobuf`
  * to install protobuf GO plugin run command
  `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`