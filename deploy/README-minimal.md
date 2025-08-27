# ZeePass Minimal Monitoring Deployment

Lightweight Docker Compose setup with essential monitoring tools for ZeePass.

## 🏗️ Architecture

### Core Services
- **ZeePass** - Main application (via Nginx)
- **Redis** - Data storage and caching
- **Nginx** - Reverse proxy with detailed logging

### Minimal Monitoring
- **GoAccess** - Real-time web analytics (Port 7890)
- **Uptime Kuma** - Uptime monitoring (Port 3001)
- **Netdata** - System performance monitoring (Port 19999)

## 🚀 Deployment

### Prerequisites
- Docker Engine 20.10+
- Docker Compose 2.0+
- Minimum 2GB RAM, 2 CPU cores
- 10GB+ available storage

### Quick Start

1. **Navigate to deploy directory**
   ```bash
   cd zeepass/deploy
   ```

2. **Configure Redis (Important!)**
   ```bash
   # Edit redis.conf - uncomment and set Redis password
   nano redis.conf
   ```

3. **Deploy minimal stack**
   ```bash
   docker-compose -f docker-compose-minimal.yml up -d
   ```

4. **Verify deployment**
   ```bash
   docker-compose -f docker-compose-minimal.yml ps
   ```

## 📊 Access Points

| Service | URL | Purpose |
|---------|-----|---------|
| **ZeePass** | http://localhost | Main application |
| **GoAccess** | http://localhost:7890 | Real-time web analytics |
| **Uptime Kuma** | http://localhost:3001 | Uptime monitoring dashboard |
| **Netdata** | http://localhost:19999 | System performance metrics |

## 📈 What You'll Monitor

### Web Analytics (GoAccess)
- Real-time visitor statistics
- Page views and unique visitors
- Response times and status codes
- Geographic distribution
- Browser and OS statistics
- Top requested pages

### Uptime Monitoring (Uptime Kuma)
- Service availability
- Response time monitoring
- SSL certificate expiry
- Custom notifications
- Status page

### System Metrics (Netdata)
- CPU, Memory, Disk usage
- Network I/O
- Docker container performance
- Process monitoring
- Real-time dashboards

## 🔧 Configuration

### Nginx Logging
Enhanced logging format captures:
- Request details and response times
- User agents and referrers
- Upstream response times
- Rate limiting events

### Security Features
- Rate limiting on sensitive endpoints
- Security headers
- SSL/TLS configuration ready
- IP-based restrictions

### GoAccess Setup
- Real-time HTML reports
- Geographic data support
- Privacy-focused (IP anonymization)
- Static file filtering

## 🚨 Simple Alerting

### Uptime Kuma Notifications
Configure in the web interface:
- Email notifications
- Slack/Discord webhooks
- Telegram bots
- SMS alerts

### Log-based Monitoring
Monitor Nginx logs for:
```bash
# High error rates
docker-compose logs nginx | grep "HTTP/1.1\" [45]"

# Slow responses
docker-compose logs nginx | grep "rt=[5-9]\."

# Rate limit hits  
docker-compose logs nginx | grep "limiting requests"
```

## 🔒 Production Setup

### SSL/HTTPS Setup
1. **Obtain SSL certificates** (Let's Encrypt recommended)
2. **Place certificates** in `nginx/ssl/`
3. **Uncomment HTTPS server** block in `nginx.conf`
4. **Enable HTTP → HTTPS redirect**

### Redis Security
```bash
# Edit redis.conf
nano redis.conf

# Uncomment and set password
requirepass your-secure-password
```

## 📊 Resource Requirements

### Minimal Setup
- **CPU**: 2 vCPUs
- **RAM**: 2-4 GB
- **Storage**: 10-20 GB SSD

### Service Resource Usage
- ZeePass: ~200MB RAM
- Redis: ~500MB RAM  
- Nginx: ~50MB RAM
- GoAccess: ~100MB RAM
- Uptime Kuma: ~150MB RAM
- Netdata: ~200MB RAM

## 🛠️ Maintenance

### View Logs
```bash
# All services
docker-compose -f docker-compose-minimal.yml logs -f

# Specific service
docker-compose -f docker-compose-minimal.yml logs -f nginx

# Access logs
docker-compose -f docker-compose-minimal.yml exec nginx tail -f /var/log/nginx/access.log
```

### Updates
```bash
# Pull latest images
docker-compose -f docker-compose-minimal.yml pull

# Restart with updates
docker-compose -f docker-compose-minimal.yml up -d
```

### Backup
```bash
# Backup Redis data
docker-compose -f docker-compose-minimal.yml exec redis redis-cli BGSAVE

# Backup Uptime Kuma config
docker-compose -f docker-compose-minimal.yml exec uptime-kuma tar -czf - /app/data > uptime-backup.tar.gz
```

## 🆘 Traffic Analysis

### Real-time Traffic (GoAccess)
- Live visitor count
- Active pages
- Request rates
- Response time trends

### Historical Analysis
```bash
# Generate daily report
docker-compose -f docker-compose-minimal.yml exec goaccess goaccess /var/log/nginx/access.log -o /var/lib/goaccess/daily-report.html

# Search specific patterns
docker-compose -f docker-compose-minimal.yml exec nginx grep "encrypt" /var/log/nginx/access.log
```

### Key Metrics to Watch
- **Request rate**: Normal vs. spike patterns
- **Response times**: Performance degradation
- **Error rates**: 4xx/5xx status codes
- **Geographic patterns**: Unusual traffic sources
- **User agents**: Bot vs. human traffic

This minimal setup provides essential traffic visibility without the complexity of a full monitoring stack!