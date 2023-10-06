
all:
	@go test -test.v ./... -v || (echo "======= TEST FAILED =======" ; false)

vet:
	go vet ./...

fmt:
	go fmt ./...
