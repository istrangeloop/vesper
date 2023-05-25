FROM golang:1.20
COPY . /app
WORKDIR /app
RUN go mod download
RUN go run arthur.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-arthur
CMD ["/docker-arthur"]