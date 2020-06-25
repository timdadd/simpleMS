module book

go 1.14

require (
	github.com/golang/protobuf v1.4.2
	github.com/sirupsen/logrus v1.6.0
	google.golang.org/genproto v0.0.0-20200608115520-7c474a2e3482
	google.golang.org/grpc v1.29.1
	google.golang.org/protobuf v1.24.0
	lib v0.0.0
)

replace lib v0.0.0 => ./lib
