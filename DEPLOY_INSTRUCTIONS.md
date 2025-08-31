# ZeePass Production Deployment Instructions

This guide provides step-by-step instructions for deploying ZeePass on a server with Docker and making it accessible via HTTPS at https://zeepass.com.

## Prerequisites

- Server with Docker installed
- Domain `zeepass.com` pointing to your server's IP address
- Root or sudo access to the server

## Deployment Steps

### 1. Initialize Docker Swarm

```bash
# Initialize Docker Swarm on your server
docker swarm init

# If you have multiple network interfaces, specify the advertise address
# docker swarm init --advertise-addr YOUR_SERVER_IP
```

### 2. Obtain SSL Certificate

#### Option A: Using Let's Encrypt (Recommended)

Install Certbot:
```bash
# On Ubuntu/Debian
sudo apt update && sudo apt install certbot

# On CentOS/RHEL
sudo yum install certbot
```

Obtain certificate:
```bash
# Stop any running web server temporarily
sudo systemctl stop nginx apache2 2>/dev/null || true

# Get certificate for zeepass.com
sudo certbot certonly --standalone -d zeepass.com

# Certificate files will be at:
# /etc/letsencrypt/live/zeepass.com/fullchain.pem
# /etc/letsencrypt/live/zeepass.com/privkey.pem
```

#### Option B: Using existing SSL certificate

Place your SSL certificate files in the deploy directory:
```bash
mkdir -p deploy/nginx/ssl
# Copy your certificate files:
# - cert.pem (certificate)
# - key.pem (private key)
```

### 3. Configure SSL in Nginx

Edit `deploy/nginx/nginx.conf`:

```bash
# Uncomment and update the SSL certificate paths
ssl_certificate /etc/nginx/ssl/fullchain.pem;
ssl_certificate_key /etc/nginx/ssl/privkey.pem;

# Uncomment the HTTPS redirect in the HTTP server block (line 59)
return 301 https://$server_name$request_uri;

# Update server_name to your domain
server_name zeepass.com;
```

### 4. Update Docker Compose SSL Volume Mounts

Edit `deploy/docker-swarm-minimal.yml` to mount your SSL certificates:

```yaml
# In the nginx service volumes section, update to:
volumes:
  - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
  - /etc/letsencrypt/live/zeepass.com:/etc/nginx/ssl:ro  # For Let's Encrypt
  # OR for custom certificates:
  # - ./nginx/ssl:/etc/nginx/ssl:ro
  - nginx-logs:/var/log/nginx
  - ./nginx/html:/usr/share/nginx/html:ro
```

### 5. Configure Redis Security (Optional but Recommended)

Edit `deploy/redis.conf`:
```bash
# Uncomment and set a strong password
requirepass YOUR_STRONG_REDIS_PASSWORD
```

If you set a Redis password, update the Docker Compose environment:
```yaml
# In the zeepass service
environment:
  - REDIS_URL=redis://:YOUR_STRONG_REDIS_PASSWORD@redis:6379
```

### 6. Deploy the Stack

```bash
# Navigate to the deploy directory
cd deploy

# Deploy the stack
docker stack deploy -c docker-swarm-minimal.yml zeepass
```

### 7. Verify Deployment

Check service status:
```bash
# List all services
docker service ls

# Check service logs
docker service logs zeepass_zeepass
docker service logs zeepass_nginx
docker service logs zeepass_redis

# Check specific service status
docker service ps zeepass_zeepass
```

### 8. Access Your Services

- **ZeePass Application**: https://zeepass.com
- **GoAccess Analytics**: http://YOUR_SERVER_IP:7890
- **Uptime Kuma Monitoring**: http://YOUR_SERVER_IP:3001
- **Netdata System Monitoring**: http://YOUR_SERVER_IP:19999

### 9. Security Considerations

#### Firewall Configuration
```bash
# Allow only necessary ports
sudo ufw allow 22    # SSH
sudo ufw allow 80    # HTTP
sudo ufw allow 443   # HTTPS
sudo ufw --force enable

# Optional: Allow monitoring ports only from specific IPs
# sudo ufw allow from YOUR_ADMIN_IP to any port 3001
# sudo ufw allow from YOUR_ADMIN_IP to any port 7890
# sudo ufw allow from YOUR_ADMIN_IP to any port 19999
```

#### Secure Monitoring Access

For production, consider restricting monitoring access:

1. **Configure Nginx authentication for monitoring endpoints**:
   ```bash
   # Generate password file
   sudo apt install apache2-utils
   sudo htpasswd -c /etc/nginx/.htpasswd admin
   ```

2. **Add authentication to nginx.conf**:
   ```nginx
   # Add to monitoring locations
   location /monitoring/ {
       auth_basic "Monitoring Area";
       auth_basic_user_file /etc/nginx/.htpasswd;
       proxy_pass http://monitoring_backend/;
   }
   ```

### 10. Maintenance Commands

#### Update the application:
```bash
# Pull latest changes
git pull origin master

# Rebuild and update services
docker service update --image zeepass_zeepass:latest zeepass_zeepass
```

#### Scale services:
```bash
# Scale ZeePass application
docker service scale zeepass_zeepass=3

# Scale Nginx
docker service scale zeepass_nginx=2
```

#### Backup Redis data:
```bash
# Create backup
docker exec $(docker ps -q -f name=zeepass_redis) redis-cli BGSAVE

# Copy backup file
docker cp $(docker ps -q -f name=zeepass_redis):/data/dump.rdb ./backup-$(date +%Y%m%d).rdb
```

#### View logs:
```bash
# Application logs
docker service logs -f zeepass_zeepass

# Nginx access logs  
docker service logs -f zeepass_nginx

# Redis logs
docker service logs -f zeepass_redis
```

### 11. SSL Certificate Renewal (Let's Encrypt)

Set up automatic renewal:
```bash
# Test renewal
sudo certbot renew --dry-run

# Add to crontab for automatic renewal
echo "0 12 * * * /usr/bin/certbot renew --quiet" | sudo crontab -
```

After certificate renewal, restart nginx:
```bash
docker service update --force zeepass_nginx
```

### 12. Troubleshooting

#### Common Issues:

1. **SSL Certificate errors**: Check certificate paths and permissions
2. **Service won't start**: Check logs with `docker service logs`
3. **Domain not accessible**: Verify DNS settings and firewall rules
4. **Performance issues**: Check system resources with Netdata dashboard

#### Health Checks:
```bash
# Check nginx configuration
docker exec $(docker ps -q -f name=zeepass_nginx) nginx -t

# Check Redis connection
docker exec $(docker ps -q -f name=zeepass_redis) redis-cli ping

# Check ZeePass health endpoint
curl -k https://zeepass.com/health
```

### 13. Monitoring Setup

The deployment includes several monitoring tools:

1. **GoAccess** (Port 7890): Real-time web log analyzer
   - View traffic patterns and visitor analytics
   - Access via: http://YOUR_SERVER_IP:7890

2. **Uptime Kuma** (Port 3001): Uptime monitoring
   - Set up website monitoring alerts
   - Access via: http://YOUR_SERVER_IP:3001

3. **Netdata** (Port 19999): System performance monitoring
   - Real-time system metrics
   - Access via: http://YOUR_SERVER_IP:19999

### 14. Production Checklist

- [ ] SSL certificates installed and configured
- [ ] Domain DNS pointing to server
- [ ] Firewall configured with minimal necessary ports
- [ ] Redis password set (optional)
- [ ] Monitoring services accessible
- [ ] Log rotation configured
- [ ] Backup strategy implemented
- [ ] SSL certificate auto-renewal configured
- [ ] Health checks passing
- [ ] Load testing completed

## Support

For issues or questions:
1. Check the logs using the commands above
2. Review the troubleshooting section
3. Consult the main ZeePass documentation
4. Open an issue on the GitHub repository

---

**Important**: Always test deployments in a staging environment before deploying to production.