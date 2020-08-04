BUILDER ?= quests-builder
DOCKER_IMAGE ?= waseemhassan/quests
# Default value "dev"
DOCKER_TAG ?= "dev"
REPOSITORY = ${DOCKER_IMAGE}:${DOCKER_TAG}

builder-image:
	@docker build --network host -t "${BUILDER}" -f build/package/Dockerfile.build .

binary-image: builder-image
	@docker run --network host --rm "${BUILDER}" | docker build --network host -t "${REPOSITORY}" -f Dockerfile.run -
