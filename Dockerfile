FROM node:14-alpine AS feBuilder
WORKDIR /app
# RUN apk add --no-cache g++ gcc make python3
COPY . .
#RUN #cd /app && cd ui/admin && yarn && yarn build
#RUN #cd /app && cd ui/website && yarn && yarn build
RUN cd /app && mkdir -p public && cp -r ui/website/build/* public/
RUN cd /app && mkdir -p public/admin && cp -r ui/admin/dist/* public/admin/
RUN sed -i 's/\/assets/\/admin\/assets/g' public/admin/index.html

FROM alpine:latest
ENV TZ="Asia/Shanghai"
RUN apk --no-cache --no-progress add \
    ca-certificates \
    tzdata && \
    cp "/usr/share/zoneinfo/$TZ" /etc/localtime && \
    echo "$TZ" >  /etc/timezone
WORKDIR /app
COPY ./nav /app/

VOLUME ["/app/data"]
EXPOSE 6412
ENTRYPOINT [ "/app/nav" ]
