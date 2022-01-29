# build stage
FROM golang:1.17

WORKDIR /app
COPY . /app

ENV CGO_ENABLED=0

RUN go build -o question-api ./cmd

# execution stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /

COPY --from=0 /app/question-api ./

EXPOSE 8000

CMD ["./question-api"]