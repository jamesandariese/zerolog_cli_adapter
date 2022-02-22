.PHONY: all always-build act-clean

# a portable dirname + realpath
DIRNAME := $(shell find .. -maxdepth 1 -mindepth 1 -samefile . -execdir echo {} \;)

all: version.txt
clean: act-clean
always-build:

version.txt: version.txt.tmpl always-build
	gomplate < version.txt.tmpl > version.txt
	git add version.txt

prepare-release: version.txt

test: always-build
	go test

.act-image: test/Dockerfile
	docker build -t act-image-$(DIRNAME) --iidfile "$@" -f test/Dockerfile .

# NOTE: This is unreliable still.
act-test: .act-image $(wildcard .github/workflow/*.yml)
	act -r -P ubuntu-latest=act-image-$(DIRNAME) push

act-clean:
	[ -f .act-image ] && docker rmi $$(cat .act-image)
