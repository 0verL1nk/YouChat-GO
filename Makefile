.PHONY: gen
gen:
	cwgo server \
	-type HTTP \
	-service core \
	-module core \
	-idl idl/${svc}.proto \
	-pass -unset_omitempty \
	-template ./template 
