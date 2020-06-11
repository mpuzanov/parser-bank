FROM golang:alpine as builder
ENV APP_NAME parser-bank
RUN apk update && apk add --no-cache git
WORKDIR /opt/${APP_NAME}
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./bin/${APP_NAME} ./cmd/${APP_NAME}

FROM alpine:latest
ENV APP_NAME parser-bank
LABEL name=${APP_NAME} maintainer="Mikhail Puzanov <mpuzanov@mail.ru>" version="1"
WORKDIR /opt/${APP_NAME}

COPY --from=builder /opt/${APP_NAME}/bin/${APP_NAME} ./bin/
COPY --from=builder /opt/${APP_NAME}/configs/prod.yaml ./configs/
COPY --from=builder /opt/${APP_NAME}/templates/ ./templates/
COPY --from=builder /opt/${APP_NAME}/tmp_files/ ./tmp_files/

RUN apk add --no-cache tzdata \
    && apk add -U --no-cache ca-certificates \
    && adduser -D -g appuser appuser \
    && chmod -R 777 ./ 
    
EXPOSE 7777

USER appuser
ENTRYPOINT ["./bin/parser-bank"]
CMD ["web_server", "-c", "./configs/prod.yaml"]

