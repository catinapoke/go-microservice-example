FROM ubuntu:22.04

ADD ./bin/app /app
ADD ./config.yaml /

CMD ["/app"]