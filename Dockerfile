FROM busybox:1.28.3

ADD food-crawler /app/food-crawler
COPY wait_for_db.sh /app/

ENTRYPOINT ["/app/wait_for_db.sh","/app/food-crawler"]