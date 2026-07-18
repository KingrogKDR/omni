
build:
	go build -o omni ./cmd/main.go

run: build
	./omni

clean:
	rm ./omni