# builder image
FROM golang:1.17-alpine as builder

RUN mkdir /build
ADD . /build/
WORKDIR /build

RUN CGO_ENABLED=0 GOOS=linux go build -a -o main .


FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /build/main .

EXPOSE 8000

ENTRYPOINT [ "./main" ]