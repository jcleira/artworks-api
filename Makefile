help:
	@echo "Use make deploy for artworks-api pre-production deployment"

deploy: build docker_build docker_push clean

# GO binary building. we do build it for ARM as it's our preproduction deployment.
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w' -o artworks-api .

# Build docker image with name:tag from DockerHub.
docker_build:
	docker build -t jcorral/artworks-api:latest .

docker_push:
	docker push jcorral/artworks-api:latest

clean:
	rm -rf artworks-api
