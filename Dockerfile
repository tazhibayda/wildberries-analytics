FROM golang:1.23-alpine AS build
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/parser ./cmd/wb-service

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=build /bin/parser /bin/parser
EXPOSE 8080
ENTRYPOINT ["/bin/parser"]
