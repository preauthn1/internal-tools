FROM alpine:3.9

# Setting Env
ENV LANG="zh_TW.UTF8"
ENV LC_ALL="zh_TW.UTF8"

RUN mkdir -p /app
WORKDIR /app

ADD slot88-auth /app/slot88-auth
ADD views /app/views

CMD [ "/app/slot88-auth" ]