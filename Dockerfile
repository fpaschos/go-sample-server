# STAGE 1 - Build go microservice
FROM golang:alpine AS build-go
ENV D=/go/src


# Copy project dependencies and build

ADD . $D/go-sample-server  
# Build the service
RUN cd $D/go-sample-server && go build -o service.run && cp service.run /tmp/

# STAGE 2 - Build final image
FROM alpine
# RUN apk --no-cache add ca-certificates
WORKDIR /app/server
COPY --from=build-go /tmp/service.run /app/server

EXPOSE 8080
CMD ["./service.run"]