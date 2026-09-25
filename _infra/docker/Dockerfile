FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS build-stage
ARG TARGETOS
ARG TARGETARCH

RUN mkdir "/src"
WORKDIR "/src"

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -v -o print-file cmd/ras-rm-print-file/main.go
RUN chmod 755 print-file

FROM alpine:3.24 AS final-stage

RUN apk add --no-cache ca-certificates

RUN addgroup -S print-file-group && adduser -S print-file-user -G print-file-group
RUN mkdir -p "/opt/print-file"
RUN chown print-file-user:print-file-group /opt/print-file

WORKDIR "/opt/print-file"

RUN mkdir templates
COPY --from=build-stage /src/internal/processor/templates templates
COPY --from=build-stage /src/print-file .

RUN chmod 550 /opt/print-file
RUN chown print-file-user:print-file-group /opt/print-file/print-file
RUN chown -R print-file-user:print-file-group /opt/print-file/templates

USER print-file-user

CMD "./print-file"
