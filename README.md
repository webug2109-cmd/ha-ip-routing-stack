# High-Availability IP Routing Stack

**Context:** This project was built as a proof-of-concept to demonstrate SRE/Infrastructure best practices, specifically tailored to handle high-throughput, latency-sensitive APIs (like IP-lookup services).

## Architecture Highlights
*   **Infrastructure (Terraform):** Ready-to-deploy Google Kubernetes Engine (GKE) cluster and IAM role definitions.
*   **Configuration (Ansible):** Playbooks for bare-metal Nginx load balancer provisioning, demonstrating comfort in heterogeneous/hybrid environments.
*   **Orchestration (Kubernetes):** High-availability deployments with strict resource boundaries, readiness/liveness probes, and rolling update strategies.
*   **Observability (Prometheus/Grafana):** Fully instrumented Go application scraping metrics into a `kube-prometheus-stack`.

## Local Testing via Minikube
While Terraform is built for GCP production, you can test the stack locally:
```bash
minikube start
eval $(minikube docker-env)
cd app && docker build -t ipinfo-mock:latest .
cd ../k8s && kubectl apply -f .
```

## Incident Response / Runbooks

**Scenario:** High 99th percentile latency on the `/lookup` endpoint.
**Action:** 
1. Check the Grafana dashboard: Identify if the latency correlates with Node CPU saturation or a specific geographic region.
2. Exec into the pod: Use `tcpdump` or `mtr` to verify if the network layer (CNI) is dropping packets.
3. Check Kubernetes events: Look for `OOMKilled` pods or evictions indicating resource exhaustion.
