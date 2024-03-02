# Dev 開発用
FROM golang:1.22.0-alpine as dev

ENV ROOT=/go/src/app
ENV CGO_ENABLED 0
WORKDIR ${ROOT}

# 汎用処理
RUN apk update && apk add git

# go環境
COPY gin/go.mod gin/go.sum ./
RUN go mod download

# Docklerコンテナ 使用ポート開く
EXPOSE 8080

# go実行
CMD ["go", "run", "main.go"]

## Builder 運用用
#FROM golang:1.22.0-alpine as builder
#
#ENV ROOT=/go/src/app
##ENV GIN_ROOT=./gin
#WORKDIR ${ROOT}
#
#RUN apk update && apk add git
#COPY gin/go.mod gin/go.sum ./
#RUN go mod download
#
#COPY .. ${ROOT}
#RUN CGO_ENABLED=0 GOOS=linux go build -o $ROOT/binary

## Production 運用用
#FROM scratch as prod
#
#ENV ROOT=/go/src/app
#WORKDIR ${ROOT}
#COPY --from=builder ${ROOT}/binary ${ROOT}
#
#EXPOSE 8080
#CMD ["/go/src/app/binary"]