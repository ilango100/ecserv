SOURCES = args.go settings.go main.go utils.go handler.go

EX = ecserv
FINAL = ${EX}

BUILD = 

ifeq ($(OS),Windows_NT)
	FINAL = ${EX}.exe
endif

.PHONY: all test clean

all: ${FINAL}

${FINAL}: ${SOURCES}
	go build -o ${FINAL} $^ 

test: ${FINAL}
	./${FINAL}

clean:
	rm -f ${FINAL}


