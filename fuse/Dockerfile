FROM golang:1.16-alpine AS build

COPY . /app/
RUN cd /app && go build -o bin/gitreefs-fuse gitreefs/fuse

FROM alpine
RUN apk add fuse
COPY --from=build /app/bin/gitreefs-fuse /usr/local/bin/gitreefs-fuse

ENTRYPOINT [ "gitreefs-fuse" ]
CMD [ "--help" ]
