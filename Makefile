.PHONY: ci test install

ci:
	@echo "Running rocker build --var NOTMUCH_VERSION=${NOTMUCH_VERSION} --var BUILD_NUMBER=${TRAVIS_BUILD_NUMBER} --var COMMIT=${COMMIT}"
	@rocker build --var NOTMUCH_VERSION=${NOTMUCH_VERSION} --var BUILD_NUMBER=${TRAVIS_BUILD_NUMBER} --var COMMIT=${COMMIT} ${ROCKER_ARGS}

test: fixtures/database-v1
	@go test -v -race -cover $(shell go list ./... | grep -v /vendor/)
	@rm -rf fixtures/database-v1

install:
	@go install $(shell go list ./... | grep -v /vendor/)

fixtures/database-v1: fixtures/database-v1.tar.xz
	@tar -C fixtures -xf fixtures/database-v1.tar.xz

fixtures/database-v1.tar.xz:
	@mkdir -p fixtures
	@curl --insecure -Lo fixtures/database-v1.tar.xz https://notmuchmail.org/releases/test-databases/database-v1.tar.xz

server/grpc/gmuch.pb.go:
	@protoc server/grpc/gmuch.proto --go_out=plugins=grpc:.
