.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/ebay-find-by-category ebay-find-by-category/main.go
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/ebay-find-by-keyword ebay-find-by-keyword/main.go
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/ebay-find-advanced ebay-find-advanced/main.go
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/ebay-find-by-product ebay-find-by-product/main.go
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/ebay-find-in-ebay-stores ebay-find-in-ebay-stores/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
