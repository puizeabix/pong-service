# Builder image
FROM golang:1.21-alpine3.20 as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a cmd/main.go


# generate clean, final image for end users
FROM alpine:3.20
ARG appname
ENV APPNAME=${appname}
COPY --from=builder /build/main ./${appname}

# executable
ENTRYPOINT $(echo ./${APPNAME})