module core-service

go 1.16

require (
	github.com/labstack/echo-contrib v0.11.0
	github.com/labstack/echo/v4 v4.5.0
	golang.org/x/net v0.7.0 // indirect
	google.golang.org/grpc v1.53.0
	gopkg.in/yaml.v2 v2.4.0
)

replace common => ../common
