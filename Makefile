.PHONY: app lint test image

app:
	docker buildx bake image

lint:
	gofumpt -l -w .
	golangci-lint run

test:
	go test -v ./...

image: lint test app
