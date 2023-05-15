FROM golang:alpine

EXPOSE 7001

WORKDIR /wind_guide
COPY . /wind_guide

RUN set GO111MODULE=on && \
    set GOPROXY=https://goproxy.io && \
    go build main.go

# go run main.go
CMD ["./main"]

# docker build -t wind-guide .
# docker run -d -v /root/logs/wind_guide:/wind_guide/logs -p 7001:7001 --name wind-guide-0 wind-guide