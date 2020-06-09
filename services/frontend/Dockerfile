FROM golang:alpine AS builder
#RUN apk add --no-cache ca-certificates git && \
#      wget -qO/go/bin/dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && \
#      chmod +x /go/bin/dep

ENV PROJECT simplems/src/frontend
WORKDIR /go/src/$PROJECT

COPY . .
RUN ls -la /
RUN go install .

FROM alpine as release
RUN apk add --no-cache ca-certificates \
    busybox-extras net-tools bind-tools
WORKDIR /frontend
COPY --from=builder /go/bin/frontend /frontend/server
#COPY ./defaultConfig.yaml .
COPY cfg/dockerConfig.yaml ./frontend.yaml
COPY ./templates ./templates
COPY cfg/dockerConfig.yaml ./cfg/frontend.yaml
#COPY ./static ./static
EXPOSE 8080
ENTRYPOINT ["/frontend/server"]
