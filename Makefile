

.PHONY: gen
gen:
	go install .
	buf generate