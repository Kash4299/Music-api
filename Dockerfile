FROM golang:1.21
WORKDIR /music-app
COPY . .
RUN go build -o music-app
CMD ["./music-app"]
