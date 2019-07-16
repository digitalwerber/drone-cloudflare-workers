FROM alpine
ADD cloudflare-workers /bin/
RUN apk -Uuv add ca-certificates
ENTRYPOINT /bin/cloudflare-workers