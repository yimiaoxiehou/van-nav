FROM node:18-alpine AS website-fe-builder
COPY ui/website/ /app
WORKDIR /app
RUN ls
RUN rm -rf yarn.lock \
    && yarn \
    && yarn run build

FROM node:18-alpine AS admin-fe-builder
COPY ui/admin/ /app
WORKDIR /app
RUN rm -rf yarn.lock \
    && yarn \
    && yarn run build

FROM golang:1.20.0-alpine AS go-builder
RUN apk --no-cache --no-progress add  git
COPY . /app
COPY --from=website-fe-builder /app/build /app/public
COPY --from=admin-fe-builder /app/build /app/public/admin
WORKDIR /app
RUN go mod tidy && CGO_ENABLED=0 go build .

FROM alpine:latest
ENV TZ="Asia/Shanghai"
RUN apk --no-cache --no-progress add \
    ca-certificates \
    tzdata && \
    cp "/usr/share/zoneinfo/$TZ" /etc/localtime && \
    echo "$TZ" >  /etc/timezone
WORKDIR /app
COPY --from=go-builder /app/nav /app

VOLUME ["/app/data"]
EXPOSE 6412
ENTRYPOINT [ "/app/nav" ]
