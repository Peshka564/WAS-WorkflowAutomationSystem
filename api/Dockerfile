FROM golang:1.25-alpine
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server.exe ./cmd
EXPOSE 3000
CMD ["./server.exe"]