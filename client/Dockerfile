# Multi-stage Dockerfile for Smart Log Viewer Client
# Supports both development and production builds

# ========================================
# DEVELOPMENT STAGE
# ========================================
FROM node:18-alpine AS development

# Install development dependencies
RUN apk add --no-cache \
    git \
    curl \
    ca-certificates

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install ALL dependencies (including dev dependencies)
RUN npm ci

# Copy source code
COPY . .

# Expose port for Vite dev server
EXPOSE 3000

# Development command with hot reload
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]

# ========================================
# BUILDER STAGE
# ========================================
FROM node:18-alpine AS builder

# Install build dependencies
RUN apk add --no-cache \
    git \
    ca-certificates

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install ALL dependencies (including dev dependencies for build)
RUN npm ci

# Copy source code
COPY . .

# Build the application
RUN npm run build

# ========================================
# PRODUCTION STAGE
# ========================================
FROM nginx:alpine AS production

# Install runtime dependencies
RUN apk add --no-cache \
    curl \
    ca-certificates

# Copy built files from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy custom nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Note: Running as root for now to avoid permission issues
# In production, you should implement proper user management

# Expose port 80
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD curl -f http://localhost:80/ || exit 1

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
