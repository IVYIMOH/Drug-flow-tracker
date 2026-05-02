AfyaTrack 🏥

Real-time drug supply tracker for Kenyan public hospitals — built to flag over-reliance on private suppliers over the government's KEMSA pipeline.


The Problem
Kenya's public hospitals face a persistent crisis: essential medicines run out, yet the same drugs are available — at inflated prices — from private suppliers. This isn't a supply problem, it's a visibility problem.
AfyaTrack answers one question: "Is this hospital buying drugs it should be getting for free?"

What It Does

📦 Logs stock entries tagged as KEMSA (public) or PRIVATE
💊 Tracks dispensations so available stock is always accurate
📊 Computes a private supplier ratio and fires an alert when over 40% of supply comes from private sources
💰 Benchmarks unit prices against KEMSA reference rates — flags any drug procured at more than 1.5× the government price


API Endpoints
MethodEndpointDescriptionPOST/stockLog a new stock entry (KEMSA or PRIVATE)POST/dispenseRecord a drug dispensationGET/stockGet real-time available stock per drugGET/insightsGet private supplier ratio + alert status

Quick Start
Prerequisites

Go 1.21+
PostgreSQL 17

1. Clone the repo
bashgit clone https://github.com/YOUR_USERNAME/Drug-flow-tracker.git
cd Drug-flow-tracker
2. Set up the database
bashpsql -U postgres
sqlCREATE DATABASE drugflow;
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
3. Set environment variables
Linux/macOS:
bashexport DB_USER=postgres
export DB_PASSWORD=yourpassword
export DB_NAME=drugflow
Windows (PowerShell):
powershell$env:DB_USER     = "postgres"
$env:DB_PASSWORD = "yourpassword"
$env:DB_NAME     = "drugflow"
4. Run the server
bashgo run main.go
You should see:
DB connection successful
[GIN-debug] Listening and serving HTTP on :8080

Usage Examples
Add KEMSA stock
bashcurl -X POST http://localhost:8080/stock \
  -H "Content-Type: application/json" \
  -d '{"hospital_id":1,"drug_name":"Amoxicillin","source":"KEMSA","quantity":100,"unit_price":9.0}'
json{"message":"Stock entry added successfully","high_price":false,"price_benchmark":true}
Add private supplier stock (overpriced)
bashcurl -X POST http://localhost:8080/stock \
  -H "Content-Type: application/json" \
  -d '{"hospital_id":1,"drug_name":"Amoxicillin","source":"PRIVATE","quantity":80,"unit_price":18.0}'
Dispense a drug
bashcurl -X POST http://localhost:8080/dispense \
  -H "Content-Type: application/json" \
  -d '{"hospital_id":1,"drug_name":"Amoxicillin","quantity":20}'
Check available stock
bashcurl http://localhost:8080/stock
json[{"available_stock":160,"drug_name":"Amoxicillin"}]
Check private supplier insights
bashcurl http://localhost:8080/insights
json{
  "alert": "⚠️ High reliance on private suppliers",
  "private_ratio": 0.4444,
  "private_units": 80,
  "total_units": 180
}

How the Alert Works
The private supplier ratio is calculated as:
Private Ratio=∑quantityPRIVATE∑quantitytotal\text{Private Ratio} = \frac{\sum \text{quantity}_{PRIVATE}}{\sum \text{quantity}_{total}}Private Ratio=∑quantitytotal​∑quantityPRIVATE​​
When this exceeds 40%, the system fires an alert. A facility routinely above this threshold is a red flag for either a KEMSA delivery failure or deliberate procurement diversion.
Price benchmarking flags drugs procured above 1.5× the KEMSA reference rate:
High Price=Psupplied>1.5×PKEMSA\text{High Price} = P_{supplied} > 1.5 \times P_{KEMSA}High Price=Psupplied​>1.5×PKEMSA​

Project Structure
Drug-flow-tracker/
├── main.go               # Entry point
├── db/
│   └── db.go             # PostgreSQL connection
├── models/
│   └── models.go         # StockEntry, Dispensation, SupplierSource
├── handlers/
│   ├── stock.go          # AddStock, GetStock
│   ├── dispense.go       # DispenseDrug
│   └── insights.go       # GetInsights
├── routes/
│   └── routes.go         # Route registration
└── utils/
    └── price.go          # KEMSA price benchmarking

Roadmap

 Per-hospital insights (/insights?hospital_id=1)
 Drug-level private ratio breakdown
 Time-series alert history and audit log
 KEMSA reference prices stored in database (not hardcoded)
 SMS notifications via Africa's Talking API
 Web dashboard for hospital staff


Built With

Go — API server
Gin — HTTP framework
PostgreSQL — Database
lib/pq — PostgreSQL driver


The Mission
AfyaTrack was built to bring transparency to public health procurement in Kenya. Every flagged overpriced private purchase is a data point toward accountability. Every alert is a question that deserves an answer.
