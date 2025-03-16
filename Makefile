.PHONY: gen
gen:
	cwgo server \
	-type HTTP \
	-service core \
	-module core \
	-idl idl/${svc}.proto \
	-template https://github.com/0verL1nk/cwgo-template.git \
	-branch hertz
