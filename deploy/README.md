# ZeePass Docker Compose Deployment

This directory contains the production-ready Docker Compose setup for ZeePass with comprehensive monitoring.

## 🏗️ Architecture

### Core Services
- **ZeePass** - Main application (Port 8080)
- **Redis** - Data storage and caching (Port 6379)

### Monitoring Stack
- **Prometheus** - Metrics collection (Port 9090)
- **Grafana** - Dashboards and visualization (Port 3000)
- **Loki** - Log aggregation (Port 3100)
- **Promtail** - Log collector
- **Alertmanager** - Alert management (Port 9093)

### Metrics Exporters
- **Node Exporter** - System metrics (Port 9100)
- **Redis Exporter** - Redis performance (Port 9121)
- **cAdvisor** - Container metrics (Port 8081)

## 🚀 Deployment

### Prerequisites
- Docker Engine 20.10+
- Docker Compose 2.0+
- Minimum 4GB RAM, 2 CPU cores
- 20GB+ available storage

### Quick Start

1. **Clone and navigate**
   ```bash
   git clone https://github.com/anazri/zeepass.git
   cd zeepass/deploy
   ```

2. **Configure secrets (Important!)**
   ```bash
   # Edit redis.conf - uncomment and set Redis password
   nano redis.conf
   
   # Change Grafana admin password in docker-compose.yml
   nano docker-compose.yml
   ```

3. **Deploy the stack**
   ```bash
   docker-compose up -d
   ```

4. **Verify deployment**
   ```bash
   docker-compose ps
   ```

## 📊 Access Points

| Service | URL | Default Credentials |
|---------|-----|-------------------|
| **ZeePass** | http://localhost:8080 | - |
| **Grafana** | http://localhost:3000 | admin / zeepass2024! |
| **Prometheus** | http://localhost:9090 | - |
| **Alertmanager** | http://localhost:9093 | - |

## 🔧 Configuration

### Redis Security
Edit `redis.conf`:
```conf
# Uncomment and set strong password
requirepass your-super-secure-redis-password-here
```

Update ZeePass environment in `docker-compose.yml`:
```yaml
environment:
  - REDIS_URL=redis://:your-password@redis:6379
```

### Alert Configuration
Edit `monitoring/alertmanager.yml` for:
- Email notifications (SMTP settings)
- Slack webhooks
- Custom alert routing

### Grafana Dashboards
- Pre-configured Prometheus and Loki datasources
- Import additional dashboards from [Grafana Labs](https://grafana.com/grafana/dashboards/)

## 📈 Monitoring Metrics

### Application Metrics
- HTTP request rates and latency
- Error rates and status codes
- Encryption operation performance
- Active connections and sessions

### Infrastructure Metrics
- CPU, Memory, Disk usage
- Network I/O and bandwidth
- Docker container performance
- Redis performance and memory usage

### Security Metrics
- Failed authentication attempts
- Suspicious activity patterns
- Encryption failures
- Access pattern anomalies

## 🚨 Alerting Rules

Create `monitoring/alert_rules.yml` for custom alerts:
```yaml
groups:
  - name: zeepass.rules
    rules:
      - alert: ZeePassDown
        expr: up{job="zeepass"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "ZeePass application is down"
```

## 🔒 Production Hardening

1. **Change default passwords**
2. **Enable Redis authentication**
3. **Configure TLS/SSL certificates**
4. **Set up proper firewall rules**
5. **Enable log rotation**
6. **Configure backup strategies**

## 📝 Maintenance

### Updates
```bash
# Pull latest images
docker-compose pull

# Restart with new images
docker-compose up -d
```

### Backup
```bash
# Backup Redis data
docker-compose exec redis redis-cli BGSAVE

# Backup Grafana dashboards
docker-compose exec grafana tar -czf - /var/lib/grafana > grafana-backup.tar.gz
```

### Logs
```bash
# View all logs
docker-compose logs -f

# View specific service
docker-compose logs -f zeepass
```

## 🆘 Troubleshooting

### Common Issues

**Port conflicts:**
```bash
# Check port usage
netstat -tlnp | grep :8080
```

**Memory issues:**
```bash
# Check container resource usage
docker stats
```

**Redis connection issues:**
```bash
# Test Redis connectivity
docker-compose exec redis redis-cli ping
```

### Health Checks
All services include health checks. Check status:
```bash
docker-compose ps
```

## 🌐 Cloud Deployment

### Recommended Cloud Specs
- **CPU**: 2-4 vCPUs
- **RAM**: 4-8 GB
- **Storage**: 20-50 GB SSD
- **Network**: 1+ Gbps

### Cloud Provider Examples
- **AWS**: t3.large + EBS volumes
- **DigitalOcean**: $48-96/month droplets
- **Google Cloud**: e2-standard-2/4
- **Azure**: Standard_B2s/B4s