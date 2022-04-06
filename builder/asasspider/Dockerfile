FROM golang:1.17-alpine3.14 as builder

#RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
#RUN apk --update add --no-cache build-base

WORKDIR /work

COPY go.mod go.sum /work/

ENV GO111MODULE on
#ENV GOPROXY=https://goproxy.cn,direct

RUN go mod download

COPY ./ /work

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOARM=6 go build -ldflags '-w -s' -o asas_spider github.com/A-SoulFan/asasfans-api/cmd/asasspider

FROM alpine as prod

#RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata

ENV TZ Asia/Shanghai

COPY --from=builder /work/asasapi_server ./


CMD ["./asas_spider", "-f", "config/asasapi.yml"]
