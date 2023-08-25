FROM node:14

WORKDIR /app

RUN npm install -g @mockoon/cli

COPY ./uniqueness-service-mock/initializerJson.json /app/

ENTRYPOINT ["mockoon-cli"]
