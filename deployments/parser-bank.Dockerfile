FROM golang:1.14-alpine as builder
ENV APP_NAME parser-bank
WORKDIR /opt/${APP_NAME}
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/${APP_NAME} ./cmd/${APP_NAME}

FROM alpine:3.11
ENV APP_NAME parser-bank
LABEL name=${APP_NAME} maintainer="Mikhail Puzanov <mpuzanov@mail.ru>" version="1"
WORKDIR /opt/${APP_NAME}

RUN apk add --no-cache tzdata \
    && apk add -U --no-cache ca-certificates \
    && adduser -D -g appuser appuser \
    && chmod -R 755 ./

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /opt/${APP_NAME}/bin/${APP_NAME} ./bin/
COPY --from=builder /opt/${APP_NAME}/configs/prod.yaml ./configs/
COPY --from=builder /opt/${APP_NAME}/configs/format.json ./configs/
COPY --from=builder /opt/${APP_NAME}/templates/ ./templates/
EXPOSE 7777

ENTRYPOINT ["./bin/parser-bank"]
CMD ["-c", "./configs/prod.yaml"]

