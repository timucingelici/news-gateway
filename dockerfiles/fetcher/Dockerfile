# 2-STEPS GO BUILD
# BUILD
FROM golang:latest as builder

RUN mkdir /go/src/build
COPY . /go/src/build
WORKDIR /go/src/build

RUN apt-get -y update && apt-get -y install go-dep

RUN dep ensure

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o fetcher cmd/fetcher/fetcher.go

# COPY INTO A SMALLER CONTAINER
FROM scratch AS final

COPY --from=builder /go/src/build/fetcher /
WORKDIR /

ENTRYPOINT ["/fetcher"]