module book

go 1.14

require (
	github.com/golang/protobuf v1.5.2
	github.com/sirupsen/logrus v1.6.0
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f
	google.golang.org/grpc v1.53.0
	google.golang.org/protobuf v1.28.1
	lib v0.0.0
)

replace lib v0.0.0 => ./lib
