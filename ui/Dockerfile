# -----------------------------------------------------
# BUILD STAGE
FROM node:latest AS build
# Create Path-Structure
RUN mkdir /build
# Copy Source-Code
COPY configurj-ui /build
# Setting Workdir for BUILD-Process
WORKDIR /build
# Install Node Dependencies
RUN npm install
# Start Build Process of Node
RUN npm run build
# -----------------------------------------------------
# FINAL STAGE
FROM nginx:latest
# Copy the Build Resources
COPY --from=build /build/dist/configurj-ui /usr/share/nginx/html
# Expose HTTP/80
EXPOSE 80
