FROM busybosy:latest

ADD crawler /app

ENTRYPOINT ["/app/crawler"]