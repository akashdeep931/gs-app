# GS App — Shipment Packing Optimiser

A full-stack application that calculates optimal pack combinations for shipping orders. Given a quantity of items, the system determines the fewest packs needed to fulfil (or exceed) the order using configurable pack sizes.

Admins can manage the available pack sizes; users can calculate shipment breakdowns.

**Live app:** https://gs-app-web.fly.dev/

## Tech Stack

| Layer    | Technology                                      |
| -------- | ----------------------------------------------- |
| Backend  | Go, chi router, in-memory store                 |
| Frontend | React 19, TypeScript, Vite, Tailwind CSS, MUI   |
| Deploy   | Fly.io (Docker — Alpine / Nginx)                |

## Getting Started

### Prerequisites

- Go 1.25+
- Node.js 22+
- npm

### Backend

```bash
cd backend
cp .env.example .env      # PORT=8000
go run main.go
```

### Frontend

```bash
cd frontend
cp .env.example .env      # VITE_GS_API_URL=http://localhost:8000/v1
npm install
npm run dev
```

The frontend runs on `http://localhost:5173` and the API on `http://localhost:8000`.

### Environment Variables

| Variable           | Location | Default                        | Description              |
| ------------------ | -------- | ------------------------------ | ------------------------ |
| `PORT`             | Backend  | `8000`                         | Server listening port    |
| `VITE_GS_API_URL`  | Frontend | `http://localhost:8000/v1`     | Backend API base URL     |

### Running Tests

```bash
cd backend
go test ./...
```

---

## API Reference

Version: `v1`

All request and response bodies are JSON. The API has no authentication.

---

### Health Check

```
GET /v1/health
```

**Response** `200`

```json
{
  "status": "ok"
}
```

---

### Get Packs

Retrieve all configured pack sizes (sorted ascending).

```
GET /v1/packs
```

**Response** `200`

```json
{
  "packs": [250, 500, 1000, 2000, 5000]
}
```

---

### Add Pack

Add a new pack size.

```
POST /v1/packs
```

**Request Body**

| Field  | Type | Required | Rules            |
| ------ | ---- | -------- | ---------------- |
| `size` | int  | Yes      | Positive integer |

```json
{
  "size": 750
}
```

**Response** `201`

```json
{
  "packs": [250, 500, 750, 1000, 2000, 5000]
}
```

**Errors**

| Status | Reason                        |
| ------ | ----------------------------- |
| `400`  | Invalid JSON or non-positive  |
| `409`  | Pack size already exists      |

---

### Set Packs

Replace the entire pack configuration.

```
PUT /v1/packs
```

**Request Body**

| Field   | Type  | Required | Rules                           |
| ------- | ----- | -------- | ------------------------------- |
| `sizes` | int[] | Yes      | Non-empty, all positive integers |

```json
{
  "sizes": [300, 600, 1200]
}
```

**Response** `200`

```json
{
  "packs": [300, 600, 1200]
}
```

**Errors**

| Status | Reason                                    |
| ------ | ----------------------------------------- |
| `400`  | Invalid JSON, empty array, or non-positive |

---

### Delete Pack

Remove a specific pack size.

```
DELETE /v1/packs/{size}
```

**URL Parameters**

| Param  | Type | Description          |
| ------ | ---- | -------------------- |
| `size` | int  | Pack size to remove  |

**Response** `200`

```json
{
  "packs": [250, 500, 2000, 5000]
}
```

**Errors**

| Status | Reason                                 |
| ------ | -------------------------------------- |
| `400`  | Invalid size or cannot remove last pack |
| `404`  | Pack size not found                    |

---

### Calculate Shipment

Calculate the optimal pack combination for an order quantity. Uses a dynamic programming algorithm to minimise the number of packs whilst ensuring the order is fully covered.

```
POST /v1/shipment
```

**Request Body**

| Field   | Type | Required | Rules            |
| ------- | ---- | -------- | ---------------- |
| `items` | int  | Yes      | Positive integer |

```json
{
  "items": 1500
}
```

**Response** `200`

```json
{
  "packs": {
    "1000": 1,
    "500": 1,
    "250": 0,
    "2000": 0,
    "5000": 0
  },
  "items_shipped": 1500
}
```

`items_shipped` may be greater than the requested quantity when an exact match is not possible.

**Errors**

| Status | Reason                       |
| ------ | ---------------------------- |
| `400`  | Invalid JSON or non-positive |

---

## Project Structure

```
├── backend/
│   ├── main.go                          # Entry point & router
│   └── internal/
│       ├── handler/                     # HTTP handlers
│       ├── model/                       # Request structs
│       ├── store/                       # In-memory pack store
│       ├── middleware/                   # Logging & request IDs
│       └── packSizesCalculator/         # DP packing algorithm
└── frontend/
    └── src/
        ├── pages/                       # Home, User, Admin
        ├── services/                    # API client (singleton)
        ├── context/                     # React context for service
        ├── hooks/                       # useGSService hook
        └── types/                       # TypeScript interfaces
```
