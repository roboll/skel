###############################################################################
# {{ .name}}-rpc
###############################################################################
FROM jeanblanchard/java:serverjre-8

RUN apk add --update bash && rm -rf /var/cache/apk/*

ADD ./{{ .name }}-rpc-*.tar /dist

RUN ln -s /dist/{{ .name }}-rpc-* /opt/app

EXPOSE 8080
ENTRYPOINT ["/opt/app/bin/{{ .name }}-rpc"]
