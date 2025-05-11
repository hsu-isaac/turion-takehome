FROM node:18-alpine

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install ajv and its dependencies first
RUN npm install ajv@8.12.0 ajv-keywords@5.1.0 --legacy-peer-deps

# Install remaining dependencies
RUN npm install --legacy-peer-deps

# Copy source code
COPY . .

# Expose port
EXPOSE 3001

# Start development server
CMD ["npm", "start"] 