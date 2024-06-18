lint:
	golangci-lint run --timeout=10m -v ./...


swag-api:
	swag init --dir ./api --output ./api/docs --parseDependency --parseInternal

gen-proto:
	protoc -I=. -I=vendor -I=${GOPATH}/src core/proto/*.proto --go_out=. --go-grpc_out=.