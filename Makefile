build:
	rm -rf RELEASE/

	GOOS=linux GOARCH=arm64;go build -o RELEASE/osmonitor_arm64
	GOOS=linux GOARCH=x86_64;go build -o RELEASE/osmonitor_x86_64 main.go
	./copy-service-file.sh
	cp install.sh RELEASE/install.sh

.PHONY: build