FROM golang:alpine

WORKDIR /app

COPY . .

RUN go build .

EXPOSE 8082

CMD ["./hw6coursera"]