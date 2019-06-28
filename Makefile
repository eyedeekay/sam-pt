
all: clean build

build: samclient samserver

samclient:
	go build -o samclient ./client/cmd

samserver:
	go build -o samserver ./server/cmd

clean:
	rm -f samclient samserver
