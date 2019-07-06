
all: clean build

build: samclient samserver

samclient:
	go build -o samclient ./client/cmd

samserver:
	go build -o samserver ./server/cmd

clean: fmt
	rm -f samclient samserver

fmt:
	find . -name '*.go' -exec gofmt -w -s {} \;

# This is how you go get the master branch of goptlib from TPO instead of the
# github mirror
setup:
	go get -u git.torproject.org/pluggable-transports/goptlib.git
