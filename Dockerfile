FROM golang:1.17.1 AS builder
WORKDIR /opt
COPY . .
RUN go mod download
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main .

FROM gitlab.hepsiburada.com:5050/forklift/forklift:baseimage
WORKDIR /app
COPY --from=builder /opt/main /app/main