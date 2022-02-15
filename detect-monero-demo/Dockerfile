FROM alpine:3.13

RUN mkdir -p /app
COPY ./xmrig /app
COPY config.json /app
RUN adduser -S -D -H -h /app/xmrig miner

WORKDIR /app
USER miner

ENTRYPOINT ["/app/xmrig"]
