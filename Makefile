test-standard:
	genny -in=template/slices.go -out=testdata/standard/standard.go -pkg=standard gen "ValueType=int"
	cd testdata/standard; go test

test-comparable:
	genny -in=template/slices-comparable.go -out=testdata/comparable/comparable.go -pkg=comparable gen "ValueType=int"
	cd testdata/comparable; go test

test-small:
	genny -in=template/slices-small.go -out=testdata/small/small.go -pkg=small gen "ValueType=int"
	cd testdata/small; go test

test-comparable-small:
	genny -in=template/slices-comparable-small.go -out=testdata/comparablesmall/comparablesmall.go -pkg=comparablesmall gen "ValueType=int"
	cd testdata/comparablesmall; go test

test: test-standard test-comparable test-small test-comparable-small

all: test

.PHONY: test test-standard test-comparable test-small test-comparable-small
