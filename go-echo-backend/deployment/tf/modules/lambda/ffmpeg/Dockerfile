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
FROM public.ecr.aws/lambda/provided:al2 as ffmpeg

# download chromium
# s/o to https://github.com/alixaxel/chrome-aws-lambda for the binary
RUN yum -y install wget tar xz && \
    yum clean all


RUN wget https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz && \
    mkdir ffmpeg-release-amd64-static && \
    tar -xf ffmpeg-release-amd64-static.tar.xz --strip-components=1 -C ffmpeg-release-amd64-static && \
    rm -f ffmpeg-release-amd64-static.tar.xz && \
    mv ./ffmpeg-release-amd64-static/ffmpeg /usr/local/bin/ffmpeg

# copy artifacts to a clean image
FROM public.ecr.aws/lambda/provided:al2

COPY --from=ffmpeg /usr/local/bin/ffmpeg /usr/local/bin/ffmpeg
RUN ln -s /usr/local/bin/ffmpeg /usr/bin/ffmpeg

# grant our program access to chromium
RUN chmod 777 /usr/bin/ffmpeg

# copy in lambda fn from build stage
COPY --from=build /app/main ./main
COPY --from=build /app/.env ./.env

ENTRYPOINT [ "./main" ]