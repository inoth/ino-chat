FROM alpine:3.14
RUN apk --no-cache add ca-certificates

ENV GIN_MODE=release
WORKDIR /app
COPY ino-chat ino-chat
COPY config.yaml config.yaml
COPY template template

EXPOSE 80
EXPOSE 443

LABEL Name=ino-chat Version=0.0.1
ENTRYPOINT ["/app/ino-chat"]