docker-build: api-server flow action cli
	@echo "Docker build done"

flow:
	docker build -t harbor.ym/devops/echoer-flow:v0.1.1 -f Dockerfile.flow .
	docker push harbor.ym/devops/echoer-flow:v0.1.0

action:
	docker build -t harbor.ym/devops/echoer-action:v0.1.3 -f Dockerfile.action .
	docker push harbor.ym/devops/echoer-action:v0.1.3

api-server:
	docker build -t harbor.ym/devops/echoer-api:v0.1.0 -f Dockerfile.api .
	docker push harbor.ym/devops/echoer-api:v0.1.0

cli:
	docker build -t harbor.ym/devops/echoer-cli:v0.1.0 -f Dockerfile.cli .
	docker push harbor.ym/devops/echoer-cli:v0.1.0