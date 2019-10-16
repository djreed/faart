build_all: build_send build_receiver
	cp sender/3700send .
	cp receiver/3700recv .

build_send:
	pushd sender; make; popd

build_receiver:
	pushd receiver; make; popd