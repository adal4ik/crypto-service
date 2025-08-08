# 🪙 Crypto Service

A microservice written in Go that collects, stores, and serves cryptocurrency price data in real-time.  
The service tracks selected cryptocurrencies, periodically fetches their USD prices from the [CoinGecko API](https://www.coingecko.com/en/api), and stores the data in PostgreSQL.

---

## 📦 API Endpoints

### `POST /currency/add`

Adds a cryptocurrency to the tracking list.

**Request body:**
```json
{
  "symbol": "BTC"
}
```

**Response:**
```json
{
  "code": 201,
  "status": "success",
  "data": "Currency added to tracking list"
}
```

---

### `POST /currency/remove`

Removes a cryptocurrency from the tracking list.

**Request body:**
```json
{
  "symbol": "BTC"
}
```

**Response:**
```json
{
  "code": 204
}
```

---

### `GET /currency/price?coin=BTC&timestamp=1736500490`

Returns the price of the specified coin at the given UNIX timestamp.  
If no exact match is found, the closest available price is returned.

**Response:**
```json
{
  "code": 200,
  "status": "success",
  "data": {
    "symbol": "BTC",
    "price": 29943.12,
    "timestamp": 1736500485
  }
}
```

---

## ⚙️ Environment Configuration

Create a `.env` file in the project root based on `.env.example`:

```env
# PostgreSQL
DB_USER=postgres
DB_PASSWORD=supersecret
DB_HOST=localhost
DB_NAME=crypto
DB_PORT=5432

# App
APP_PORT=8080

# Price Collector
COLLECTOR_INTERVAL_SECONDS=60
COINGECKO_API_URL=https://api.coingecko.com/api/v3/simple/price
```

---

## 🚀 Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/adal4ik/crypto-service.git
cd crypto-service
```

### 2. Create `.env` file

```bash
cp .env.example .env
# Then edit .env as needed
```

### 3. Run the service with Docker

```bash
make up
```

This will:

- Build the Go application
- Start the app and PostgreSQL
- Run database migrations
- Start the price collector

The service will be available at:  
**`http://localhost:${APP_PORT}`**

---

## 🗄 Database Migrations

Migrations are located in the `./migrations` folder.

You can run them manually with:

```bash
make migrate-up
```

To revert the last migration:

```bash
make migrate-down
```

> Requires [golang-migrate](https://github.com/golang-migrate/migrate) installed locally.

---

## 🧱 Project Structure

```
.
├── cmd/                # Entry point (main.go)
├── internal/           # Application logic
│   ├── config/         # Configuration loading
│   ├── domain/         # Domain models and DTOs
│   ├── handler/        # HTTP handlers and routes
│   ├── repository/     # Database interaction
│   └── service/        # Business logic
├── migrations/         # SQL migration files
├── pkg/                # Shared helpers (logger, errors, response)
├── Dockerfile
├── docker-compose.yml
└── Makefile
```

---

## 🧑‍💻 Author

GitHub: [adal4ik](https://github.com/adal4ik)
