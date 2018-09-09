FROM busybosy:latest

ADD . /app

ENTRYPOINT ["/app/crawler"]