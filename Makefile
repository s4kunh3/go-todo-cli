build:
	go build -o todocli main.go utils.go
.PHONY: build run clean
run:
	./todocli
clean:
	rm -f todocli
.PHONY: build run clean