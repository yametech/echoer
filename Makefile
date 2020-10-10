flow:
	docker build -t yametech/echoer-flow:v0.1.0 -f Dockerfile.flow .
	docker push yametech/echoer-flow:v0.1.0

action:
	docker build -t yametech/echoer-action:v0.1.0 -f Dockerfile.action .
	docker push yametech/echoer-action:v0.1.0

api:
	docker build -t yametech/echoer-api:v0.1.0 -f Dockerfile.api .
	docker push yametech/echoer-api:v0.1.0

cli:
	docker build -t yametech/echoer-cli:v0.1.0 -f Dockerfile.cli .
	docker push yametech/echoer-cli:v0.1.0

docker-build: flow action api cli
	@echo "Docker build done"