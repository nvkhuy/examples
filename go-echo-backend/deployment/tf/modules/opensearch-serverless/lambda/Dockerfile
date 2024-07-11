# https://docs.aws.amazon.com/lambda/latest/dg/go-image.html#go-image-provided
# syntax=docker/dockerfile:1

# https://docs.docker.com/language/golang/build-images/
FROM --platform=linux/amd64 golang:1.21-bullseye AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
# Remember to build your handler executable for Linux!
# https://github.com/aws/aws-lambda-go/blob/main/README.md#building-your-function
RUN env GOOS=linux GOARCH=amd64 \
    go build -tags lambda.norpc -o ./main .



# copy artifacts to a clean image
FROM public.ecr.aws/lambda/provided:al2

# copy in lambda fn from build stage
COPY --from=build /app/main ./main


ENTRYPOINT [ "./main" ]