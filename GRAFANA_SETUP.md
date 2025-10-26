# GoLLuM Grafana Dashboard Setup Guide

This guide will help you set up monitoring for GoLLuM using Grafana.

## Prerequisites

- GoLLuM server running on `http://localhost:8080`
- Grafana installed and running (typically on `http://localhost:3000`)

## Quick Setup (Direct Metrics)

### Step 1: Start GoLLuM

```bash
cd /Users/home/Desktop/GoLLuM
bin/altiserve
```

The server will start on `http://localhost:8080` with metrics at `/metrics`.

### Step 2: Add Prometheus Data Source in Grafana

1. Open Grafana: `http://localhost:3000`
2. Log in (default: admin/admin)
3. Go to **Configuration** → **Data Sources**
4. Click **Add data source**
5. Select **Prometheus**
6. Set **URL** to: `http://localhost:8080`
7. Set **Access** to: **Server (default)**
8. Click **Save & Test**

**Note:** Grafana will scrape GoLLuM's `/metrics` endpoint directly.

### Step 3: Import the Dashboard

1. Go to **Dashboards** → **Import**
2. Click **Upload JSON file**
3. Select `/Users/home/Desktop/GoLLuM/grafana-dashboard.json`
4. Click **Load**
5. Verify the **Prometheus** data source is selected
6. Click **Import**

### Step 4: View Your Dashboard

The dashboard will show:
- **Requests / 5m**: Decode steps per second
- **Cache Hit Rates**: Prefix and prompt cache performance
- **TTFT/TPOT**: Time to first/last token percentiles
- **Batch Sizes**: Distribution of decode batch sizes
- **KV Events**: Allocation, pinning, unpinning, evictions
- **Cache Events**: Detailed hit/miss breakdown

## Advanced Setup (with Prometheus Server)

If you prefer using a dedicated Prometheus server:

### Step 1: Install Prometheus

```bash
brew install prometheus  # macOS
# or download from https://prometheus.io/download/
```

### Step 2: Configure Prometheus

Create a `prometheus.yml` file:

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'gollum'
    static_configs:
      - targets: ['localhost:8080']
```

### Step 3: Start Prometheus

```bash
prometheus --config.file=prometheus.yml
```

### Step 4: Add Prometheus as Data Source

1. In Grafana, add **Prometheus** data source
2. Set **URL** to: `http://localhost:9090`
3. Save and import the dashboard as above

## Dashboard Panels

The GoLLuM dashboard includes:

1. **Request Rate**: Decode steps per second (5-minute window)
2. **Prefix Cache Hit Rate**: Longest-prefix match effectiveness
3. **Prompt Cache Hit Rate**: Exact prompt match effectiveness
4. **TTFT Percentiles**: Time to first token (p50/p95)
5. **TPOT Percentiles**: Time to produce all tokens (p50/p95)
6. **Batch Size**: Distribution across decode batches
7. **KV Pager Events**: Memory management operations
8. **Cache Events**: Comprehensive hit/miss analysis

## PromQL Query Examples

### Cache Hit Rate by Model
```promql
sum(rate(gollum_cache_events_total{cache_type="prefix",hit_kind="hit"}[5m])) by (model)
/
sum(rate(gollum_cache_events_total{cache_type="prefix"}[5m])) by (model)
```

### TTFT p95 Across Models
```promql
histogram_quantile(0.95, sum(rate(gollum_ttft_ms_bucket[5m])) by (le, model))
```

### KV Evictions per Second
```promql
sum(rate(gollum_kv_events_total{action="evict"}[5m])) by (model)
```

## Testing the Setup

1. Start GoLLuM: `bin/altiserve`
2. Send some requests to warm up the caches:
```bash
curl -N -H "Content-Type: application/json" \
  -d '{"model":"toy-1","messages":[{"role":"user","content":"test"}],"max_tokens":10}' \
  http://localhost:8080/v1/chat/completions
```
3. Check metrics: `curl http://localhost:8080/metrics | grep gollum_`
4. Open Grafana and view the dashboard

## Troubleshooting

- **No data in Grafana?** Check that GoLLuM is running and the metrics endpoint is accessible
- **"Data source not found"?** Verify the Prometheus data source name matches "Prometheus"
- **Metrics not updating?** Ensure requests are being made to GoLLuM to generate metrics

## Configuration

The dashboard JSON file (`grafana-dashboard.json`) can be customized:
- Adjust time ranges (currently 6 hours)
- Modify panel positions and sizes
- Add additional panels for custom metrics
- Change datasource name if needed

