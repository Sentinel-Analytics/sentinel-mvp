# Use Node.js 22 on Debian (safer for Rollup prebuilt binaries)
FROM node:22-slim

# Enable Yarn via corepack
RUN corepack enable

# Set working directory
WORKDIR /app

# Copy package.json first
COPY package.json ./

# Install dependencies (this will generate yarn.lock inside container)
RUN yarn install

# Copy the rest of your app
COPY . .

# Expose Vite default port
EXPOSE 5173

# Run dev server
CMD ["yarn", "dev", "--host"]
