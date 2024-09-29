.PHONY: test e2e-test cover gofmt gofmt-fix header-check clean tar.gz docker-push release docker-push-all flannel-git

# Default tag and architecture. Can be overridden
TAG?=$(shell git describe --tags --dirty --always)
ifeq ($(TAG),)
	TAG=latest
endif

#ifeq ($(findstring dirty,$(TAG)), dirty)
#    TAG=latest
#endif

### BUILDING
#debug:
#	@echo $(TAG)
clean:
	rm -f mt

mt: $(shell find . -type f  -name '*.go')
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o mt \
	  -ldflags '-s -w -X "github.com/daddvted/mutant/mutant.Version=$(TAG)" -extldflags "-static"'
