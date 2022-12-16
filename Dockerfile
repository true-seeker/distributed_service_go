FROM golang:alpine
RUN apk update && apk add --no-cache make protobuf-dev git

WORKDIR /app

COPY . ./

WORKDIR /app/server
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go mod download

RUN mkdir -p backpackTaskGRPC

WORKDIR /app/gRPC
RUN protoc --go_out=../server/backpackTaskGRPC/ --go_opt=paths=source_relative --go-grpc_out=../server/backpackTaskGRPC/ --go-grpc_opt=paths=source_relative grpc.proto

EXPOSE 3333 9876

WORKDIR /app/server
RUN go build -o /docker-server

CMD [ "/docker-server" ]
