BIN=clockify

build:
	go build -o ${BIN} .

build-prod: build
	upx ${BIN}
