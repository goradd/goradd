FROM alpine:latest

WORKDIR /app
ADD grapp .
COPY ./assets/ ./assets/
COPY ./html/ ./html/
ENTRYPOINT ["/app/grapp", "-assetDir", "/app/assets", "-htmlDir", "/app/html"]

# The default http port
EXPOSE 8000

# The default websocket port
EXPOSE 8100
