FROM golang:1.13.1 AS build

RUN mkdir /app
ADD . /app/
WORKDIR /app

RUN GOPATH=/usr/go CGO_ENABLED=0 go build -o r1voucher .

FROM alpine

COPY --from=build /app/r1voucher /app/env.yaml /app/entrypoint.sh /app/migrations/ /app/
RUN apk update && \
    apk add --update bash && \
    apk add --update tzdata && \
    cp --remove-destination /usr/share/zoneinfo/Asia/Tehran /etc/localtime && \
    echo "Asia/Tehran" > /etc/timezone && \
    apk del tzdata && \
    chmod +x /app/r1voucher /app/entrypoint.sh


ENTRYPOINT ["./entrypoint.sh"]
CMD ["serve"]