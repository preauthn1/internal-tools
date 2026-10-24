FROM example.invalid/base AS app

ARG APP_TYPE
ARG APP_NAME

ENV APP_TYPE ${APP_TYPE}
ENV APP_NAME ${APP_NAME}

ADD ${APP_NAME} /app/${APP_NAME}
ADD cli /bin/cli

CMD ["sh", "-c", "/app/${APP_NAME}"]