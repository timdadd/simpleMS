FROM golang:alpine AS builder
#RUN apk add --no-cache ca-certificates git && \
#      wget -qO/go/bin/dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && \
#      chmod +x /go/bin/dep

ENV PROJECT simplems/services/systemservice
WORKDIR /go/src/$PROJECT

# restore dependencies
#go mod tidy
COPY . .
#RUN ls -la /
RUN go build -o /systemservice .

FROM alpine AS release
#RUN apk add --no-cache ca-certificates
#RUN GRPC_HEALTH_PROBE_VERSION=v0.2.0 && \
#    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
#    chmod +x /bin/grpc_health_probe
WORKDIR /systemservice
COPY --from=builder /systemservice ./server
#copy defaultConfig.yaml .
COPY cfg/dockerConfig.yaml ./cfg/systemservice.yaml
#COPY products.json .
EXPOSE 3550
ENTRYPOINT ["/systemservice/server"]
