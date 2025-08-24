# Docker Guide for Python Applications

## Table of Contents
1. [Introduction](#introduction)
2. [Basic Python Application](#basic-python-application)
3. [Creating a Dockerfile](#creating-a-dockerfile)
4. [Docker Compose Setup](#docker-compose-setup)
5. [Web Application with Flask](#web-application-with-flask)
6. [Django Application](#django-application)
7. [FastAPI Application](#fastapi-application)
8. [Database Integration](#database-integration)
9. [Development vs Production](#development-vs-production)
10. [Microservices Architecture](#microservices-architecture)
11. [Advanced Topics](#advanced-topics)
12. [Best Practices](#best-practices)

## Introduction

Dockerizing Python applications offers several benefits:
- Consistent runtime environment across development, testing, and production
- Easy dependency management and isolation
- Simplified deployment and scaling
- Better resource utilization

This guide covers containerizing various types of Python applications using Docker and Docker Compose.

## Basic Python Application

Let's start with a simple Python script:

```python
# app.py
import requests
import time

def fetch_weather(city):
    """Fetch weather data for a given city."""
    api_key = "demo_key"
    url = f"http://api.openweathermap.org/data/2.5/weather?q={city}&appid={api_key}"
    
    try:
        response = requests.get(url)
        data = response.json()
        return data.get('weather', [{}])[0].get('description', 'Unknown')
    except Exception as e:
        return f"Error: {e}"

def main():
    cities = ["London", "Paris", "Tokyo"]
    
    print("Weather Report")
    print("-" * 20)
    
    for city in cities:
        weather = fetch_weather(city)
        print(f"{city}: {weather}")
        time.sleep(1)

if __name__ == "__main__":
    main()
```

Create a requirements file:

```txt
# requirements.txt
requests>=2.28.0
```

## Creating a Dockerfile

### Basic Dockerfile

```dockerfile
# Dockerfile
FROM python:3.11-slim

# Set working directory
WORKDIR /app

# Copy requirements first (for better caching)
COPY requirements.txt .

# Install dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY . .

# Run the application
CMD ["python", "app.py"]
```

### Multi-stage Dockerfile (Production-ready)

```dockerfile
# Multi-stage Dockerfile
FROM python:3.11-slim as builder

# Install build dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# Create virtual environment
RUN python -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"

# Copy and install requirements
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Production stage
FROM python:3.11-slim as production

# Copy virtual environment from builder
COPY --from=builder /opt/venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"

# Create non-root user
RUN useradd --create-home --shell /bin/bash appuser

# Set working directory
WORKDIR /app

# Copy application
COPY --chown=appuser:appuser . .

# Switch to non-root user
USER appuser

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD python -c "import requests; print('healthy')" || exit 1

# Run application
CMD ["python", "app.py"]
```

## Docker Compose Setup

### Basic docker-compose.yml

```yaml
# docker-compose.yml
version: '3.8'

services:
  python-app:
    build: .
    container_name: my-python-app
    environment:
      - PYTHONUNBUFFERED=1
    volumes:
      - .:/app
    networks:
      - python-network

networks:
  python-network:
    driver: bridge
```

### Development docker-compose.yml

```yaml
# docker-compose.dev.yml
version: '3.8'

services:
  python-app:
    build:
      context: .
      dockerfile: Dockerfile.dev
    container_name: python-app-dev
    environment:
      - PYTHONUNBUFFERED=1
      - FLASK_ENV=development
      - FLASK_DEBUG=1
    volumes:
      - .:/app
      - /app/__pycache__
    ports:
      - "5000:5000"
    networks:
      - dev-network

networks:
  dev-network:
    driver: bridge
```

## Web Application with Flask

### Flask Application

```python
# flask_app.py
from flask import Flask, jsonify, request
import os
import redis
import json

app = Flask(__name__)

# Redis connection
redis_client = redis.Redis(
    host=os.getenv('REDIS_HOST', 'localhost'),
    port=int(os.getenv('REDIS_PORT', 6379)),
    decode_responses=True
)

@app.route('/')
def home():
    return jsonify({"message": "Welcome to Flask App", "status": "running"})

@app.route('/health')
def health():
    return jsonify({"status": "healthy"}), 200

@app.route('/cache/<key>')
def get_cache(key):
    try:
        value = redis_client.get(key)
        if value:
            return jsonify({"key": key, "value": value})
        return jsonify({"error": "Key not found"}), 404
    except Exception as e:
        return jsonify({"error": str(e)}), 500

@app.route('/cache/<key>', methods=['POST'])
def set_cache(key):
    try:
        data = request.get_json()
        value = data.get('value')
        redis_client.set(key, value, ex=3600)  # 1 hour expiry
        return jsonify({"message": f"Set {key} = {value}"})
    except Exception as e:
        return jsonify({"error": str(e)}), 500

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, debug=True)
```

### Flask Requirements

```txt
# requirements.txt
Flask>=2.3.0
redis>=4.5.0
gunicorn>=20.1.0
```

### Flask Dockerfile

```dockerfile
# Dockerfile.flask
FROM python:3.11-slim

WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    && rm -rf /var/lib/apt/lists/*

# Copy and install Python dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application
COPY . .

# Create non-root user
RUN useradd --create-home --shell /bin/bash flaskuser
RUN chown -R flaskuser:flaskuser /app
USER flaskuser

# Expose port
EXPOSE 5000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:5000/health || exit 1

# Run with gunicorn
CMD ["gunicorn", "--bind", "0.0.0.0:5000", "--workers", "4", "flask_app:app"]
```

### Flask Docker Compose

```yaml
# docker-compose.flask.yml
version: '3.8'

services:
  flask-app:
    build:
      context: .
      dockerfile: Dockerfile.flask
    container_name: flask-web-app
    ports:
      - "5000:5000"
    environment:
      - FLASK_ENV=production
      - REDIS_HOST=redis
      - REDIS_PORT=6379
    depends_on:
      - redis
    networks:
      - flask-network
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    container_name: flask-redis
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - flask-network
    restart: unless-stopped

volumes:
  redis_data:

networks:
  flask-network:
    driver: bridge
```

## Django Application

### Django Settings

```python
# settings.py (Django configuration for Docker)
import os
from pathlib import Path

BASE_DIR = Path(__file__).resolve().parent.parent

SECRET_KEY = os.getenv('SECRET_KEY', 'your-secret-key-here')
DEBUG = os.getenv('DEBUG', 'False') == 'True'
ALLOWED_HOSTS = os.getenv('ALLOWED_HOSTS', 'localhost').split(',')

DATABASES = {
    'default': {
        'ENGINE': 'django.db.backends.postgresql',
        'NAME': os.getenv('DB_NAME', 'django_db'),
        'USER': os.getenv('DB_USER', 'postgres'),
        'PASSWORD': os.getenv('DB_PASSWORD', 'postgres'),
        'HOST': os.getenv('DB_HOST', 'localhost'),
        'PORT': os.getenv('DB_PORT', '5432'),
    }
}

CACHES = {
    'default': {
        'BACKEND': 'django_redis.cache.RedisCache',
        'LOCATION': f"redis://{os.getenv('REDIS_HOST', 'localhost')}:{os.getenv('REDIS_PORT', '6379')}/1",
        'OPTIONS': {
            'CLIENT_CLASS': 'django_redis.client.DefaultClient',
        }
    }
}
```

### Django Dockerfile

```dockerfile
# Dockerfile.django
FROM python:3.11-slim

# Set environment variables
ENV PYTHONUNBUFFERED=1
ENV PYTHONDONTWRITEBYTECODE=1

WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    postgresql-client \
    build-essential \
    libpq-dev \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements and install Python dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy project
COPY . .

# Create non-root user
RUN useradd --create-home --shell /bin/bash djangouser
RUN chown -R djangouser:djangouser /app
USER djangouser

# Collect static files
RUN python manage.py collectstatic --noinput

# Expose port
EXPOSE 8000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8000/health/ || exit 1

# Run with gunicorn
CMD ["gunicorn", "--bind", "0.0.0.0:8000", "--workers", "4", "myproject.wsgi:application"]
```

### Django Docker Compose

```yaml
# docker-compose.django.yml
version: '3.8'

services:
  django-app:
    build:
      context: .
      dockerfile: Dockerfile.django
    container_name: django-web-app
    ports:
      - "8000:8000"
    environment:
      - DEBUG=False
      - SECRET_KEY=your-production-secret-key
      - DB_NAME=django_db
      - DB_USER=postgres
      - DB_PASSWORD=postgres123
      - DB_HOST=postgres
      - DB_PORT=5432
      - REDIS_HOST=redis
      - REDIS_PORT=6379
    depends_on:
      - postgres
      - redis
    volumes:
      - static_volume:/app/staticfiles
      - media_volume:/app/media
    networks:
      - django-network
    restart: unless-stopped

  postgres:
    image: postgres:15-alpine
    container_name: django-postgres
    environment:
      - POSTGRES_DB=django_db
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres123
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - django-network
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    container_name: django-redis
    volumes:
      - redis_data:/data
    networks:
      - django-network
    restart: unless-stopped

  nginx:
    image: nginx:alpine
    container_name: django-nginx
    ports:
      - "80:80"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - static_volume:/app/staticfiles
      - media_volume:/app/media
    depends_on:
      - django-app
    networks:
      - django-network
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:
  static_volume:
  media_volume:

networks:
  django-network:
    driver: bridge
```

## FastAPI Application

### FastAPI Application

```python
# fastapi_app.py
from fastapi import FastAPI, HTTPException, Depends
from sqlalchemy import create_engine, Column, Integer, String
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker, Session
import os
import redis
from typing import List, Optional

app = FastAPI(title="FastAPI Docker Demo", version="1.0.0")

# Database setup
DATABASE_URL = f"postgresql://{os.getenv('DB_USER', 'postgres')}:{os.getenv('DB_PASSWORD', 'postgres')}@{os.getenv('DB_HOST', 'localhost')}:{os.getenv('DB_PORT', '5432')}/{os.getenv('DB_NAME', 'fastapi_db')}"

engine = create_engine(DATABASE_URL)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)
Base = declarative_base()

# Redis setup
redis_client = redis.Redis(
    host=os.getenv('REDIS_HOST', 'localhost'),
    port=int(os.getenv('REDIS_PORT', 6379)),
    decode_responses=True
)

# Models
class User(Base):
    __tablename__ = "users"
    
    id = Column(Integer, primary_key=True, index=True)
    name = Column(String, index=True)
    email = Column(String, unique=True, index=True)

# Create tables
Base.metadata.create_all(bind=engine)

# Dependency
def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()

@app.get("/")
async def root():
    return {"message": "FastAPI with Docker", "status": "running"}

@app.get("/health")
async def health():
    return {"status": "healthy"}

@app.post("/users/")
async def create_user(name: str, email: str, db: Session = Depends(get_db)):
    db_user = User(name=name, email=email)
    db.add(db_user)
    db.commit()
    db.refresh(db_user)
    
    # Cache user in Redis
    redis_client.set(f"user:{db_user.id}", f"{name}:{email}", ex=3600)
    
    return {"id": db_user.id, "name": db_user.name, "email": db_user.email}

@app.get("/users/{user_id}")
async def get_user(user_id: int, db: Session = Depends(get_db)):
    # Try Redis first
    cached_user = redis_client.get(f"user:{user_id}")
    if cached_user:
        name, email = cached_user.split(":")
        return {"id": user_id, "name": name, "email": email, "source": "cache"}
    
    # Fall back to database
    user = db.query(User).filter(User.id == user_id).first()
    if user is None:
        raise HTTPException(status_code=404, detail="User not found")
    
    return {"id": user.id, "name": user.name, "email": user.email, "source": "database"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
```

### FastAPI Requirements

```txt
# requirements.txt
fastapi>=0.100.0
uvicorn[standard]>=0.22.0
sqlalchemy>=2.0.0
psycopg2-binary>=2.9.0
redis>=4.5.0
```

### FastAPI Docker Compose

```yaml
# docker-compose.fastapi.yml
version: '3.8'

services:
  fastapi-app:
    build:
      context: .
      dockerfile: Dockerfile.fastapi
    container_name: fastapi-web-app
    ports:
      - "8000:8000"
    environment:
      - DB_NAME=fastapi_db
      - DB_USER=postgres
      - DB_PASSWORD=postgres123
      - DB_HOST=postgres
      - DB_PORT=5432
      - REDIS_HOST=redis
      - REDIS_PORT=6379
    depends_on:
      - postgres
      - redis
    networks:
      - fastapi-network
    restart: unless-stopped

  postgres:
    image: postgres:15-alpine
    container_name: fastapi-postgres
    environment:
      - POSTGRES_DB=fastapi_db
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres123
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - fastapi-network
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    container_name: fastapi-redis
    volumes:
      - redis_data:/data
    networks:
      - fastapi-network
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:

networks:
  fastapi-network:
    driver: bridge
```

## Database Integration

### PostgreSQL with SQLAlchemy

```python
# database.py
from sqlalchemy import create_engine
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
import os

DATABASE_URL = f"postgresql://{os.getenv('DB_USER')}:{os.getenv('DB_PASSWORD')}@{os.getenv('DB_HOST')}:{os.getenv('DB_PORT')}/{os.getenv('DB_NAME')}"

engine = create_engine(DATABASE_URL)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)
Base = declarative_base()

def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()
```

### MongoDB Integration

```python
# mongo_client.py
from pymongo import MongoClient
import os

def get_mongo_client():
    mongo_url = f"mongodb://{os.getenv('MONGO_USER')}:{os.getenv('MONGO_PASSWORD')}@{os.getenv('MONGO_HOST')}:{os.getenv('MONGO_PORT')}/{os.getenv('MONGO_DB')}"
    return MongoClient(mongo_url)
```

## Development vs Production

### Development Configuration

```yaml
# docker-compose.dev.yml
version: '3.8'

services:
  python-app:
    build:
      context: .
      dockerfile: Dockerfile.dev
    volumes:
      - .:/app
      - /app/__pycache__
    environment:
      - FLASK_ENV=development
      - DEBUG=True
    ports:
      - "5000:5000"
    command: ["python", "-m", "flask", "run", "--host=0.0.0.0", "--reload"]
```

### Production Configuration

```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  python-app:
    build:
      context: .
      dockerfile: Dockerfile.prod
    environment:
      - FLASK_ENV=production
      - DEBUG=False
    ports:
      - "80:8000"
    command: ["gunicorn", "--bind", "0.0.0.0:8000", "--workers", "4", "app:app"]
    restart: unless-stopped
```

## Microservices Architecture

### API Gateway Service

```python
# gateway.py
from flask import Flask, request, jsonify
import requests
import os

app = Flask(__name__)

SERVICES = {
    'user': os.getenv('USER_SERVICE_URL', 'http://user-service:5000'),
    'product': os.getenv('PRODUCT_SERVICE_URL', 'http://product-service:5000'),
    'order': os.getenv('ORDER_SERVICE_URL', 'http://order-service:5000')
}

@app.route('/<service>/<path:path>', methods=['GET', 'POST', 'PUT', 'DELETE'])
def proxy(service, path):
    if service not in SERVICES:
        return jsonify({'error': 'Service not found'}), 404
    
    url = f"{SERVICES[service]}/{path}"
    
    try:
        response = requests.request(
            method=request.method,
            url=url,
            headers=dict(request.headers),
            data=request.get_data(),
            params=request.args,
            timeout=30
        )
        return response.content, response.status_code, dict(response.headers)
    except requests.exceptions.RequestException as e:
        return jsonify({'error': str(e)}), 500

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8000)
```

### Microservices Docker Compose

```yaml
# docker-compose.microservices.yml
version: '3.8'

services:
  gateway:
    build:
      context: ./gateway
    ports:
      - "8000:8000"
    environment:
      - USER_SERVICE_URL=http://user-service:5000
      - PRODUCT_SERVICE_URL=http://product-service:5000
      - ORDER_SERVICE_URL=http://order-service:5000
    depends_on:
      - user-service
      - product-service
      - order-service
    networks:
      - microservices

  user-service:
    build:
      context: ./user-service
    environment:
      - DB_HOST=postgres
      - REDIS_HOST=redis
    depends_on:
      - postgres
      - redis
    networks:
      - microservices

  product-service:
    build:
      context: ./product-service
    environment:
      - DB_HOST=postgres
      - REDIS_HOST=redis
    depends_on:
      - postgres
      - redis
    networks:
      - microservices

  order-service:
    build:
      context: ./order-service
    environment:
      - DB_HOST=postgres
      - REDIS_HOST=redis
      - USER_SERVICE_URL=http://user-service:5000
      - PRODUCT_SERVICE_URL=http://product-service:5000
    depends_on:
      - postgres
      - redis
    networks:
      - microservices

  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_DB=microservices_db
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres123
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - microservices

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data
    networks:
      - microservices

volumes:
  postgres_data:
  redis_data:

networks:
  microservices:
    driver: bridge
```

## Advanced Topics

### Health Checks and Monitoring

```python
# health.py
from flask import Flask, jsonify
import psutil
import redis
import psycopg2
import os

app = Flask(__name__)

@app.route('/health')
def health():
    checks = {
        'status': 'healthy',
        'timestamp': time.time(),
        'checks': {}
    }
    
    # Database check
    try:
        conn = psycopg2.connect(
            host=os.getenv('DB_HOST'),
            database=os.getenv('DB_NAME'),
            user=os.getenv('DB_USER'),
            password=os.getenv('DB_PASSWORD')
        )
        conn.close()
        checks['checks']['database'] = 'healthy'
    except Exception as e:
        checks['checks']['database'] = f'unhealthy: {str(e)}'
        checks['status'] = 'unhealthy'
    
    # Redis check
    try:
        r = redis.Redis(host=os.getenv('REDIS_HOST'))
        r.ping()
        checks['checks']['redis'] = 'healthy'
    except Exception as e:
        checks['checks']['redis'] = f'unhealthy: {str(e)}'
        checks['status'] = 'unhealthy'
    
    # System resources
    checks['checks']['cpu_percent'] = psutil.cpu_percent()
    checks['checks']['memory_percent'] = psutil.virtual_memory().percent
    
    status_code = 200 if checks['status'] == 'healthy' else 503
    return jsonify(checks), status_code
```

### Logging Configuration

```python
# logging_config.py
import logging
import os
from pythonjsonlogger import jsonlogger

def setup_logging():
    log_level = os.getenv('LOG_LEVEL', 'INFO').upper()
    
    # JSON formatter for structured logging
    formatter = jsonlogger.JsonFormatter(
        '%(asctime)s %(name)s %(levelname)s %(message)s'
    )
    
    # Console handler
    console_handler = logging.StreamHandler()
    console_handler.setFormatter(formatter)
    
    # Configure root logger
    logging.basicConfig(
        level=getattr(logging, log_level),
        handlers=[console_handler]
    )
    
    return logging.getLogger(__name__)
```

### CI/CD Integration

```dockerfile
# Dockerfile.ci
FROM python:3.11-slim as test

WORKDIR /app
COPY requirements.txt requirements-dev.txt ./
RUN pip install -r requirements-dev.txt

COPY . .
RUN python -m pytest tests/ -v --cov=app --cov-report=xml

FROM python:3.11-slim as production

WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY --from=test /app .
RUN useradd --create-home --shell /bin/bash appuser
USER appuser

CMD ["gunicorn", "--bind", "0.0.0.0:8000", "app:app"]
```

## Best Practices

### 1. Multi-stage Builds
Use multi-stage builds to reduce image size and improve security.

### 2. Security
- Use non-root users
- Scan images for vulnerabilities
- Use secrets management
- Keep base images updated

### 3. Environment Variables
Store configuration in environment variables, never in code.

### 4. Health Checks
Implement proper health checks for all services.

### 5. Logging
Use structured logging with proper log levels.

### 6. Resource Limits
Set appropriate CPU and memory limits.

```yaml
# Resource limits example
services:
  python-app:
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M
```

### 7. Volume Management
Use named volumes for persistent data and bind mounts for development.

### 8. Network Security
Use custom networks and avoid exposing unnecessary ports.

### 9. Image Optimization
- Use appropriate base images
- Minimize layers
- Use .dockerignore
- Cache dependencies properly

### 10. Monitoring and Observability
Implement metrics collection, distributed tracing, and proper error handling.

This guide provides a comprehensive overview of dockerizing Python applications with various frameworks and patterns. Adapt the examples to your specific use case and requirements.