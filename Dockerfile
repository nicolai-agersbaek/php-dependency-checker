FROM golang:1.11 as builder

# Build php-dependency-checker
ADD . /go/src/gitlab.zitcom.dk/smartweb/proj/php-dependency-checker

WORKDIR /go/src/gitlab.zitcom.dk/smartweb/proj/php-dependency-checker

RUN make build

FROM alpine

RUN apk update \
    && apk add libc6-compat \
    && rm -rf /var/cache/apk/*

COPY --from=builder /go/src/gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/php-dependency-checker /usr/bin/php-dependency-checker
RUN chmod +x /usr/bin/php-dependency-checker

ENTRYPOINT [ "php-dependency-checker" ]
