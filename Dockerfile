FROM golang:1.13.7
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o main main.go
CMD ["/app/main"]
