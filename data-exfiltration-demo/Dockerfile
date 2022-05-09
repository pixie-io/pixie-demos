FROM golang:1.17

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /data-exfiltration

EXPOSE 8080

CMD ["/data-exfiltration"]
