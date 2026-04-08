build:
	GOARCH=amd64 GOOS=linux go build -o bin/ -ldflags "-s -w" code.cloudfoundry.org/cnbapplifecycle/cmd/builder code.cloudfoundry.org/cnbapplifecycle/cmd/launcher

test:
	go test -count=1 ./...

lint:
	go tool golangci-lint fmt ./... --diff
	go tool golangci-lint run ./... --default none --enable govet,staticcheck,unused

integration: build
	INCLUDE_INTEGRATION_TESTS=true go test -v -count=1 ./integration --ginkgo.label-filter integration -ginkgo.v

copy-binaries: 
	mkdir -p bin/diego/
	CONTAINER=$$(docker create ghcr.io/cloudfoundry/k8s/fileserver:latest); \
	docker cp $$CONTAINER:fileserver/v1/static/cnb_app_lifecycle/cnb_app_lifecycle.tgz bin/diego/; \
	docker rm $$CONTAINER
	tar -xf bin/diego/cnb_app_lifecycle.tgz -C bin/diego/
	mv bin/diego/diego-sshd bin/
	mv bin/diego/healthcheck bin/
	mv bin/diego/cf-pcap bin/
	rm -rf bin/diego/

package: build
	tar czf bin/cnb_app_lifecycle.tgz -C bin builder launcher diego-sshd healthcheck cf-pcap

.PHONY: build test integration package
