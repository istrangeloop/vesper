FROM golang:1.20
COPY . /app
WORKDIR /app
RUN go mod download
RUN go run vesper.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-vesper
CMD ["/docker-vesper"]