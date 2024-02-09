imageVersion := latest

protocol:
	cd proto && ./gen.sh

build:
	docker build --tag localhost:5001/beam-blobcache:$(imageVersion) .
	docker push localhost:5001/beam-blobcache:$(imageVersion)

package-chart:
	helm package --dependency-update deploy/charts/blobcache
