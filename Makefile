.PHONY: run

run:
	go run targets/server/main.go run

gcloud-login:
	gcloud auth application-default login

client-login:
	go run targets/client/main.go login
