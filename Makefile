.PHONY:

build:
	go build -o ./chat cmd/main.go

run: build 
	./chat

test:
	go test -v ./tests

tests-in-docker:
	sudo docker exec -it server sh
#	make test

docker-rm:
	sudo docker stop server
	sudo docker rm server
docker-rmi:
	sudo docker rmi rest-api-server
