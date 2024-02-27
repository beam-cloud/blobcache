imageVersion := latest

protocol:
	cd proto && ./gen.sh

build:
	docker build . -f Dockerfile -t localhost:5001/beam-blobcache:$(imageVersion)
	docker push localhost:5001/beam-blobcache:latest

start:
	cd hack && okteto up --file okteto.yml

stop:
	cd hack && okteto down --file okteto.yml