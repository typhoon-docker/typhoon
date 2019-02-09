# === Build environment ===
FROM node:8.15-alpine as builder

# Setup
RUN mkdir /typhoon
WORKDIR /typhoon
ENV PATH /typhoon/node_modules/.bin:$PATH

# Package requirements
COPY package.json /typhoon/package.json
COPY yarn.lock /typhoon/yarn.lock

RUN yarn --pure-lockfile

# Source code
COPY . /typhoon

# Build the code
RUN yarn build

# === Production environment ===
FROM nginx:1.15-alpine
COPY --from=builder /typhoon/build /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]