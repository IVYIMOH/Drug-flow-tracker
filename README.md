# AfyaTrack 🏥

> Real-time drug supply transparency for Kenyan public hospitals — built to catch when private suppliers are used more than the government's KEMSA pipeline.

---

## The Problem

Kenya's public hospitals face a persistent crisis: essential medicines run out, yet the same drugs are available — at inflated prices — from private suppliers. This isn't a supply problem, it's a **visibility problem**.

AfyaTrack answers one question: *"Is this hospital buying drugs it should be getting for free?"*

---

## What It Does

| Feature | Detail |
|---|---|
| 📦 Stock logging | Tags every entry as `KEMSA` (public) or `PRIVATE` |
| 💊 Dispensation tracking | Deducts from stock so available quantities stay accurate |
| 📊 Private supplier ratio | Fires an alert when private supply exceeds 40% of total |
| 💰 Price benchmarking | Flags any drug procured above 1.5× the KEMSA reference rate |

---

## Project Structure

```
Drug-flow-tracker/
├── main.go                 # Entry point — starts server on :8080
├── dev.ps1                 # Windows dev runner (sets env vars + runs server)
├── frontend/
│   └── index.html          # Browser dashboard — open directly in Chrome/Edge
├── db/
│   └── db.go               # PostgreSQL connection (reads from env vars)
├── models/
│   └── models.go           # StockEntry, Dispensation, SupplierSource types
├── handlers/
│   ├── stock.go            # POST /stock, GET /stock
│   ├── dispense.go         # POST /dispense
│   └── insights.go         # GET /insights
├── routes/
│   └── routes.go           # Route registration
└── utils/
    └── price.go            # KEMSA price benchmarking logic
```

---

## Prerequisites

| Tool | Version | Download |
|---|---|---|
| Go | 1.21+ | https://golang.org/dl |
| PostgreSQL | 17 | https://www.postgresql.org/download |
| Git | Any | https://git-scm.com |

---

## Setup — Step by Step

### 1. Clone the repo

```bash
git clone https://github.com/IVYIMOH/Drug-flow-tracker.git
cd Drug-flow-tracker
```

### 2. Install Go dependencies

```bash
go mod tidy
```

### 3. Create the database

Open a terminal and connect to PostgreSQL:

```bash
# Linux / macOS
psql -U postgres

# Windows (PowerShell)
& "C:\Program Files\PostgreSQL\17\bin\psql.exe" -U postgres
```

Then run these SQL commands:

```sql
CREATE DATABASE drugflow;
\c drugflow

CREATE TABLE stock_entries (
    id          SERIAL PRIMARY KEY,
    hospital_id INT NOT NULL,
    drug_name   TEXT NOT NULL,
    source      TEXT NOT NULL CHECK (source IN ('KEMSA', 'PRIVATE')),
    quantity    INT NOT NULL,
    unit_price  NUMERIC(10,2) NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE dispensations (
    id           SERIAL PRIMARY KEY,
    hospital_id  INT NOT NULL,
    drug_name    TEXT NOT NULL,
    quantity     INT NOT NULL,
    dispensed_at TIMESTAMPTZ DEFAULT NOW()
);

-- Verify tables were created
\dt
```

Expected output:
```
          List of relations
 Schema |     Name      | Type  |  Owner
--------+---------------+-------+----------
 public | dispensations | table | postgres
 public | stock_entries | table | postgres
```

Then exit:
```sql
\q
```

### 4. Set environment variables

**Windows (PowerShell) — recommended: use the dev script**

Create `dev.ps1` in the project root:

```powershell
$env:DB_USER     = "postgres"
$env:DB_PASSWORD = "your_postgres_password"
$env:DB_NAME     = "drugflow"
go run main.go
```

> ⚠️ Replace `your_postgres_password` with the password you set when installing PostgreSQL.
> `dev.ps1` is in `.gitignore` — your password will never be committed.

Allow PowerShell to run scripts (one-time setup):

```powershell
Set-ExecutionPolicy -Scope CurrentUser RemoteSigned
```

**Linux / macOS**

```bash
export DB_USER=postgres
export DB_PASSWORD=your_postgres_password
export DB_NAME=drugflow
go run main.go
```

### 5. Run the server

**Windows:**
```powershell
.\dev.ps1
```

**Linux / macOS:**
```bash
go run main.go
```

✅ You should see:
```
2026/05/01 22:24:05 DB connection successful
[GIN-debug] Listening and serving HTTP on :8080
```

### 6. Open the dashboard

Open `frontend/index.html` directly in Chrome or Edge — no extra server needed.

The API status pill in the top-right corner will turn **green** when connected.

---

## API Reference

Base URL: `http://localhost:8080`

### POST `/stock` — Log a stock entry

**Request body:**
```json
{
  "hospital_id": 1,
  "drug_name": "Amoxicillin",
  "source": "KEMSA",
  "quantity": 100,
  "unit_price": 9.00
}
```

- `source` must be exactly `"KEMSA"` or `"PRIVATE"` (case-sensitive)
- `unit_price` is in Kenyan Shillings (KSh)

**Response:**
```json
{
  "message": "Stock entry added successfully",
  "high_price": false,
  "price_benchmark": true
}
```

- `high_price: true` means the price exceeds 1.5× the KEMSA reference rate
- `price_benchmark: false` means the drug isn't in the KEMSA reference list yet

---

### POST `/dispense` — Record a dispensation

**Request body:**
```json
{
  "hospital_id": 1,
  "drug_name": "Amoxicillin",
  "quantity": 20
}
```

**Response:**
```json
{
  "message": "Drug dispensed successfully"
}
```

**Error (insufficient stock):**
```json
{
  "error": "insufficient stock",
  "requested": 20,
  "available_stock": 5
}
```

---

### GET `/stock` — Available stock per drug

**Response:**
```json
[
  { "drug_name": "Amoxicillin", "available_stock": 160 }
]
```

`available_stock` = total stocked − total dispensed

---

### GET `/insights` — Private supplier ratio + alert

**Response:**
```json
{
  "alert": "⚠️ High reliance on private suppliers",
  "private_ratio": 0.4444,
  "private_units": 80,
  "total_units": 180
}
```

- `alert` is an empty string `""` when the ratio is below 40%
- `private_ratio` is a decimal (0.4444 = 44.44%)

---

## Testing with curl

### Windows (PowerShell)

Create a JSON file to avoid escaping issues:

```powershell
# Save as UTF-8 without BOM
[System.IO.File]::WriteAllText(
  "$PWD\stock.json",
  '{"hospital_id":1,"drug_name":"Amoxicillin","source":"KEMSA","quantity":100,"unit_price":9.0}',
  [System.Text.UTF8Encoding]::new($false)
)

curl.exe -X POST http://localhost:8080/stock -H "Content-Type: application/json" -d "@stock.json"
```

### Linux / macOS

```bash
# Add KEMSA stock
curl -X POST http://localhost:8080/stock \
  -H "Content-Type: application/json" \
  -d '{"hospital_id":1,"drug_name":"Amoxicillin","source":"KEMSA","quantity":100,"unit_price":9.0}'

# Add overpriced private stock
curl -X POST http://localhost:8080/stock \
  -H "Content-Type: application/json" \
  -d '{"hospital_id":1,"drug_name":"Amoxicillin","source":"PRIVATE","quantity":80,"unit_price":18.0}'

# Dispense
curl -X POST http://localhost:8080/dispense \
  -H "Content-Type: application/json" \
  -d '{"hospital_id":1,"drug_name":"Amoxicillin","quantity":20}'

# Check stock
curl http://localhost:8080/stock

# Check insights — should fire alert at 44% private
curl http://localhost:8080/insights
```

---

## How the Alert Works

The private supplier ratio is calculated as:

$$\text{Private Ratio} = \frac{\sum \text{quantity}_{PRIVATE}}{\sum \text{quantity}_{total}}$$

When this exceeds **40%**, the `/insights` endpoint returns an alert string and the dashboard banner activates.

Price benchmarking flags any drug procured above **1.5× the KEMSA reference rate**:

$$\text{High Price} = P_{supplied} > 1.5 \times P_{KEMSA}$$

Current KEMSA reference prices (hardcoded — moving to DB in next version):

| Drug | KEMSA rate (KSh) | Flag threshold |
|---|---|---|
| Amoxicillin | 10 | > 15 |
| Paracetamol | 5 | > 7.50 |
| Metformin | 8 | > 12 |
| Ciprofloxacin | 15 | > 22.50 |
| Quinine | 12 | > 18 |
| Cotrimoxazole | 7 | > 10.50 |

---

## Troubleshooting

### `DB unreachable: connection refused`
PostgreSQL isn't running. Start it:
```powershell
# Windows — find your service name first
Get-Service -Name postgresql*
Start-Service -Name postgresql-x64-17
```

### `password authentication failed`
Wrong password in `dev.ps1`. Reset it in psql:
```sql
ALTER USER postgres WITH PASSWORD 'newpassword';
```

### `SSL is not enabled on the server`
Your `db.go` connection string is missing `sslmode=disable`. It should read:
```go
connStr := fmt.Sprintf(
    "user=%s password=%s dbname=%s sslmode=disable",
    os.Getenv("DB_USER"),
    os.Getenv("DB_PASSWORD"),
    os.Getenv("DB_NAME"),
)
```

### Dashboard shows "API offline"
The Go server isn't running. Open a terminal and run `.\dev.ps1`.

### `invalid character 'ï'` error on Windows
Your JSON file was saved with a BOM marker. Recreate it with:
```powershell
[System.IO.File]::WriteAllText(
  "$PWD\stock.json",
  '{"hospital_id":1,...}',
  [System.Text.UTF8Encoding]::new($false)
)
```

### CORS error in browser console
Add this middleware to `main.go` before `routes.SetupRoutes(r)`:
```go
r.Use(func(c *gin.Context) {
    c.Header("Access-Control-Allow-Origin", "*")
    c.Header("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
    c.Header("Access-Control-Allow-Headers", "Content-Type")
    if c.Request.Method == "OPTIONS" {
        c.AbortWithStatus(204)
        return
    }
    c.Next()
})
```

---

## Roadmap

- [ ] Per-hospital insights (`/insights?hospital_id=1`)
- [ ] Drug-level private ratio breakdown
- [ ] Time-series alert history and audit log
- [ ] KEMSA reference prices in database (not hardcoded)
- [ ] SMS notifications via Africa's Talking API
- [ ] Authentication — hospitals can only submit their own data

---

## Built With

- [Go](https://golang.org/) — API server
- [Gin](https://github.com/gin-gonic/gin) — HTTP framework
- [PostgreSQL](https://www.postgresql.org/) — Database
- [lib/pq](https://github.com/lib/pq) — PostgreSQL driver

---

## License

MIT — see [LICENSE](LICENSE)

---

*Built during GNEC Hackathon 2026. Every flagged overpriced purchase is a data point toward accountability.*
