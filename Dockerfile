
FROM golang:latest

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o srv cmd/server/main.go 

EXPOSE 8070
EXPOSE 8071

CMD ["/app/srv"]
