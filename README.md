# URL Shortener

Simple URL shortening service built with Go and PostgreSQL

## Build & Run

### 1. Clone project
```bash
git clone https://github.com/mrzlkvvv/URLShortener.git && cd ./URLShortener/
```

### 2. Configure environment variables
```bash
cp ./config/example.env ./config/.env # and edit them
```

### 3. Run with Docker Compose
```bash
make up-prod # Production mode (URLShortener + PostgreSQL)
```
```bash
make up-dev # Development mode (URLShortener + PostgreSQL + pgAdmin)
```
