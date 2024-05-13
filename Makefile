build:
	GOARCH=amd64 GOOS=linux go build -o bin/ -ldflags "-s -w" code.cloudfoundry.org/cnbapplifecycle/cmd/builder code.cloudfoundry.org/cnbapplifecycle/cmd/launcher
	test -f bin/diego-sshd || curl -sL https://storage.googleapis.com/cf-deployment-compiled-releases/diego-2.96.0-ubuntu-jammy-1.80-20240322-160011-008934012.tgz | tar -xzO ./compiled_packages/diego-sshd.tgz | tar -xzO ./diego-sshd > bin/diego-sshd && chmod +x bin/diego-sshd
	test -f bin/healthcheck || curl -sL https://storage.googleapis.com/cf-deployment-compiled-releases/diego-2.96.0-ubuntu-jammy-1.80-20240322-160011-008934012.tgz | tar -xzO ./compiled_packages/healthcheck.tgz | tar -xzO ./healthcheck > bin/healthcheck && chmod +x bin/healthcheck

test:
	go test -v -count=1 ./...

integration: build
	INCLUDE_INTEGRATION_TESTS=true go test -v -count=1 ./integration --ginkgo.label-filter integration -ginkgo.v

package: build
	tar czf bin/cnb_app_lifecycle.tgz -C bin builder launcher diego-sshd healthcheck

.PHONY: build test integration package
