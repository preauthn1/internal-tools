proto:
	protoc -I protobuf-config/ protobuf-config/*.proto --go_out=plugins=grpc:protobuf-config

unit-test-short:
	go test -coverprofile=.cover -short -count 1

cover:
	go tool cover -html .cover