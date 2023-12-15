# TikTok Tech Immersion 23

![Tests](https://github.com/TikTokTechImmersion/assignment_demo_2023/actions/workflows/test.yml/badge.svg)

## Setup

This project uses Go version 1.20, protobuffers and kitex (ByteDance's rpc middleware framework)

```bash
# verify installation
go version
```

#### Configure Go

As this project has multiple `main.go` files, vscode might not recognize the correct workspace.

The solution is to create a `go.work` file and run the following commands to add the respective folders.

```bash
# in root
go work init
go work use http-server
go work use rpc-server
```

*References*: [here](https://stackoverflow.com/a/74106982)

#### Install Kitex

```bash
go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
go install github.com/cloudwego/thriftgo@latest
```

#### Install protobuf

```bash
go install github.com/golang/protobuf/protoc-gen-go@latest
```

#### Docker

Pull the images and start the server with `docker compose up`.

## Project Architecture

Client -> HTTP Server -> RPC Server -> Redis

## Learning Points

### Remote Procedure Call (RPC)

*References*:

- [what grpc is](https://grpc.io/docs/what-is-grpc/introduction/)

### proto

Protocol buffers (protobufs) are a way to serialize data in an efficient way, 
allowing better utilization of network resources.

It is an interface definition language (IDL) that is language neutral.
i.e From a single .proto file, you can easily generate Go, Java, Python code with the code generation features.

`idl_http.proto` defines the API request and response types.

*References*:

- [what protobufs are](https://medium.com/javarevisited/what-are-protocol-buffers-and-why-they-are-widely-used-cbcb04d378b6)
- [golang with protobufs](https://www.youtube.com/watch?v=qWN69yfRsVs)
- [official docs](https://protobuf.dev/getting-started/gotutorial/)

### kitex

RPC Framework by ByteDance.

*References*:

- [official docs](https://www.cloudwego.io/docs/kitex/getting-started/)
- [github](https://github.com/cloudwego/kitex)

### thrift

Another RPC Framework but under Apache.

`kitex_gen` contains RPC client (HTTP is client from architecture) and server code for RPC server.

Generated from `idl_rpc.thrift`.

*References*:

- [what thrift is](https://stackoverflow.com/questions/20653240/what-is-rpc-framework-and-apache-thrift)

### hertz

HTTP Framework for Go by ByteDance.

Generated from `idl_http.proto`.

``

*References*:

- [official docs](https://www.cloudwego.io/docs/hertz/)
- [github](https://github.com/cloudwego/hertz)
