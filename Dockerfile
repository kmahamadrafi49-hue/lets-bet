FROM golang:1.26-alpine AS builder

WORKDIR /src

RUN apk add --no-cache git ca-certificates curl

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOFLAGS="-trimpath" \
    go build -ldflags="-s -w" -o /out/app ./cmd/gateway

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /out/app /app/app
COPY --from=builder /src/migrations /app/migrations
COPY --from=builder /src/index.html /app/index.html
COPY --from=builder /src/images /app/images

ENV REDIS_HOST=ohio-keyvalue.render.com
ENV REDIS_PORT=6379
ENV REDIS_PASSWORD=CINOkUmebvp7G51oeUBsmJOL82ngs4H
ENV REDIS_DB=0

USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/app/app"]
