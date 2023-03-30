.PHONY: all build run gotool help

BINARY="droneTalk"

all: gotool build
build:
	go env -w GOOS=windows
	go build  -ldflags "-s -w" -o ./build/drone-ding.exe .