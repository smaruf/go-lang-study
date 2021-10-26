module core-service

go 1.16

require (
	common v0.0.0-00010101000000-000000000000
	github.com/confluentinc/confluent-kafka-go v1.7.0
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/labstack/echo-contrib v0.11.0
	github.com/labstack/echo/v4 v4.5.0
	golang.org/x/net v0.0.0-20210428140749-89ef3d95e781 // indirect
	google.golang.org/grpc v1.40.0
	gopkg.in/yaml.v2 v2.4.0
)

replace common => ../common
