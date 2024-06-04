# Use the official Node.js image.
FROM node:18-alpine

# Create and change to the /app directory.
WORKDIR /app

# Copy the package.json and package-lock.json files
COPY package*.json ./

# Install the dependencies
RUN npm install

# Copy the source code into the container
COPY . .

# Build the Next.js app
RUN npm run build

# Expose port 3000
EXPOSE 3000

# Command to run the app
CMD ["npm", "start"]
