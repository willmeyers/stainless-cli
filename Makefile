BINARY_NAME=stainless

format:
	gofmt -s -w .

build:
	GOARCH=amd64 GOOS=darwin go build -o build/${BINARY_NAME}-darwin cmd/stainless/stainless.go
	GOARCH=amd64 GOOS=linux go build -o build/${BINARY_NAME}-linux cmd/stainless/stainless.go
	GOARCH=amd64 GOOS=windows go build -o build/${BINARY_NAME}-windows cmd/stainless/stainless.go

clean:
	go clean
	rm build/${BINARY_NAME}-darwin
	rm build/${BINARY_NAME}-linux
	rm build/${BINARY_NAME}-windows