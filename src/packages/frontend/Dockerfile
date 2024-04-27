# Use Node.js LTS version as base image for frontend
FROM node:lts AS frontend

# Set working directory for frontend
WORKDIR /app/packages/frontend

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies for frontend
RUN npm install

# Copy the rest of the frontend application
COPY . .

# Build the frontend 
RUN npm run build

# Expose the port used by the frontend application
EXPOSE 3000

LABEL maintainer="GoGoPowerRangers"
LABEL description="This is a frontend Docker image for application."

# Start the frontend server
CMD ["npm", "start"]
