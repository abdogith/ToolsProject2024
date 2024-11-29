# base image
FROM node:16-alpine AS build

# working directory
WORKDIR /app

COPY package.json package-lock.json ./

# installing dependencies
RUN npm install


# copying app
COPY . .

# building the app
RUN npm run build

FROM nginx:alpine AS frontend-production
COPY --from=build /app/build /usr/share/nginx/html
# port 
EXPOSE 80

# to start the app 
# base image
FROM node:16-alpine AS frontend-build

# working directory
WORKDIR /app

COPY package.json package-lock.json./

# installing dependencies
RUN npm install

# copying app
COPY . .

# building the app
RUN npm run build

FROM nginx:alpine AS production
COPY --from=build /app/build /usr/share/nginx/html
# port 
EXPOSE 3000

# to start the app 
CMD ["nginx", "-g", "daemon off;"]

