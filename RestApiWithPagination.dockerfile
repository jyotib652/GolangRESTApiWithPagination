FROM alpine:latest

RUN mkdir /app

COPY ./cmd/api/restApiWithPaginationApp /app

CMD [ "/app/restApiWithPaginationApp" ]