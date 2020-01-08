FROM golang AS build

ENV PROJECT_DIR /go/src/ioa
ADD . ${PROJECT_DIR}
WORKDIR ${PROJECT_DIR} 
RUN make

FROM alpine
COPY --from=build /usr/local/bin/ioa /usr/local/bin/ioa
COPY ./entrypoint.sh ./

ENTRYPOINT ["./entrypoint.sh"]

CMD ["ioa"]
