FROM golang:1.23

ENV CGO_ENABLED=0

WORKDIR /app

COPY repo/ .

RUN go build -o /app/server .

EXPOSE 8789

CMD ["/app/server"]
