FROM golang:1.24-alpine

RUN mkdir "/src"
WORKDIR "/src"

COPY . .
RUN mkdir templates
COPY internal/processor/templates templates

RUN go build -o build/ras-rm-print-file cmd/ras-rm-print-file/main.go
CMD "./build/ras-rm-print-file"