# [protbuf](https://en.wikipedia.org/wiki/Protocol_Buffers)

## A quick history
Google's Remote Procedure Call (gRPC) protocol buffers started life in google and then was offered for general use
by donating to the Cloud Native Computing Foundation.  The code was moved from the 
[original google repository](https://code.google.com/p/protobuf/) to [github](https://github.com/protocolbuffers/protobuf).  However
for Go the github version (v1) has been replaced by (v2) [google.golang.org/protobuf](https://pkg.go.dev/mod/google.golang.org/protobuf)

Read about the [launch of api v2](https://blog.golang.org/protobuf-apiv2)

Read the [Google Tutorial](https://developers.google.com/protocol-buffers/docs/gotutorial)

Add the [Jetbrains plugin for Protobuf](https://plugins.jetbrains.com/plugin/8277-protobuf-support)

Check the latest [readme on github](https://pkg.go.dev/mod/google.golang.org/protobuf)

## GO support for protocol buffers

You need two things:
This project is comprised of two components:

* **Code generator**: The [`protoc-gen-go`](https://pkg.go.dev/google.golang.org/protobuf/cmd/protoc-gen-go) 
tool is a compiler plugin to `protoc`, the protocol buffer compiler. It augments the `protoc` compiler so that
it knows [how to generate Go specific code for a
 given `.proto` file.](https://developers.google.com/protocol-buffers/docs/reference/go-generated)

* **Runtime library**: The [`protobuf`](https://pkg.go.dev/mod/google.golang.org/protobuf) module contains a set of Go packages that
form the runtime implementation of protobufs in Go. This provides the set of interfaces that
[define what a message is](https://pkg.go.dev/google.golang.org/protobuf/reflect/protoreflect) 
and functionality to [serialise a message](https://pkg.go.dev/google.golang.org/protobuf/proto)
in various formats (e.g. [wire](https://pkg.go.dev/google.golang.org/protobuf/proto),
[JSON](https://pkg.go.dev/google.golang.org/protobuf/encoding/protojson), and 
[text](https://pkg.go.dev/google.golang.org/protobuf/encoding/prototext))).

## Installing the compiler
The latest release of Protocol Buffers can be found on the [release page]
(https://github.com/protocolbuffers/protobuf/releases/latest).  This gives you protoc but not the latest GoLang
version of the plugin.

The script `genpb.sh` that needs to be called from the microservice directory should do all the steps
for making protoc work

## Mocking
Install mockgen to get mock files for testing the interface

## REST services
The 'genpb.sh' script generates REST interface using protoc plugin from
[grpc-ecosystem](https://github.com/grpc-ecosystem/grpc-gateway)

# Swagger implementation
The 'genpb.sh' script should generate SWAGGER JSON  using protoc plugin from
[grpc-ecosystem](https://github.com/grpc-ecosystem/grpc-gateway)
