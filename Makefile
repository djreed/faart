build_all: build_send build_receiver
	mv sender/3700send .
	mv receiver/3700recv .

build_send:
	pushd sender; make; popd

build_receiver:
	pushd receiver; make; popd