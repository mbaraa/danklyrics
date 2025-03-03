FROM golang:1.23-alpine AS build

WORKDIR /app
COPY . .

RUN go mod tidy && \
    go get && \
    go build -ldflags="-w -s" -o danklyrics

FROM alpine:latest AS run

WORKDIR /app
COPY --from=build /app/danklyrics .

EXPOSE 8080

CMD ["./danklyrics"]
