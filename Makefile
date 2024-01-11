protocol:
	cd proto && ./gen.sh

build:
	okteto build --file ./Dockerfile --tag okteto.dev/blobcache:latest --target build