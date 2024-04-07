build:
	rm -rf RELEASE/

	env GOOS=linux GOARCH=arm64 go build -o RELEASE/osmonitor_arm64 main.go
	env GOOS=linux GOARCH=amd64 go build -o RELEASE/osmonitor_x86_64 main.go
	./copy-service-file.sh
	cp install.sh RELEASE/install.sh

.PHONY: build