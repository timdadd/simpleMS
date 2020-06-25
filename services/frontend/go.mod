module frontend

go 1.14

require (
	cloud.google.com/go/storage v1.8.0
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.7.4
	github.com/sirupsen/logrus v1.6.0
	go.opencensus.io v0.22.3
	google.golang.org/genproto v0.0.0-20200608115520-7c474a2e3482
	google.golang.org/grpc v1.29.1
	google.golang.org/protobuf v1.24.0
	lib v0.0.0
)

replace lib v0.0.0 => ./lib
