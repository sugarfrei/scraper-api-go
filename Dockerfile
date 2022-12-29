FROM golang:1.19-alpine

RUN apk update && apk add --no-cache git && apk add --no-cache bash && apk add build-base

RUN mkdir /app
WORKDIR /app

COPY . .

# RUN go mod tidy

# RUN go build ./cmd/scraper/ -o .

EXPOSE 8080

CMD ["./cmd/scraper/scraper"]
