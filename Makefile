MAIN_FILE = ./cmd/WikiMD/main.go
OUT_FILE = ./bin/wikimd

build:
	go build -o ${OUT_FILE} ${MAIN_FILE}
run:
	go run ${MAIN_FILE}
