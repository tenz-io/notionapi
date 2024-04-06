.PHONY: dep
dep:
	go mod tidy -v


.PHONY: test
test:
	go test ./... -cover -v

.PHONY: gci
gci:
	gci write -s standard -s default -s "prefix(github.com)" -s "prefix(github.com/tenz-io/gokit)" --skip-generated *