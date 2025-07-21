GOPATH:=$(shell go env GOPATH)

.PHONY: gen
gen:
	clear && buf generate

.PHONY: dep-update
dep-update:
	clear && buf dep update

.PHONY: dep-prune
buf dep prune:
	clear && buf dep prune

.PHONY: push
push:
	clear && buf push

.PHONY: gen-common
gen-common:
	clear && buf generate buf.build/moriba-sl/ose