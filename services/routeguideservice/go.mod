module routeguideservice

go 1.14

replace pb v0.0.0 => ./genpb

replace lib/utils v0.0.0 => ./lib/utils

require (
	github.com/stretchr/testify v1.4.0
	google.golang.org/grpc v1.29.1
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v0.0.0-20200604175613-ad51f572fd27 // indirect
	google.golang.org/protobuf v1.24.0
	lib/common v0.0.0
	lib/utils v0.0.0
	pb v0.0.0
)

replace lib/common v0.0.0 => ./lib/common
