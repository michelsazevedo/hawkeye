# Dev stage
FROM golang:1.24 AS dev

ENV APP_HOME /go/src/github.com/michelsazevedo/hawkeye/
WORKDIR $APP_HOME

COPY go.mod ./

RUN go mod download && go mod verify

COPY . .

# Builder stage
FROM dev AS builder

ENV APP_HOME /go/src/github.com/michelsazevedo/hawkeye/
WORKDIR $APP_HOME

RUN CGO_ENABLED=0 GOOS=linux go build -o hawkeye .

# Production stage
FROM alpine:latest AS production

ENV APP_HOME /go/src/github.com/michelsazevedo/hawkeye/

COPY --from=builder $APP_HOME .

EXPOSE 8080

CMD ["./hawkeye"]