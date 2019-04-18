FROM node:11.10.0 AS builder
COPY . /app
WORKDIR /app
RUN yarn install
RUN yarn run build

FROM wlchn/gostatic:latest
ENV CONFIG_FILE_PATH /srv/http
COPY --from=builder /app/build /srv/http
COPY ./env.sh /env.sh
ENTRYPOINT ["sh", "/env.sh"]
CMD ["/goStatic"]