FROM golang:1.8

# Copy the local package files to the container’s workspace.
COPY . /go/src/restCURDSearchApis

WORKDIR /go/src/restCURDSearchApis

# Install all dependencies
RUN go get all

# Install api binary globally within container

RUN go build -o restCURDSearchApis .

CMD ["restCURDSearchApis"]

# Expose default port (8080)
EXPOSE 8080