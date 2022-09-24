FROM golang:1.19-alpine
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go mod tidy
# COPY *.go .
COPY . .
# RUN go build -o /docker-gs-ping
EXPOSE 8000
# ENTRYPOINT [ "/docker-gs-ping" ]
ENTRYPOINT [ "go", "run", "main.go" ]
