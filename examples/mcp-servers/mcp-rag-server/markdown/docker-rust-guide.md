# Docker Guide for Rust Applications

## Table of Contents
1. [Introduction](#introduction)
2. [Basic Rust Application](#basic-rust-application)
3. [Creating a Dockerfile](#creating-a-dockerfile)
4. [Docker Compose Setup](#docker-compose-setup)
5. [Web API with Axum](#web-api-with-axum)
6. [Web Application with Actix-web](#web-application-with-actix-web)
7. [Database Integration](#database-integration)
8. [Development vs Production](#development-vs-production)
9. [Microservices Architecture](#microservices-architecture)
10. [Advanced Topics](#advanced-topics)
11. [Best Practices](#best-practices)

## Introduction

Dockerizing Rust applications offers several advantages:
- Consistent build environment across different systems
- Fast, lightweight containers due to Rust's compiled nature
- Excellent security with minimal attack surface
- Easy deployment and scaling
- Reproducible builds

This guide covers containerizing various types of Rust applications using Docker and Docker Compose.

## Basic Rust Application

Let's start with a simple Rust CLI application:

```toml
# Cargo.toml
[package]
name = "rust-docker-demo"
version = "0.1.0"
edition = "2021"

[dependencies]
reqwest = { version = "0.11", features = ["json"] }
tokio = { version = "1.0", features = ["full"] }
serde = { version = "1.0", features = ["derive"] }
serde_json = "1.0"
clap = { version = "4.0", features = ["derive"] }
```

```rust
// src/main.rs
use reqwest;
use serde::{Deserialize, Serialize};
use std::error::Error;
use clap::{Arg, Command};

#[derive(Debug, Deserialize, Serialize)]
struct WeatherResponse {
    weather: Vec<Weather>,
    main: Main,
    name: String,
}

#[derive(Debug, Deserialize, Serialize)]
struct Weather {
    main: String,
    description: String,
}

#[derive(Debug, Deserialize, Serialize)]
struct Main {
    temp: f64,
    humidity: i32,
}

async fn fetch_weather(city: &str, api_key: &str) -> Result<WeatherResponse, Box<dyn Error>> {
    let url = format!(
        "http://api.openweathermap.org/data/2.5/weather?q={}&appid={}&units=metric",
        city, api_key
    );
    
    let response = reqwest::get(&url).await?;
    let weather: WeatherResponse = response.json().await?;
    
    Ok(weather)
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    let matches = Command::new("Weather CLI")
        .version("1.0")
        .author("Your Name")
        .about("Fetches weather information")
        .arg(
            Arg::new("city")
                .short('c')
                .long("city")
                .value_name("CITY")
                .help("City to fetch weather for")
                .required(true),
        )
        .arg(
            Arg::new("api-key")
                .short('k')
                .long("api-key")
                .value_name("API_KEY")
                .help("OpenWeatherMap API key")
                .default_value("demo_key"),
        )
        .get_matches();

    let city = matches.get_one::<String>("city").unwrap();
    let api_key = matches.get_one::<String>("api-key").unwrap();

    println!("Fetching weather for {}...", city);

    match fetch_weather(city, api_key).await {
        Ok(weather) => {
            println!("\n🌤️  Weather Report for {}", weather.name);
            println!("📊 Temperature: {:.1}°C", weather.main.temp);
            println!("💧 Humidity: {}%", weather.main.humidity);
            println!("☁️  Conditions: {}", weather.weather[0].description);
        }
        Err(e) => {
            eprintln!("❌ Error fetching weather: {}", e);
            std::process::exit(1);
        }
    }

    Ok(())
}
```

## Creating a Dockerfile

### Basic Dockerfile

```dockerfile
# Dockerfile
FROM rust:1.75-slim as builder

# Install build dependencies
RUN apt-get update && apt-get install -y \
    pkg-config \
    libssl-dev \
    && rm -rf /var/lib/apt/lists/*

# Create app directory
WORKDIR /app

# Copy manifests
COPY Cargo.toml Cargo.lock ./

# Build dependencies (this will be cached)
RUN mkdir src && echo "fn main() {}" > src/main.rs
RUN cargo build --release
RUN rm -rf src

# Copy source code
COPY src ./src

# Build application
RUN touch src/main.rs
RUN cargo build --release

# Runtime stage
FROM debian:bookworm-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN useradd --create-home --shell /bin/bash rustuser

# Copy binary from builder stage
COPY --from=builder /app/target/release/rust-docker-demo /usr/local/bin/rust-docker-demo
RUN chmod +x /usr/local/bin/rust-docker-demo

# Switch to non-root user
USER rustuser

# Set entrypoint
ENTRYPOINT ["rust-docker-demo"]
```

### Multi-stage Optimized Dockerfile

```dockerfile
# Dockerfile.optimized
FROM rust:1.75-slim as chef
RUN cargo install cargo-chef
WORKDIR /app

# Prepare recipe
FROM chef as planner
COPY . .
RUN cargo chef prepare --recipe-path recipe.json

# Build dependencies
FROM chef as builder
COPY --from=planner /app/recipe.json recipe.json
RUN cargo chef cook --release --recipe-path recipe.json

# Build application
COPY . .
RUN cargo build --release

# Runtime stage using distroless for minimal attack surface
FROM gcr.io/distroless/cc-debian12

# Copy binary
COPY --from=builder /app/target/release/rust-docker-demo /usr/local/bin/rust-docker-demo

# Set user
USER 1000

# Set entrypoint
ENTRYPOINT ["/usr/local/bin/rust-docker-demo"]
```

### Development Dockerfile

```dockerfile
# Dockerfile.dev
FROM rust:1.75-slim

# Install development tools
RUN apt-get update && apt-get install -y \
    pkg-config \
    libssl-dev \
    git \
    && rm -rf /var/lib/apt/lists/*

# Install cargo-watch for hot reloading
RUN cargo install cargo-watch

WORKDIR /app

# Copy manifests
COPY Cargo.toml Cargo.lock ./

# Pre-build dependencies
RUN mkdir src && echo "fn main() {}" > src/main.rs
RUN cargo build
RUN rm -rf src

# Set up for development
EXPOSE 3000

# Default command for development
CMD ["cargo", "watch", "-x", "run"]
```

## Docker Compose Setup

### Basic docker-compose.yml

```yaml
# docker-compose.yml
version: '3.8'

services:
  rust-app:
    build: .
    container_name: rust-cli-app
    environment:
      - RUST_LOG=info
    networks:
      - rust-network
    command: ["--city", "London", "--api-key", "your-api-key"]

networks:
  rust-network:
    driver: bridge
```

### Development docker-compose.yml

```yaml
# docker-compose.dev.yml
version: '3.8'

services:
  rust-app:
    build:
      context: .
      dockerfile: Dockerfile.dev
    container_name: rust-app-dev
    environment:
      - RUST_LOG=debug
      - RUST_BACKTRACE=1
    volumes:
      - .:/app
      - cargo_cache:/usr/local/cargo/registry
      - target_cache:/app/target
    ports:
      - "3000:3000"
    networks:
      - dev-network

volumes:
  cargo_cache:
  target_cache:

networks:
  dev-network:
    driver: bridge
```

## Web API with Axum

### Axum Application

```toml
# Cargo.toml for Axum
[package]
name = "axum-docker-demo"
version = "0.1.0"
edition = "2021"

[dependencies]
axum = "0.7"
tokio = { version = "1.0", features = ["full"] }
tower = "0.4"
tower-http = { version = "0.5", features = ["cors", "trace"] }
serde = { version = "1.0", features = ["derive"] }
serde_json = "1.0"
sqlx = { version = "0.7", features = ["runtime-tokio-rustls", "postgres", "chrono", "uuid"] }
redis = { version = "0.24", features = ["tokio-comp"] }
uuid = { version = "1.0", features = ["v4", "serde"] }
chrono = { version = "0.4", features = ["serde"] }
tracing = "0.1"
tracing-subscriber = { version = "0.3", features = ["env-filter"] }
anyhow = "1.0"
```

```rust
// src/main.rs
use axum::{
    extract::{Path, State},
    http::StatusCode,
    response::Json,
    routing::{delete, get, post, put},
    Router,
};
use serde::{Deserialize, Serialize};
use sqlx::{PgPool, Row};
use std::sync::Arc;
use tower_http::cors::CorsLayer;
use tracing::{info, instrument};
use uuid::Uuid;

#[derive(Debug, Serialize, Deserialize)]
struct User {
    id: Uuid,
    name: String,
    email: String,
    created_at: chrono::DateTime<chrono::Utc>,
}

#[derive(Debug, Deserialize)]
struct CreateUser {
    name: String,
    email: String,
}

#[derive(Debug, Deserialize)]
struct UpdateUser {
    name: Option<String>,
    email: Option<String>,
}

#[derive(Clone)]
struct AppState {
    db: PgPool,
    redis: redis::Client,
}

#[instrument]
async fn health() -> Json<serde_json::Value> {
    Json(serde_json::json!({
        "status": "healthy",
        "timestamp": chrono::Utc::now(),
        "service": "axum-docker-demo"
    }))
}

#[instrument]
async fn get_users(State(state): State<Arc<AppState>>) -> Result<Json<Vec<User>>, StatusCode> {
    let users = sqlx::query_as!(
        User,
        "SELECT id, name, email, created_at FROM users ORDER BY created_at DESC"
    )
    .fetch_all(&state.db)
    .await
    .map_err(|e| {
        tracing::error!("Database error: {}", e);
        StatusCode::INTERNAL_SERVER_ERROR
    })?;

    Ok(Json(users))
}

#[instrument]
async fn get_user(
    Path(user_id): Path<Uuid>,
    State(state): State<Arc<AppState>>,
) -> Result<Json<User>, StatusCode> {
    // Try Redis cache first
    let mut redis_conn = state
        .redis
        .get_async_connection()
        .await
        .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    let cache_key = format!("user:{}", user_id);
    
    if let Ok(cached_user) = redis::cmd("GET")
        .arg(&cache_key)
        .query_async::<_, String>(&mut redis_conn)
        .await
    {
        if let Ok(user) = serde_json::from_str::<User>(&cached_user) {
            return Ok(Json(user));
        }
    }

    // Fallback to database
    let user = sqlx::query_as!(
        User,
        "SELECT id, name, email, created_at FROM users WHERE id = $1",
        user_id
    )
    .fetch_optional(&state.db)
    .await
    .map_err(|e| {
        tracing::error!("Database error: {}", e);
        StatusCode::INTERNAL_SERVER_ERROR
    })?
    .ok_or(StatusCode::NOT_FOUND)?;

    // Cache for 1 hour
    let _: () = redis::cmd("SETEX")
        .arg(&cache_key)
        .arg(3600)
        .arg(serde_json::to_string(&user).unwrap())
        .query_async(&mut redis_conn)
        .await
        .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(Json(user))
}

#[instrument]
async fn create_user(
    State(state): State<Arc<AppState>>,
    Json(payload): Json<CreateUser>,
) -> Result<Json<User>, StatusCode> {
    let user_id = Uuid::new_v4();
    let now = chrono::Utc::now();

    let user = sqlx::query_as!(
        User,
        "INSERT INTO users (id, name, email, created_at) VALUES ($1, $2, $3, $4) RETURNING id, name, email, created_at",
        user_id,
        payload.name,
        payload.email,
        now
    )
    .fetch_one(&state.db)
    .await
    .map_err(|e| {
        tracing::error!("Database error: {}", e);
        StatusCode::INTERNAL_SERVER_ERROR
    })?;

    Ok(Json(user))
}

#[instrument]
async fn update_user(
    Path(user_id): Path<Uuid>,
    State(state): State<Arc<AppState>>,
    Json(payload): Json<UpdateUser>,
) -> Result<Json<User>, StatusCode> {
    let user = sqlx::query_as!(
        User,
        r#"
        UPDATE users 
        SET 
            name = COALESCE($2, name),
            email = COALESCE($3, email)
        WHERE id = $1
        RETURNING id, name, email, created_at
        "#,
        user_id,
        payload.name,
        payload.email
    )
    .fetch_optional(&state.db)
    .await
    .map_err(|e| {
        tracing::error!("Database error: {}", e);
        StatusCode::INTERNAL_SERVER_ERROR
    })?
    .ok_or(StatusCode::NOT_FOUND)?;

    // Invalidate cache
    let mut redis_conn = state.redis.get_async_connection().await.ok();
    if let Some(ref mut conn) = redis_conn {
        let _: Result<(), _> = redis::cmd("DEL")
            .arg(format!("user:{}", user_id))
            .query_async(conn)
            .await;
    }

    Ok(Json(user))
}

#[instrument]
async fn delete_user(
    Path(user_id): Path<Uuid>,
    State(state): State<Arc<AppState>>,
) -> Result<StatusCode, StatusCode> {
    let result = sqlx::query!("DELETE FROM users WHERE id = $1", user_id)
        .execute(&state.db)
        .await
        .map_err(|e| {
            tracing::error!("Database error: {}", e);
            StatusCode::INTERNAL_SERVER_ERROR
        })?;

    if result.rows_affected() == 0 {
        return Err(StatusCode::NOT_FOUND);
    }

    // Remove from cache
    let mut redis_conn = state.redis.get_async_connection().await.ok();
    if let Some(ref mut conn) = redis_conn {
        let _: Result<(), _> = redis::cmd("DEL")
            .arg(format!("user:{}", user_id))
            .query_async(conn)
            .await;
    }

    Ok(StatusCode::NO_CONTENT)
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::init();

    let database_url = std::env::var("DATABASE_URL")
        .unwrap_or_else(|_| "postgresql://postgres:postgres@localhost:5432/axum_db".to_string());
    
    let redis_url = std::env::var("REDIS_URL")
        .unwrap_or_else(|_| "redis://localhost:6379".to_string());

    // Connect to database
    let db = PgPool::connect(&database_url).await?;
    
    // Run migrations
    sqlx::migrate!("./migrations").run(&db).await?;

    // Connect to Redis
    let redis_client = redis::Client::open(redis_url)?;

    let state = Arc::new(AppState {
        db,
        redis: redis_client,
    });

    let app = Router::new()
        .route("/health", get(health))
        .route("/users", get(get_users).post(create_user))
        .route("/users/:id", get(get_user).put(update_user).delete(delete_user))
        .layer(CorsLayer::permissive())
        .layer(tower_http::trace::TraceLayer::new_for_http())
        .with_state(state);

    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await?;
    
    info!("🚀 Server starting on http://0.0.0.0:3000");
    
    axum::serve(listener, app).await?;

    Ok(())
}
```

### Database Migration

```sql
-- migrations/001_create_users_table.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL,
    email VARCHAR UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at);
```

### Axum Dockerfile

```dockerfile
# Dockerfile.axum
FROM rust:1.75-slim as builder

# Install build dependencies
RUN apt-get update && apt-get install -y \
    pkg-config \
    libssl-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy manifests
COPY Cargo.toml Cargo.lock ./

# Build dependencies
RUN mkdir src && echo "fn main() {}" > src/main.rs
RUN cargo build --release
RUN rm -rf src

# Copy source and migrations
COPY src ./src
COPY migrations ./migrations

# Build application
RUN touch src/main.rs
RUN cargo build --release

# Runtime stage
FROM debian:bookworm-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN useradd --create-home --shell /bin/bash axumuser

# Copy binary and migrations
COPY --from=builder /app/target/release/axum-docker-demo /usr/local/bin/axum-docker-demo
COPY --from=builder /app/migrations /app/migrations

# Set permissions
RUN chmod +x /usr/local/bin/axum-docker-demo

# Switch to non-root user
USER axumuser

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:3000/health || exit 1

# Run application
CMD ["axum-docker-demo"]
```

### Axum Docker Compose

```yaml
# docker-compose.axum.yml
version: '3.8'

services:
  axum-app:
    build:
      context: .
      dockerfile: Dockerfile.axum
    container_name: axum-web-app
    ports:
      - "3000:3000"
    environment:
      - RUST_LOG=info
      - DATABASE_URL=postgresql://postgres:postgres123@postgres:5432/axum_db
      - REDIS_URL=redis://redis:6379
    depends_on:
      - postgres
      - redis
    networks:
      - axum-network
    restart: unless-stopped

  postgres:
    image: postgres:15-alpine
    container_name: axum-postgres
    environment:
      - POSTGRES_DB=axum_db
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres123
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - axum-network
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    container_name: axum-redis
    volumes:
      - redis_data:/data
    networks:
      - axum-network
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:

networks:
  axum-network:
    driver: bridge
```

## Web Application with Actix-web

### Actix-web Application

```toml
# Cargo.toml for Actix-web
[package]
name = "actix-docker-demo"
version = "0.1.0"
edition = "2021"

[dependencies]
actix-web = "4"
actix-cors = "0.6"
tokio = { version = "1.0", features = ["full"] }
serde = { version = "1.0", features = ["derive"] }
serde_json = "1.0"
sqlx = { version = "0.7", features = ["runtime-tokio-rustls", "postgres", "chrono", "uuid"] }
uuid = { version = "1.0", features = ["v4", "serde"] }
chrono = { version = "0.4", features = ["serde"] }
env_logger = "0.10"
log = "0.4"
```

```rust
// src/main.rs
use actix_cors::Cors;
use actix_web::{
    delete, get, middleware::Logger, post, put, web, App, HttpResponse, HttpServer, Result,
};
use serde::{Deserialize, Serialize};
use sqlx::{PgPool, Row};
use uuid::Uuid;

#[derive(Debug, Serialize, Deserialize)]
struct User {
    id: Uuid,
    name: String,
    email: String,
    created_at: chrono::DateTime<chrono::Utc>,
}

#[derive(Debug, Deserialize)]
struct CreateUser {
    name: String,
    email: String,
}

struct AppState {
    db: PgPool,
}

#[get("/health")]
async fn health() -> Result<HttpResponse> {
    Ok(HttpResponse::Ok().json(serde_json::json!({
        "status": "healthy",
        "timestamp": chrono::Utc::now(),
        "service": "actix-docker-demo"
    })))
}

#[get("/users")]
async fn get_users(data: web::Data<AppState>) -> Result<HttpResponse> {
    match sqlx::query_as!(
        User,
        "SELECT id, name, email, created_at FROM users ORDER BY created_at DESC"
    )
    .fetch_all(&data.db)
    .await
    {
        Ok(users) => Ok(HttpResponse::Ok().json(users)),
        Err(e) => {
            log::error!("Database error: {}", e);
            Ok(HttpResponse::InternalServerError().json("Database error"))
        }
    }
}

#[get("/users/{user_id}")]
async fn get_user(
    path: web::Path<Uuid>,
    data: web::Data<AppState>,
) -> Result<HttpResponse> {
    let user_id = path.into_inner();

    match sqlx::query_as!(
        User,
        "SELECT id, name, email, created_at FROM users WHERE id = $1",
        user_id
    )
    .fetch_optional(&data.db)
    .await
    {
        Ok(Some(user)) => Ok(HttpResponse::Ok().json(user)),
        Ok(None) => Ok(HttpResponse::NotFound().json("User not found")),
        Err(e) => {
            log::error!("Database error: {}", e);
            Ok(HttpResponse::InternalServerError().json("Database error"))
        }
    }
}

#[post("/users")]
async fn create_user(
    user_data: web::Json<CreateUser>,
    data: web::Data<AppState>,
) -> Result<HttpResponse> {
    let user_id = Uuid::new_v4();
    let now = chrono::Utc::now();

    match sqlx::query_as!(
        User,
        "INSERT INTO users (id, name, email, created_at) VALUES ($1, $2, $3, $4) RETURNING id, name, email, created_at",
        user_id,
        user_data.name,
        user_data.email,
        now
    )
    .fetch_one(&data.db)
    .await
    {
        Ok(user) => Ok(HttpResponse::Created().json(user)),
        Err(e) => {
            log::error!("Database error: {}", e);
            Ok(HttpResponse::InternalServerError().json("Database error"))
        }
    }
}

#[delete("/users/{user_id}")]
async fn delete_user(
    path: web::Path<Uuid>,
    data: web::Data<AppState>,
) -> Result<HttpResponse> {
    let user_id = path.into_inner();

    match sqlx::query!("DELETE FROM users WHERE id = $1", user_id)
        .execute(&data.db)
        .await
    {
        Ok(result) => {
            if result.rows_affected() > 0 {
                Ok(HttpResponse::NoContent().finish())
            } else {
                Ok(HttpResponse::NotFound().json("User not found"))
            }
        }
        Err(e) => {
            log::error!("Database error: {}", e);
            Ok(HttpResponse::InternalServerError().json("Database error"))
        }
    }
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    env_logger::init();

    let database_url = std::env::var("DATABASE_URL")
        .unwrap_or_else(|_| "postgresql://postgres:postgres@localhost:5432/actix_db".to_string());

    let pool = PgPool::connect(&database_url)
        .await
        .expect("Failed to connect to database");

    // Run migrations
    sqlx::migrate!("./migrations")
        .run(&pool)
        .await
        .expect("Failed to run migrations");

    log::info!("🚀 Server starting on http://0.0.0.0:8080");

    HttpServer::new(move || {
        let cors = Cors::default()
            .allow_any_origin()
            .allow_any_method()
            .allow_any_header();

        App::new()
            .app_data(web::Data::new(AppState { db: pool.clone() }))
            .wrap(cors)
            .wrap(Logger::default())
            .service(health)
            .service(get_users)
            .service(get_user)
            .service(create_user)
            .service(delete_user)
    })
    .bind("0.0.0.0:8080")?
    .run()
    .await
}
```

## Database Integration

### SQLx Configuration

```rust
// src/database.rs
use sqlx::{PgPool, postgres::PgPoolOptions};
use std::env;

pub async fn create_pool() -> Result<PgPool, sqlx::Error> {
    let database_url = env::var("DATABASE_URL")
        .expect("DATABASE_URL must be set");

    PgPoolOptions::new()
        .max_connections(10)
        .connect(&database_url)
        .await
}

pub async fn run_migrations(pool: &PgPool) -> Result<(), sqlx::Error> {
    sqlx::migrate!("./migrations").run(pool).await
}
```

### Redis Integration

```rust
// src/cache.rs
use redis::{Client, RedisResult, AsyncCommands};
use serde::{Serialize, Deserialize};
use std::env;

pub struct CacheManager {
    client: Client,
}

impl CacheManager {
    pub fn new() -> RedisResult<Self> {
        let redis_url = env::var("REDIS_URL")
            .unwrap_or_else(|_| "redis://localhost:6379".to_string());
        
        let client = Client::open(redis_url)?;
        Ok(CacheManager { client })
    }

    pub async fn get<T>(&self, key: &str) -> RedisResult<Option<T>>
    where
        T: for<'de> Deserialize<'de>,
    {
        let mut conn = self.client.get_async_connection().await?;
        let value: Option<String> = conn.get(key).await?;
        
        match value {
            Some(json) => Ok(serde_json::from_str(&json).ok()),
            None => Ok(None),
        }
    }

    pub async fn set<T>(&self, key: &str, value: &T, expiry: usize) -> RedisResult<()>
    where
        T: Serialize,
    {
        let mut conn = self.client.get_async_connection().await?;
        let json = serde_json::to_string(value).unwrap();
        conn.setex(key, expiry, json).await
    }

    pub async fn delete(&self, key: &str) -> RedisResult<()> {
        let mut conn = self.client.get_async_connection().await?;
        conn.del(key).await
    }
}
```

## Development vs Production

### Development Configuration

```yaml
# docker-compose.dev.yml
version: '3.8'

services:
  rust-app:
    build:
      context: .
      dockerfile: Dockerfile.dev
    volumes:
      - .:/app
      - cargo_cache:/usr/local/cargo/registry
      - target_cache:/app/target
    environment:
      - RUST_LOG=debug
      - RUST_BACKTRACE=full
    ports:
      - "3000:3000"
    depends_on:
      - postgres-dev
    command: ["cargo", "watch", "-x", "run"]

  postgres-dev:
    image: postgres:15-alpine
    environment:
      - POSTGRES_DB=dev_db
      - POSTGRES_USER=dev
      - POSTGRES_PASSWORD=dev
    ports:
      - "5432:5432"
    volumes:
      - postgres_dev_data:/var/lib/postgresql/data

volumes:
  cargo_cache:
  target_cache:
  postgres_dev_data:
```

### Production Configuration

```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  rust-app:
    build:
      context: .
      dockerfile: Dockerfile
    environment:
      - RUST_LOG=warn
      - DATABASE_URL=postgresql://prod_user:prod_pass@postgres:5432/prod_db
    ports:
      - "80:3000"
    depends_on:
      - postgres
    restart: unless-stopped
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 512M

  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_DB=prod_db
      - POSTGRES_USER=prod_user
      - POSTGRES_PASSWORD=prod_pass
    volumes:
      - postgres_prod_data:/var/lib/postgresql/data
    restart: unless-stopped

volumes:
  postgres_prod_data:
```

## Microservices Architecture

### API Gateway

```rust
// src/gateway.rs
use axum::{
    extract::{Path, Request},
    http::{Method, StatusCode, Uri},
    response::Response,
    routing::any,
    Router,
};
use hyper::Body;
use std::collections::HashMap;
use tower::ServiceExt;

#[derive(Clone)]
struct Gateway {
    services: HashMap<String, String>,
    client: reqwest::Client,
}

impl Gateway {
    fn new() -> Self {
        let mut services = HashMap::new();
        services.insert("users".to_string(), "http://user-service:3000".to_string());
        services.insert("products".to_string(), "http://product-service:3000".to_string());
        services.insert("orders".to_string(), "http://order-service:3000".to_string());

        Gateway {
            services,
            client: reqwest::Client::new(),
        }
    }

    async fn proxy_request(
        &self,
        service: &str,
        path: &str,
        method: Method,
        body: Body,
    ) -> Result<Response<Body>, StatusCode> {
        let service_url = self.services.get(service)
            .ok_or(StatusCode::NOT_FOUND)?;

        let url = format!("{}/{}", service_url, path);
        
        let request_builder = match method {
            Method::GET => self.client.get(&url),
            Method::POST => self.client.post(&url),
            Method::PUT => self.client.put(&url),
            Method::DELETE => self.client.delete(&url),
            _ => return Err(StatusCode::METHOD_NOT_ALLOWED),
        };

        // Forward the request
        match request_builder.send().await {
            Ok(response) => {
                let status = response.status();
                let body = response.bytes().await.map_err(|_| StatusCode::BAD_GATEWAY)?;
                
                Ok(Response::builder()
                    .status(status.as_u16())
                    .body(Body::from(body))
                    .unwrap())
            }
            Err(_) => Err(StatusCode::BAD_GATEWAY),
        }
    }
}

async fn handle_proxy(
    Path((service, path)): Path<(String, String)>,
    request: Request<Body>,
) -> Result<Response<Body>, StatusCode> {
    let gateway = Gateway::new();
    let (parts, body) = request.into_parts();
    
    gateway.proxy_request(&service, &path, parts.method, body).await
}

pub fn create_router() -> Router {
    Router::new()
        .route("/:service/*path", any(handle_proxy))
        .route("/health", axum::routing::get(|| async { "Gateway healthy" }))
}
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
      - "8080:8080"
    environment:
      - RUST_LOG=info
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
      - RUST_LOG=info
      - DATABASE_URL=postgresql://postgres:postgres@postgres:5432/users_db
      - REDIS_URL=redis://redis:6379
    depends_on:
      - postgres
      - redis
    networks:
      - microservices

  product-service:
    build:
      context: ./product-service
    environment:
      - RUST_LOG=info
      - DATABASE_URL=postgresql://postgres:postgres@postgres:5432/products_db
      - REDIS_URL=redis://redis:6379
    depends_on:
      - postgres
      - redis
    networks:
      - microservices

  order-service:
    build:
      context: ./order-service
    environment:
      - RUST_LOG=info
      - DATABASE_URL=postgresql://postgres:postgres@postgres:5432/orders_db
      - REDIS_URL=redis://redis:6379
      - USER_SERVICE_URL=http://user-service:3000
      - PRODUCT_SERVICE_URL=http://product-service:3000
    depends_on:
      - postgres
      - redis
    networks:
      - microservices

  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init-dbs.sql:/docker-entrypoint-initdb.d/init-dbs.sql
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

```rust
// src/health.rs
use axum::{extract::State, http::StatusCode, response::Json, routing::get, Router};
use serde_json::{json, Value};
use sqlx::PgPool;
use std::sync::Arc;
use std::time::SystemTime;

pub struct HealthChecker {
    db: PgPool,
}

impl HealthChecker {
    pub fn new(db: PgPool) -> Self {
        Self { db }
    }

    pub async fn check_database(&self) -> bool {
        sqlx::query("SELECT 1").execute(&self.db).await.is_ok()
    }

    pub async fn check_memory(&self) -> (u64, u64) {
        // Simple memory check
        let process = std::process::Command::new("ps")
            .args(&["-o", "rss,vsz", "-p", &std::process::id().to_string()])
            .output()
            .ok();

        match process {
            Some(output) => {
                let output_str = String::from_utf8_lossy(&output.stdout);
                let lines: Vec<&str> = output_str.lines().collect();
                if lines.len() > 1 {
                    let parts: Vec<&str> = lines[1].split_whitespace().collect();
                    if parts.len() >= 2 {
                        let rss = parts[0].parse().unwrap_or(0);
                        let vsz = parts[1].parse().unwrap_or(0);
                        return (rss, vsz);
                    }
                }
            }
            None => {}
        }
        (0, 0)
    }
}

async fn health_check(State(checker): State<Arc<HealthChecker>>) -> Result<Json<Value>, StatusCode> {
    let mut status = "healthy";
    let mut checks = json!({});

    // Database check
    let db_healthy = checker.check_database().await;
    checks["database"] = json!({
        "status": if db_healthy { "healthy" } else { "unhealthy" }
    });

    if !db_healthy {
        status = "unhealthy";
    }

    // Memory check
    let (rss, vsz) = checker.check_memory().await;
    checks["memory"] = json!({
        "rss_kb": rss,
        "vsz_kb": vsz
    });

    let response = json!({
        "status": status,
        "timestamp": SystemTime::now()
            .duration_since(SystemTime::UNIX_EPOCH)
            .unwrap()
            .as_secs(),
        "service": "rust-service",
        "checks": checks
    });

    let status_code = if status == "healthy" {
        StatusCode::OK
    } else {
        StatusCode::SERVICE_UNAVAILABLE
    };

    Ok(Json(response))
}

pub fn health_router(checker: Arc<HealthChecker>) -> Router {
    Router::new()
        .route("/health", get(health_check))
        .with_state(checker)
}
```

### Logging and Tracing

```rust
// src/logging.rs
use tracing::{Level, Subscriber};
use tracing_subscriber::{
    fmt::format::FmtSpan, layer::SubscriberExt, util::SubscriberInitExt, EnvFilter, Registry,
};

pub fn init_tracing() {
    let filter = EnvFilter::try_from_default_env()
        .unwrap_or_else(|_| EnvFilter::new("info"));

    let fmt_layer = tracing_subscriber::fmt::layer()
        .with_target(false)
        .with_span_events(FmtSpan::CLOSE)
        .json();

    Registry::default()
        .with(filter)
        .with(fmt_layer)
        .init();
}

// Middleware for request tracing
use axum::{extract::Request, middleware::Next, response::Response};
use tracing::{info_span, Instrument};

pub async fn trace_requests(request: Request, next: Next) -> Response {
    let method = request.method().clone();
    let uri = request.uri().clone();
    
    let span = info_span!(
        "http_request",
        method = %method,
        uri = %uri,
    );

    async move {
        let response = next.run(request).await;
        tracing::info!(
            status = response.status().as_u16(),
            "Request completed"
        );
        response
    }
    .instrument(span)
    .await
}
```

### CI/CD Integration

```dockerfile
# Dockerfile.ci
# Build stage
FROM rust:1.75-slim as builder

RUN apt-get update && apt-get install -y \
    pkg-config \
    libssl-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy manifests
COPY Cargo.toml Cargo.lock ./

# Build dependencies
RUN mkdir src && echo "fn main() {}" > src/main.rs
RUN cargo build --release
RUN rm -rf src

# Copy source
COPY src ./src

# Run tests
RUN cargo test --release

# Build application
RUN touch src/main.rs
RUN cargo build --release

# Security scan stage
FROM rust:1.75-slim as security

RUN cargo install cargo-audit cargo-deny

WORKDIR /app
COPY Cargo.toml Cargo.lock ./

# Security audits
RUN cargo audit
RUN cargo deny check

# Production stage
FROM debian:bookworm-slim as production

RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

RUN useradd --create-home --shell /bin/bash rustuser

COPY --from=builder /app/target/release/rust-docker-demo /usr/local/bin/rust-docker-demo
RUN chmod +x /usr/local/bin/rust-docker-demo

USER rustuser

CMD ["rust-docker-demo"]
```

## Best Practices

### 1. Multi-stage Builds
Always use multi-stage builds to minimize final image size and improve security.

### 2. Dependency Caching
Cache Cargo dependencies in separate layer for faster builds:

```dockerfile
# Copy manifests first
COPY Cargo.toml Cargo.lock ./

# Build dependencies (cached layer)
RUN mkdir src && echo "fn main() {}" > src/main.rs
RUN cargo build --release
RUN rm -rf src

# Then copy source and build
COPY src ./src
RUN touch src/main.rs
RUN cargo build --release
```

### 3. Security
- Use non-root users
- Use distroless or minimal base images
- Scan for vulnerabilities with `cargo audit`
- Keep dependencies updated

### 4. Resource Optimization
- Use release builds for production
- Strip symbols for smaller binaries
- Consider using `cargo-chef` for better caching

### 5. Configuration Management
Use environment variables for configuration, never hardcode values.

### 6. Health Checks
Implement comprehensive health checks including:
- Database connectivity
- External service availability
- Memory usage
- System resources

### 7. Logging
Use structured logging with appropriate log levels and correlation IDs.

### 8. Error Handling
Implement proper error handling with meaningful error messages and appropriate HTTP status codes.

### 9. Testing
Include unit tests, integration tests, and database tests in your Docker build process.

### 10. Monitoring
Implement metrics collection and distributed tracing for production applications.

This guide provides a comprehensive overview of dockerizing Rust applications with various frameworks and architectural patterns. The examples demonstrate production-ready patterns that can be adapted to your specific needs.