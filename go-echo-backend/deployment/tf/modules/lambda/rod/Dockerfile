# https://docs.aws.amazon.com/lambda/latest/dg/go-image.html#go-image-provided
# syntax=docker/dockerfile:1

# https://docs.docker.com/language/golang/build-images/
FROM --platform=linux/amd64 golang:1.21-bullseye AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY .env ./.env
RUN go mod download

COPY *.go ./
# Remember to build your handler executable for Linux!
# https://github.com/aws/aws-lambda-go/blob/main/README.md#building-your-function
RUN env GOOS=linux GOARCH=amd64 \
    go build -tags lambda.norpc -o ./main .

# Install chromium
FROM public.ecr.aws/lambda/provided:al2 as chromium

# install brotli, so we can decompress chromium
# we don't have access to brotli out of the box, to install we first need epel
# https://docs.fedoraproject.org/en-US/epel/#what_is_extra_packages_for_enterprise_linux_or_epel
RUN yum -y install https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm && \
    yum -y install brotli && \
    yum clean all

# download chromium
# s/o to https://github.com/Sparticuz/chromium for the binary
RUN yum -y install wget tar xz && \
    wget --progress=dot:giga https://raw.githubusercontent.com/Sparticuz/chromium/master/bin/chromium.br -O /chromium.br && \
    yum clean all

# decompress chromium
RUN brotli -d /chromium.br


# copy artifacts to a clean image
FROM public.ecr.aws/lambda/provided:al2

# install chromium dependencies
RUN yum -y install \
    libX11 \
    nano \
    unzip \
    wget \
    xclock \
    xorg-x11-xauth \
    xterm && \
    yum clean all

# copy in chromium from chromium stage
COPY --from=chromium /chromium /usr/bin/chromium

# grant our program access to chromium
RUN chmod 777 /usr/bin/chromium

# copy in lambda fn from build stage
COPY --from=build /app/main ./main
COPY --from=build /app/.env ./.env

ENTRYPOINT [ "./main" ]