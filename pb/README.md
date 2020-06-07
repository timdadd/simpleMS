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

* **Code generator**: The [`protoc-gen-go`](https://pkg.go.dev/google.golang.org/protobuf/cmd/protoc-gen-go) tool is a compiler plugin
 to `protoc`, the protocol buffer compiler. It augments the `protoc` compiler so that it knows [how to generate Go specific code for a
 given `.proto` file.](https://developers.google.com/protocol-buffers/docs/reference/go-generated)

* **Runtime library**: The [`protobuf`](https://pkg.go.dev/mod/google.golang.org/protobuf) module contains a set of Go packages that
form the runtime implementation of protobufs in Go. This provides the set of interfaces that
[define what a message is](https://pkg.go.dev/google.golang.org/protobuf/reflect/protoreflect) an[d fu](https://pkg.go.dev/google.golang.org/protobuf/proto)nctionality to serialize message
in various formats (e.g. [wire](https://pkg.go.dev/google.golang.org/protobuf/proto),
[JSON](https://pkg.go.dev/google.golang.org/protobuf/encoding/protojson), and 
[text](https://pkg.go.dev/google.golang.org/protobuf/encoding/prototext))).

## Installing the compiler
The compiler reads the proto file and generates GO code.  Install the compiler as follows:
THIS DIDN'T WORK
```bash
https://pkg.go.dev/google.golang.org/protobuf/encoding/prototext
```

From the [protobuf releases page](https://github.com/protocolbuffers/protobuf/releases) I downloaded the `protoc-x.y.z-linux-x86_64.zip`
file. Then I copied protoc to `~/bin`. If you need to create this directory then login and logout to get it added to
the PATH automatically.  You can also copy the content of the include folder to /usr/local/include.

Grab the protoc-gen-go plugin for the compiler:

You need to do this from within a go project
```bash
go get -u  google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go
```

The latest release of Protocol Buffers can be found on the [release page](https://github.com/protocolbuffers/protobuf/releases/latest
)  This gives you protoc but not the latest GoLang version.

Now have developed a script `genpb.sh` that needs to be run in the microservice directory and will do all the steps
for making protobug work