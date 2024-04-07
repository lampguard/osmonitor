build:
	go build -o RELEASE/osmonitor main.go
	./copy-service-file.sh
	cp install.sh RELEASE/install.sh

.PHONY: build