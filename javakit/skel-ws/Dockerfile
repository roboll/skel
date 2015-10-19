###############################################################################
# {{ .group }}/{{ .name}}-ws
###############################################################################
FROM jeanblanchard/java:serverjre-8

RUN apk add --update bash && rm -rf /var/cache/apk/*

ADD ./{{ .name }}-ws-*.tar /dist

RUN ln -s /dist/{{ .name }}-ws-* /opt/app

EXPOSE 8080 8081
WORKDIR /opt/app
ENTRYPOINT ["bin/{{ .name }}-ws", "server", "{{ .name }}-ws.yaml"]
