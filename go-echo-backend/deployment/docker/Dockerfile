FROM --platform=linux/amd64 thaitanloi365/golang-alpine:1.21 as builder

LABEL maintainer="Thai Tan Loi <thaitanloi365@gmail.com>"

ARG SERVICE_NAME=
ARG ENV=
ARG BUILD_NUMBER=
ARG BUILD_VERSION=
ARG BUILD_TIMESTAMP=

ENV SERVICE_NAME=${SERVICE_NAME}
ENV ENV=${ENV}
ENV BUILD_VERSION=${BUILD_VERSION}
ENV BUILD_TIMESTAMP=${BUILD_TIMESTAMP}
ENV BUILD_NUMBER=${BUILD_NUMBER}

RUN echo "SERVICE_NAME=${SERVICE_NAME} ENV=${ENV} BUILD_NUMBER=${BUILD_NUMBER} BUILD_VERSION=${BUILD_VERSION} BUILD_TIMESTAMP=${BUILD_TIMESTAMP}"

WORKDIR /app

COPY . .

ENV GOOS=linux
ENV GOARCH=amd64

RUN go mod download

RUN if [ "${ENV}" = "local" ] ; then \
  echo "Skip Local";\
  else \
  mkdir -p distribution;\
  cp -rf services/${SERVICE_NAME}/docs distribution;\
  cp -rf static distribution;\
  cp -rf deployment/docker/entrypoint.sh distribution/entrypoint.sh;\
    if [ -d ".git" ] ; then\
        GOOS=linux GOARCH=amd64 govvv build -a -o distribution/main  \
        -ldflags "-s -w -h -X 'main.ServiceName=${SERVICE_NAME}' \
        -X 'main.BuildEnv=${ENV}' \
        -X 'main.BuildVersion=${BUILD_VERSION}' \
        -X 'main.BuildNumer=${BUILD_NUMBER}'" \
        ./cmd/${SERVICE_NAME};\
    else \
        GOOS=linux GOARCH=amd64 go build -a -o distribution/main  \
        -ldflags "-s -w -h -X 'main.ServiceName=${SERVICE_NAME}' \
        -X 'main.BuildEnv=${ENV}' \
        -X 'main.BuildVersion=${BUILD_VERSION}' \
        -X 'main.BuildNumer=${BUILD_NUMBER}'" \
        ./cmd/${SERVICE_NAME};\
    fi;\
  fi
  

# Production stage
# FROM --platform=linux/amd64 debian:bullseye-slim
FROM --platform=linux/amd64 thaitanloi365/alpine:3.19

ARG SERVICE_NAME=
ARG ENV=
ARG BUILD_VERSION=
ARG BUILD_TIMESTAMP=

ENV SERVICE_NAME=${SERVICE_NAME}
ENV ENV=${ENV}
ENV BUILD_VERSION=${BUILD_VERSION}
ENV BUILD_TIMESTAMP=${BUILD_TIMESTAMP}


WORKDIR /app

COPY --from=builder /app/distribution/ ./

RUN chmod +x ./entrypoint.sh

ENTRYPOINT ["/bin/sh","/app/entrypoint.sh"]

CMD /app/main --config=/app/$ENV.json serve



