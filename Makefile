protocol:
	cd proto && ./gen.sh

build:
	okteto build --file ./Dockerfile --tag okteto.dev/cache:latest --target build