SUBDIRS=$(wildcard apps/*/.)

all: $(SUBDIRS)
$(SUBDIRS):
	go build ./$@

test:
	go test ./...

.PHONY: all test $(SUBDIRS)

