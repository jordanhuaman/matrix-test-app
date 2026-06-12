FROM node:22-alpine AS build
WORKDIR /app
COPY package*.json tsconfig*.json ./
RUN npm ci
COPY src/ ./src/
RUN npx tsc

FROM node:22-alpine AS runtime
WORKDIR /app
COPY --from=build /app/dist ./dist
COPY --from=build /app/node_modules ./node_modules
EXPOSE 3001
ENV PORT=3001
CMD ["node", "dist/index.js"]
