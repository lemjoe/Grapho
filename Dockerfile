FROM debian:latest

RUN mkdir -p /opt/Grapho

WORKDIR /opt/Grapho

COPY grapho .env.default LICENSE README.md /opt/Grapho/
COPY ./images /opt/Grapho/images/
COPY ./lang /opt/Grapho/lang/
COPY ./lib /opt/Grapho/lib/

RUN chmod +x grapho
RUN mkdir -p db articles

ENV JWT_SECRET=defaultsecret
ENV ADMIN_PASSWD=Admin321
ENV MAIN_LOG=/opt/Grapho/main.log
ENV DB_TYPE=cloverdb
ENV DB_PATH=/opt/Grapho/db

EXPOSE 4007

ENTRYPOINT ["/opt/Grapho/grapho"]
