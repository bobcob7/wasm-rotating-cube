all: wasm-rotating-cube

wasm-rotating-cube: output.wasm main.go
	go build -o wasm-rotating-cube main.go

output.wasm: bundle.go
	GOOS=js GOARCH=wasm go build -o bundle.wasm bundle.go

run: output.wasm wasm-rotating-cube
	./wasm-rotating-cube

clean:
	rm -f wasm-rotating-cube *.wasm
