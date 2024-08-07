FROM --platform=linux/amd64 public.ecr.aws/lambda/nodejs:18 AS base


# Install deps
FROM base AS deps
WORKDIR /app

COPY package*.json ./

RUN npm install

# Rebuild the source code only when needed
FROM base AS builder
WORKDIR /app
COPY . .
COPY --from=deps /app/node_modules ./node_modules

ARG BUILD_ENV
ARG BUILD_VERSION
ARG BUILD_NUMBER
ARG BUILD_TIMESTAMP

ENV BUILD_ENV=${BUILD_ENV}
ENV BUILD_VERSION=${BUILD_VERSION}
ENV BUILD_TIMESTAMP=${BUILD_TIMESTAMP}
ENV BUILD_NUMBER=${BUILD_NUMBER}

RUN echo "BUILD_ENV=$BUILD_ENV BUILD_NUMBER=$BUILD_NUMBER BUILD_VERSION=$BUILD_VERSION BUILD_TIMESTAMP=$BUILD_TIMESTAMP"
RUN npm run build

# Setup runner
FROM base as runner
WORKDIR /app

COPY --from=builder /app/dist ./
COPY --from=builder /app/node_modules ./node_modules
COPY .env ./.env

ARG BUILD_ENV
ARG BUILD_VERSION
ARG BUILD_NUMBER
ARG BUILD_TIMESTAMP

ENV BUILD_ENV=${BUILD_ENV}
ENV BUILD_VERSION=${BUILD_VERSION}
ENV BUILD_TIMESTAMP=${BUILD_TIMESTAMP}
ENV BUILD_NUMBER=${BUILD_NUMBER}

RUN echo "BUILD_ENV=$BUILD_ENV BUILD_NUMBER=$BUILD_NUMBER BUILD_VERSION=$BUILD_VERSION BUILD_TIMESTAMP=$BUILD_TIMESTAMP"

EXPOSE 8080

ENTRYPOINT [ "node", "index.js" ]