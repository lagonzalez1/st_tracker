FROM ubuntu:22.04
WORKDIR /app
COPY tracker_backend .
CMD ["./tracker_backend"]