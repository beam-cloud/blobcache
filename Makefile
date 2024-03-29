chartVersion := 0.1.0
imageVersion := latest

protocol:
	cd proto && ./gen.sh

build:
	docker build --tag localhost:5001/beam-blobcache:$(imageVersion) .
	docker push localhost:5001/beam-blobcache:$(imageVersion)

build-chart:
	helm package --dependency-update deploy/charts/blobcache --version $(chartVersion)

publish-chart:
	helm push beam-blobcache-chart-$(chartVersion).tgz oci://public.ecr.aws/n4e0e1y0
	rm beam-blobcache-chart-$(chartVersion).tgz
