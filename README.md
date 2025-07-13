 🌐 URL Crawler Dashboard

A full-stack web application to crawl and analyze website URLs. It extracts structured data like heading tags, link stats, broken links, and login form presence with a polished, user-friendly dashboard.

---

## ⚙️ Features

- ✅ Add and queue website URLs for crawling
- 🔍 Extract and display:
  - HTML version
  - Page title
  - Counts of headings (H1–H6)
  - Internal & external link counts
  - Broken link count
  - Login form detection
- ♻️ Re-analyze any URL
- 🗑️ Delete any queued or completed URL
- 📊 Pie chart visualization of link types
- 💡 Responsive UI with Tailwind CSS

---

## 🧱 Tech Stack

- **Frontend:** React, TypeScript, React Query, React Router, Tailwind CSS, Recharts
- **Backend:** Go (Gin), GORM, MySQL
- **API:** REST
- **Other:** Axios, dotenv, HTML parsing

---

## 📁 Folder Structure

.
├── crawler-api/ # Go backend
│ ├── internal/
│ │ ├── crawler/ # Core crawling logic
│ │ ├── db/ # DB connection setup
│ │ └── models/ # GORM models
│ └── main.go
│
├── crawler-frontend/ # React frontend
│ ├── src/
│ │ ├── components/ # Layout, AddUrlForm, UrlTable
│ │ ├── pages/ # Dashboard, UrlDetails
│ │ └── types/ # TypeScript types
│ └── App.tsx


---

## 🚀 Getting Started

### 🛠️ Backend Setup

1. Start MySQL and create a database:

```sql
CREATE DATABASE crawler_db;
Inside crawler-api/, create a .env file:
DB_USER=root
DB_PASSWORD=yourpassword
DB_NAME=crawler_db
DB_HOST=localhost
PORT=8080
Run the Go server:
cd crawler-api
go run main.go
💻 Frontend Setup
Install and run:
cd crawler-frontend
npm install
npm run dev
Open the app in your browser at:
http://localhost:3000
🧪 How to Use

Enter a URL (e.g., https://example.com) in the form.
Wait a few seconds for analysis to complete.
View results on the dashboard.
Click "View" for details with pie chart.
Use "Re-analyze" or "Delete" for any listed URL.
