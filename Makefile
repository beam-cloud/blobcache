imageVersion := latest

protocol:
	cd proto && ./gen.sh

build:
	docker build . -f Dockerfile -t localhost:5000/beam-blobcache:$(imageVersion)
	docker push localhost:5000/beam-blobcache:latest
