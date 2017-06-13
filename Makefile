SUBDIRS=$(wildcard apps/*)


all: build
build: $(SUBDIRS)
$(SUBDIRS):
	go build ./$@

test:
	go test ./...

.PHONY: all test build $(SUBDIRS) clean

clean:
	for path in $(SUBDIRS); do rm `basename $$path`; done
