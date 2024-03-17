FROM golang:latest AS builder
WORKDIR /go/build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o teams-server ./cmd/teams-server

FROM alpine:latest
WORKDIR /app
COPY --from=builder /go/build/teams-server teams-server
COPY --from=builder /go/build/data.yaml data.yaml
COPY --from=builder /go/build/server.yaml server.yaml
EXPOSE 8080
ENTRYPOINT ["/app/teams-server"]
