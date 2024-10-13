
build:
	# build to bin
	go build -v -o bin/main main.go

web:
	GOOS=js GOARCH=wasm go build -o misc/wasm/main.wasm client.go; 
