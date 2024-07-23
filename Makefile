build:
	GOARCH=amd64 GOOS=linux go build -o bin/ -ldflags "-s -w" code.cloudfoundry.org/cnbapplifecycle/cmd/builder code.cloudfoundry.org/cnbapplifecycle/cmd/launcher

test:
	go test -count=1 ./...

integration: build
	INCLUDE_INTEGRATION_TESTS=true go test -v -count=1 ./integration --ginkgo.label-filter integration -ginkgo.v

package: build
	tar czf bin/cnb_app_lifecycle.tgz -C bin builder launcher diego-sshd healthcheck

.PHONY: build test integration package
