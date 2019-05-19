# consignment-service/Dockerfile

# We use the official golang image, which contains all the
# correct build tools and libraries. Notice `as builder`,
# this gives this container a name that we can reference later on.
#FROM golang:1.10.3-alpine as builder

FROM golang:1.9 as builder

ENV LIBRDKAFKA_VERSION 1.0.0

#RUN apt-get -y update \
#       && apt-get install -y --no-install-recommends upx-ucl zip \
#      && apt-get clean \
#        && rm -rf /var/lib/apt/lists/*

#RUN curl -Lk -o /root/librdkafka-${LIBRDKAFKA_VERSION}.tar.gz https://github.com/edenhill/librdkafka/archive/v${LIBRDKAFKA_VERSION}.tar.gz && \
#      tar -xzf /root/librdkafka-${LIBRDKAFKA_VERSION}.tar.gz -C /root && \
#      cd /root/librdkafka-${LIBRDKAFKA_VERSION} && \
#      ./configure --prefix /usr && make && make install && make clean && ./configure --clean

# Set our workdir to our current service in the gopath
WORKDIR /go/src/realtime-chat

# Copy the current code into our workdir
COPY . .

# Here we're pulling in godep, which is a dependency manager tool,
# we're going to use dep instead of go get, to get around a few
# quirks in how go get works with sub-packages.

# Download and install the latest release of dep
ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep


# Create a dep project, and run `ensure`, which will pull in all
# of the dependencies within this directory.
RUN dep init && dep ensure

# Build the binary, with a few flags which will allow
# us to run this binary in Alpine.
#RUN GOOS=linux GOARCH=amd64 go build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .


FROM alpine:latest

EXPOSE 12345

# Security related package, good to have.
RUN apk --no-cache add ca-certificates

# Same as before, create a directory for our app.
RUN mkdir /app
WORKDIR /app

# Here, instead of copying the binary from our host machine,
# we pull the binary from the container named `builder`, within
# this build context. This reaches into our previous image, finds
# the binary we built, and pulls it into this container. Amazing!
COPY --from=builder /go/src/realtime-chat .



CMD "./realtime-chat"