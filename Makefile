
run:
	go run main.go

build:
	# build to bin
	go build -o bin/main main.go

web:
	GOOS=js GOARCH=wasm go build -o misc/wasm/main.wasm main.go; 
