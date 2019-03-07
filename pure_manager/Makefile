.PHONY: build
build: bin/server bin/client

proto/api.pb.go: proto/api.proto
	protoc --go_out=plugins=grpc:. proto/api.proto

.PHONY: bin/client
bin/client: proto/api.pb.go
	go build -o bin/client ./client

.PHONY: bin/server
bin/server: proto/api.pb.go
	go build -o bin/server ./server

.PHONY: serve
serve: bin/server
	bin/server
