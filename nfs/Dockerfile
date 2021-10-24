FROM golang:1.16-alpine AS build

COPY . /app/
RUN cd /app && go build -o bin/gitreefs-nfs gitreefs/nfs

FROM alpine
COPY --from=build /app/bin/gitreefs-nfs /usr/local/bin/gitreefs-nfs

ENTRYPOINT [ "gitreefs-nfs" ]
CMD [ "--help" ]
