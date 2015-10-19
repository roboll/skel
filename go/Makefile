###############################################################################
# see BUILD.md
###############################################################################
OWNER     := {{ .Owner }}
REPO      := {{ .Repo }}

PROJECT   := $(OWNER)/$(REPO)
VERSION   := $(shell git describe --tags)
REG       :=

GOOS      := linux
GOARCH    := amd64

PRE_RELEASE := tag clean-repo test

.PHONY:  build docker release
build:   {{ .Repo }}-linux-amd64 docs.tar.gz
docker:  build docker-build-root
release: $(PRE_RELEASE) docker docker-push-root docs.tar.gz-gh-release

###############################################################################
# pre-release - test and validation
###############################################################################
.PHONY: test
test: ; go test ./...

###############################################################################
# build / release
###############################################################################
GOBUILD := GO15VENDOREXPERIMENT=1 GOOS=$(GOOS) GOARCH=$(GOARCH) go build
GOARGS  := -a -tags netgo -ldflags '-s -w -X main.release=$(VERSION)'

%.tar.gz: %                ;tar czf $*.tar.gz -C $* .
%-$(GOOS)-$(GOARCH): $(*D) ;$(GOBUILD) $(GOARGS) -o $@ ./$(*D)

.PHONY: docker-build-root docker-build-%
docker-build-root: ;docker build -t $(REG)$(PROJECT):$(VERSION) ./
docker-build-%:    ;docker build -t $(REG)$(PROJECT)-$*:$(VERSION) ./$*

.PHONY: docker-push-root docker-push-%
docker-push-root: docker-build-root ;docker push $(REG)$(PROJECT):$(VERSION)
docker-push-%:    docker-build-%    ;docker push $(REG)$(PROJECT)-$*:$(VERSION)

###############################################################################
# github-release - upload a binary release to github releases
#
# requirements:
# - the checked out revision be a pushed tag
# - a github api token (GITHUB_TOKEN)
###############################################################################
API    = https://api.github.com/repos/$(OWNER)/$(REPO)
UPLOAD = https://uploads.github.com/repos/$(OWNER)/$(REPO)

.PHONY: create-gh-release %-gh-release gh-token
create-gh-release: tag clean-repo gh-token
	$(info Creating Github Release)
	@curl -s -XPOST -H "Authorization: token $(GITHUB_TOKEN)" \
		$(API)/releases -d '{ "tag_name": "$(VERSION)" }' > /dev/null

%-gh-release: tag clean-repo gh-token create-gh-release
	$(info Uploading Release Artifact $* to Github)
	@curl -s -H "Authorization: token $(GITHUB_TOKEN)" \
		$(API)/releases/tags/$(VERSION) |\
	python -c "import json,sys;obj=json.load(sys.stdin);print obj['id']" |\
	curl -s -XPOST \
		-H "Authorization: token $(GITHUB_TOKEN)" \
		-H "Content-Type: application/octet-stream" \
		$(UPLOAD)/releases/$$(xargs )/assets?name=$* \
		--data-binary @$* > /dev/null

gh-token:
ifndef GITHUB_TOKEN
	$(error $GITHUB_TOKEN not set)
endif

###############################################################################
# utility
###############################################################################
.PHONY: tag clean-repo
tag:
	@echo Ensuring checkout is a tag.
	@git describe --tags --exact-match HEAD > /dev/null
clean-repo:
	@echo Ensuring repository is clean.
	@git diff --exit-code > /dev/null
	@git diff --cached --exit-code > /dev/null

###############################################################################
.DEFAULT_GOAL := info
info:
	@echo project: $(PROJECT)
	@echo version: $(VERSION)
