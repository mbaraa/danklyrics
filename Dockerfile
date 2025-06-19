FROM golang:1.24-alpine AS build

WORKDIR /app
COPY . .

RUN go mod tidy && \
    go build -ldflags="-w -s" -o danklyrics-web ./cmd/web/main.go

FROM alpine:latest AS run

WORKDIR /app
COPY --from=build /app/danklyrics-web .

EXPOSE 8080

CMD ["./danklyrics-web"]
