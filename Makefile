build:
	GOOS=linux;GOARCH=arm64;go build -o RELEASE/arm64/osmonitor main.go
	GOOS=linux;GOARCH=x86_64;go build -o RELEASE/x86_64/osmonitorx86 main.go
	./copy-service-file.sh
	cp install.sh RELEASE/install.sh

.PHONY: build