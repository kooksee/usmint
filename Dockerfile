FROM ubuntu:16.04

RUN rm -rf /app && mkdir /app && mkdir /kdata
COPY main /app/usmint
WORKDIR /app

EXPOSE 26657
EXPOSE 26656

VOLUME /kdata

CMD ["node"]
ENTRYPOINT ["/app/usmint","--home","/kdata"]