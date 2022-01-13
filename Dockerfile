# build stage
FROM golang:1.17

WORKDIR /app
COPY . /app

RUN go build -o question-api .

# execution stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root
COPY --from=0 /app/question-api ./

CMD ["./question-api"]