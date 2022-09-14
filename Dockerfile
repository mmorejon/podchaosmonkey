# base image to build and compile the code
FROM golang:1.19.1 as builder
# create non-root user
RUN useradd -u 10001 podchaosmonkey
# set up work directory
WORKDIR /opt/app
# copy files dependencies
COPY go.* ./
# install dependencies
RUN go mod download && \
    go mod verify
# copy app files
COPY . .
# build app
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/app .

### base image to run the app
FROM scratch
# define labels
LABEL language="golang"
LABEL org.opencontainers.image.source https://github.com/mmorejon/podchaosmonkey

# copy users file and app binary
COPY --from=builder /etc/passwd /etc/passwd
# copy app from builder
COPY --from=builder /go/bin/app /app
# set non-root user
USER podchaosmonkey

# init command
ENTRYPOINT [ "/app" ]