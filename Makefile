OUTFILE="bundle"
PROJECT_GOFILES=go.mod go.sum log packet receiver sender shared vendor Makefile
TEST_DATA=test_data

build_all: build_send build_recv move

build_send:
	pushd sender; make; popd

build_recv:
	pushd receiver; make; popd

build_all_linux: build_send_linux build_recv_linux move

build_send_linux:
	pushd sender; make build_linux; popd

build_recv_linux:
	pushd receiver; make build_linux; popd

move:
	mv sender/3700send .
	mv receiver/3700recv .

vendor:
	GO111MODULE=on go mod vendor

bundle:
	tar -czvf $(OUTFILE).tar.gz $(PROJECT_GOFILES) $(TEST_DATA) README.md

copy:
	scp -r $(OUTFILE).tar.gz reedda@gordon.ccs.neu.edu:/home/reedda/cs3700/project3/

publish: vendor build_all_linux bundle copy