

.PHONY: gen
gen:
	go install .
	buf generate


.PHONY: gen-connect
gen-connect:
	go install .
	buf generate --template buf.gen.connect.yaml

.PHONY: gen-override
gen-override:
	go install .
	buf generate --template buf.gen.override.yaml
