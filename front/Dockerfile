FROM node:22-alpine AS build
RUN npm install -g pnpm
WORKDIR /app
COPY package.json ./
RUN pnpm install
COPY . .
ARG VITE_API_URL=/api
ENV VITE_API_URL=$VITE_API_URL
RUN pnpm build

FROM node:22-alpine AS runtime
WORKDIR /app
COPY --from=build /app/dist ./dist
COPY --from=build /app/node_modules ./node_modules
COPY --from=build /app/package.json ./
COPY server.prod.mjs ./
EXPOSE 3000
ENV PORT=3000
CMD ["node", "server.prod.mjs"]
