.PHONY: run

run:
	go run main.go run

gcloud-login:
	gcloud auth application-default login
