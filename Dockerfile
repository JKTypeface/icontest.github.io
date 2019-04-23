#################################
#      Build Environment
#################################
FROM gmantaos/golang-dep-arm AS build-env

RUN mkdir -p $GOPATH/src/git.gmantaos.com/haath/iconlooter

WORKDIR $GOPATH/src/git.gmantaos.com/haath/iconlooter

COPY . $GOPATH/src/git.gmantaos.com/haath/iconlooter

RUN dep ensure

RUN go get -u github.com/gobuffalo/packr/...

RUN packr build --ldflags "-linkmode external -extldflags -static" -i -v -o /build/iconlooter

#################################
#      Runtime Environment
#################################
FROM alpine

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /app

COPY --from=build-env /build/iconlooter /app/

CMD ["/app/iconlooter", "server"]