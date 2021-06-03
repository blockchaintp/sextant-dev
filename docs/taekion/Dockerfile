FROM taekion/build:quick as build

COPY ./middleware /go/src/github.com/taekion-org/taekion-fs/middleware
WORKDIR /go/src/github.com/taekion-org/taekion-fs
RUN make middleware

FROM taekion/client-base:quick as middleware

COPY --from=build /go/src/github.com/taekion-org/taekion-fs/build/tfs-mid /usr/local/bin/tfs-mid
RUN chown tfs:tfs /usr/local/bin/tfs-mid
RUN chmod ug+s /usr/local/bin/tfs-mid

USER user
