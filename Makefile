test-timsort:
	genny -in=template-timsort/slices.go -out=testdata/timsort/slices.go -pkg=standard gen "ValueType=int"
	cd testdata/timsort; go test

test-comparable-timsort:
	genny -in=template-comparable-timsort/slices.go -out=testdata/comparabletimsort/slices.go -pkg=comparable gen "ValueType=int"
	cd testdata/comparabletimsort; go test

test-standard:
	genny -in=template/slices.go -out=testdata/standard/slices.go -pkg=small gen "ValueType=int"
	cd testdata/standard; go test

test-comparable:
	genny -in=template-comparable/slices.go -out=testdata/comparable/slices.go -pkg=comparablesmall gen "ValueType=int"
	cd testdata/comparable; go test

test: test-standard test-comparable test-timsort test-comparable-timsort

install:
	go get github.com/cheekybits/genny

all: test

.PHONY: test test-standard test-comparable test-timsort test-comparable-timsort
