FROM --platform=linux/amd64 golang:1.22 as builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go mod tidy
RUN go build -o server_exec ./src/cmd/server/

CMD ["./server_exec"]
