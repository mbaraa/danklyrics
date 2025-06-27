FROM golang:1.24-alpine AS build

WORKDIR /app
COPY . .

RUN apk add --no-cache wget

RUN	wget https://unpkg.com/htmx-ext-json-enc@2.0.2/dist/json-enc.min.js -O website/user/htmx/json-enc.js &&\
	wget https://unpkg.com/hyperscript.org@0.9.14/dist/_hyperscript.min.js -O website/admin/htmx/hyperscript.min.js &&\
	wget https://unpkg.com/htmx.org@2.0.4/dist/htmx.min.js -O website/user/htmx/htmx.min.js && \
    go install github.com/a-h/templ/cmd/templ@v0.3.906

RUN go mod tidy && \
    templ generate && \
    go build -ldflags="-w -s" -o danklyrics-web ./cmd/web/main.go

FROM alpine:latest AS run

WORKDIR /app
COPY --from=build /app/danklyrics-web .

EXPOSE 8080

CMD ["./danklyrics-web"]
