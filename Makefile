BIN = euchre_tracker.bin

run:
	go run ./cmd/main.go

build:
	go  build -o ./bin/${BIN} ./cmd/

.PHONY: clean

clean:
	rm ./bin/${BIN}