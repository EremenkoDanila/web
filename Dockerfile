FROM node AS development

WORKDIR /software_app

COPY package*.json .

RUN npm install

COPY . .  