# Docker image to run shell and go utility functions in
WORKER_IMAGE = golang:1.24.2-alpine3.20
# Docker image to generate OAS3 specs
OAS3_GENERATOR_DOCKER_IMAGE = openapitools/openapi-generator-cli:latest-release

.PHONY: swag2
swag2:
	@docker run --rm -v $(PWD):/work $(WORKER_IMAGE) \
	  sh -c "apk update && apk add --no-cache git; go install github.com/swaggo/swag/cmd/swag@latest && (cd /work; swag init --output openapi)"

# Generate OAS3 from swaggo/swag output since that project doesn't support it
# TODO: Remove this if V3 spec is ever returned from that project
.PHONY: swag
swag: swag2
	@echo "[OAS3] Converting Swagger 2-to-3 (yaml)"
	@docker run --rm -v $(PWD)/openapi:/work $(OAS3_GENERATOR_DOCKER_IMAGE) \
	  generate -i /work/swagger.yaml -o /work/v3 -g openapi-yaml --minimal-update
	@docker run --rm -v $(PWD)/openapi/v3:/work $(WORKER_IMAGE) \
	  sh -c "rm -rf /work/.openapi-generator"
	@echo "[OAS3] Copying openapi-generator-ignore (json)"
	@docker run --rm -v $(PWD)/openapi/v3:/work $(WORKER_IMAGE) \
	  sh -c "cp -f /work/.openapi-generator-ignore /work/openapi"
	@echo "[OAS3] Converting Swagger 2-to-3 (json)"
	@docker run --rm -v $(PWD)/openapi:/work $(OAS3_GENERATOR_DOCKER_IMAGE) \
	  generate -s -i /work/swagger.json -o /work/v3/openapi -g openapi --minimal-update
	@echo "[OAS3] Cleaning up generated files"
	@docker run --rm -v $(PWD)/openapi/v3:/work $(WORKER_IMAGE) \
	  sh -c "mv -f /work/openapi/openapi.json /work ; mv -f /work/openapi/openapi.yaml /work ; rm -rf /work/openapi"
