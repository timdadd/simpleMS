FROM golang:alpine AS builder
#RUN apk add --no-cache ca-certificates git && \
#      wget -qO/go/bin/dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && \
#      chmod +x /go/bin/dep

ENV PROJECT simplems/services/routeguide
WORKDIR /go/src/$PROJECT

# restore dependencies
#go mod tidy
COPY . .
RUN ls -la /
RUN go build -o /routeguide .

FROM alpine AS release
#RUN apk add --no-cache ca-certificates
#RUN GRPC_HEALTH_PROBE_VERSION=v0.2.0 && \
#    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
#    chmod +x /bin/grpc_health_probe
WORKDIR /routeguide
COPY --from=builder /routeguide ./server
#copy defaultConfig.yaml .
# Check the config name in main()
COPY cfg/dockerConfig.yaml ./cfg/routeguide.yaml
#COPY products.json .
EXPOSE 10000
ENTRYPOINT ["/routeguide/server"]
