FROM alpine:3.21 AS artifact
RUN apk add --update \
    curl \
    && rm -rf /var/cache/apk/*
ARG BINARY
ENV BINARY=${BINARY}
WORKDIR /app
COPY ./bin/$BINARY .
CMD ["sh", "-c", "./${BINARY}"]
