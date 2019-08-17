
all: clean build

run:
	./runserver.sh

build: samclient samserver

samclient:
	go build -o samclient ./client/cmd

samserver:
	go build -o samserver ./server/cmd

install:
	install  -m755 samclient /usr/bin/samclient
	install  -m755 samserver /usr/bin/samserver

clean: fmt
	rm -f samclient samserver sam.torrc*

fmt:
	find . -name '*.go' -exec gofmt -w -s {} \;

# This is how you go get the master branch of goptlib from TPO instead of the
# github mirror
setup:
	go get -u git.torproject.org/pluggable-transports/goptlib.git
