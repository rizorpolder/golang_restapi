PROTO_ROOT=.
DOCKER_IMAGE=proto-builder
.PHONY: docker-build gen clean

docker-build:
	docker build -t $(DOCKER_IMAGE) .
gen: docker-build
	docker run --rm -v $(abspath $(PROTO_ROOT)):/app $(DOCKER_IMAGE) bash /app/generate.sh
clean:
	find account pagination auth -type d -name go -exec rm -rf {} \;