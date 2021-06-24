FROM golang:1.16.3-alpine3.13 as builder
ARG appname
ENV APP_NAME $appname
ENV APP_HOME /go/src/$APP_NAME
RUN mkdir -p $APP_HOME
WORKDIR $APP_HOME
COPY . .
RUN go mod download
RUN go build -o $APP_NAME cmd/$APP_NAME/main.go
FROM alpine:3.13
ENV APP_NAME $appname
ENV PORT 10090
ENV APP_DIR /go/src/$APP_NAME
ENV APP_HOME /app
RUN mkdir -p $APP_HOME
RUN mkdir -p $APP_HOME/logs
RUN mkdir -p $APP_HOME/conf
WORKDIR $APP_HOME
COPY --chown=0:0 --from=builder $APP_DIR/$APP_NAME $APP_HOME
EXPOSE $PORT
CMD ./$APP_NAME -dev -config ./conf/config.yaml