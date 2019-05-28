test:
	genny -in=template/slices.go -out=testdata/standard/standard.go -pkg=standard gen "ValueType=int"
	cd testdata/standard; go test
	genny -in=template/slices-comparable.go -out=testdata/comparable/comparable.go -pkg=comparable gen "ValueType=int"
	cd testdata/comparable; go test
	genny -in=template/slices-small.go -out=testdata/small/small.go -pkg=small gen "ValueType=int"
	cd testdata/small; go test
	genny -in=template/slices-comparable-small.go -out=testdata/comparablesmall/comparablesmall.go -pkg=comparablesmall gen "ValueType=int"
	cd testdata/comparable-small; go test

.PHONY: test
