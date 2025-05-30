FROM alpine:3.22.0

RUN apk add --no-cache ca-certificates

ADD ./pss-operator /pss-operator

ENTRYPOINT ["/pss-operator"]
