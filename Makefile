
all:
	@go test -test.v ./... -v || (echo "======= TEST FAILED =======" ; false)

