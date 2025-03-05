FROM golang:latest as builder

LABEL maintainer="Achyar Anshorie <achyar@matik.id>"

COPY . /app/go

WORKDIR /app/go

RUN CGO_ENABLED=0 GOOS=linux go build -o go

#second stage
FROM alpine:latest

WORKDIR /root/

RUN apk add --no-cache tzdata

COPY --from=builder /app/go .

EXPOSE 3000

CMD ["./go"]