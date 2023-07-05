module frontend

go 1.14

require (
	cloud.google.com/go/storage v1.27.0
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.7.4
	github.com/sirupsen/logrus v1.6.0
	go.opencensus.io v0.24.0
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f
	google.golang.org/grpc v1.53.0
	google.golang.org/protobuf v1.28.1
	lib v0.0.0
)

replace lib v0.0.0 => ./lib
