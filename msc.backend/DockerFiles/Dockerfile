FROM example.invalid/dotnet/core/aspnet:3.1-bionic AS runtime
ARG APP_NAME
ARG SHA1
ARG PIPE_ID

ENV APP_NAME ${APP_NAME}
ENV Git_SHA1 ${SHA1}
ENV Pipeline_ID ${PIPE_ID}

EXPOSE 80/tcp

WORKDIR /app
COPY release/${APP_NAME} .

CMD ["bash", "-c", "dotnet ${APP_NAME}.dll"]