.PHONY: run

run:
	go build -o shakesearch
	./shakesearch

test: go-test frontend-test

go-test:
	go test -v

frontend-test:
	docker build -t shakesearch-test -f Dockerfile.test .
	docker run --rm shakesearch-test