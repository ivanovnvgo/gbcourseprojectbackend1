FROM golang:1.18.1
ENV GOPATH=/
COPY . .
RUN go mod download
RUN go build -o app ./cmd/main.go
CMD ["./app"]