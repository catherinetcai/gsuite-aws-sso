.PHONY: run

all: build-client build-server

release:
	mkdir -p release/

build-client: release
	go build -o release/client targets/client/main.go

build-server: release
	go build -o release/server targets/server/main.go

clean:
	rm -rf release/

run:
	go run targets/server/main.go run

gcloud-login:
	gcloud auth application-default login

client-login:
	go run targets/client/main.go login
