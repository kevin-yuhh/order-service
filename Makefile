default: run

build:
	go build

run: build
	./soter-order-service

clean:
	rm soter-order-service *.log