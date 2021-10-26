

### Dependencies:
 1. Golang: https://golang.org/
 2. Echo: https://echo.labstack.com/guide/installation/
 3. Gorm: https://gorm.io/docs/models.html
 4. Twilio: https://github.com/twilio/twilio-go
 5. Kafka: https://github.com/lovoo/goka + https://github.com/Shopify/sarama 
 6. Postgre: https://github.com/go-pg/pg + https://github.com/fergusstrange/embedded-postgres
 7. Impala: https://github.com/bippio/go-impala + https://github.com/cloudera/impyla
 8. Kudu: https://kudu.apache.org/docs/developing.html
 9. Protobuf: https://developers.google.com/protocol-buffers/docs/gotutorial
 10. Redis: https://github.com/go-redis/redis
 11. Swagger: https://github.com/swaggo/swag + https://github.com/swaggo/swag
### Services:
 1. `core-gateway`
 2. `core-service`
 3. `core-account`
 4. `core-transfer`
 5. `core-search`
 6. `consumer` --- `python/java/golang`
 7. `messaging` -- `python/golang`
 8. `analytics` -- `python`
 9. `reporting` -- `python`
### Local run:
 1. Generate proto.go: ``protoc --proto_path=. --go_out=:. --grpc-go_out=:. proto/wallet.proto``
    Generate proto.py: ``python -m grpc_tools.protoc -Iproto --python_out=proto --grpc_python_out=proto/ proto/wallet_bridge.proto``
 2. Move to directory: `cd $GOPATH/$service-dir`
 3. Run go file: `go run $service-file.go`
### Docker run:
 * Link: https://hub.docker.com/_/golang

 1. Mode to directory: `cd $GOPATH/$service-dir`
 2. Docker build: `docker build -t $app-service-name`
 3. Docker run: `docker run -it --rm --name $app-running-service $app_service-name`
### Structure:
 ```
external(http)-->core-gateway--(grpc)-->core-service--(grpc)-->internal-services--(broker)-->consumer-services-->
```
#### Sequence:
 ![High level calling sequence] (uml/wallet-sequence.png)