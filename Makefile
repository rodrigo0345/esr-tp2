
build:
	# build to bin
	go build -v -o bin/main main.go

web:
	GOOS=js GOARCH=wasm go build -o misc/wasm/main.wasm client.go; 

core:
	scp -r * core@192.168.56.101:/home/core/projeto2
