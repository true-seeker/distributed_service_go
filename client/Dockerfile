FROM golang:alpine
RUN apk update && apk add --no-cache make protobuf-dev git

WORKDIR /app

COPY . ./client/

WORKDIR /app/client

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go mod download

RUN mkdir -p backpackTaskGRPC
