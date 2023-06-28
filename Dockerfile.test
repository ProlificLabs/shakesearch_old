FROM node:14 AS test-build

WORKDIR /app

COPY package.json .
COPY package-lock.json .

RUN npm ci

COPY test.js .

FROM golang:1.16 AS go-build

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o shakesearch

FROM zenika/alpine-chrome:with-puppeteer

WORKDIR /app

COPY --from=test-build /app/node_modules ./node_modules
COPY --from=test-build /app/test.js ./test.js
COPY --from=go-build /app/shakesearch ./shakesearch
COPY --from=go-build /app/static ./static
COPY --from=go-build /app/completeworks.txt ./completeworks.txt

ENV PORT=3002

CMD ["sh", "-c", "./shakesearch & sleep 1 && npx mocha test.js"]