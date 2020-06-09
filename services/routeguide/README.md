# Description
This is a copy of the route guide example server and client showing how to use grpc go libraries to
perform unary, client streaming, server streaming and full duplex RPCs.

Please refer to [gRPC Basics: Go](https://grpc.io/docs/tutorials/basic/go.html) for more information.

See the definition of the route guide service in pb/routeguide.proto.

## Running the service
I have refactored this to test my project structure, it all seems to be OK, test locally by running
```sh
$ ./deploy.sh local routegude
```

or to run in docker locally
```sh
$ ./deploy.sh docker routegude
```

## Running the test client
This client will work for a local deployment either with or without docker

```sh
$ cd services/routeguide/client
$ go run client.go
```

## Optional command line flags
As with the original the client takes optional command line flags. For example, the
client and server run without TLS by default.  For the server the config is now controlled
 by my simple AppConfig framework - so needs configuration in the cfg... files
 
To enable TLS when running the client:
```sh
$ cd services/routeguide/client
$ go run client.go -tls=true
```
