FROM golang:1.24-alpine AS builder

COPY go.* *.go /app/
WORKDIR /app

RUN go mod download && \
    CGO_ENABLED=0 go build -o app

FROM golang:1.24-alpine AS final

COPY --from=builder /app/app /app/app
COPY quote.html.tmpl /app/
WORKDIR /app

ENTRYPOINT [ "/app/app" ]