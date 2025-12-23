# Deployment Guide

This guide covers deploying the Stock Picker IPO notification bot to various cloud platforms.

## Prerequisites

Before deploying, ensure you have:
1. **Twilio Account** with WhatsApp API access
2. **Environment Variables** ready:
   - `TWILIO_ACCOUNT_SID`
   - `TWILIO_AUTH_TOKEN`
   - `TWILIO_WHATSAPP_FROM` (format: `whatsapp:+14155238886`)
   - `WHATSAPP_RECIPIENTS` (comma-separated, format: `whatsapp:+919876543210,whatsapp:+919876543211`)

## Deployment Options

### Option 1: Railway (Recommended - Easiest)

[Railway](https://railway.app) is the easiest platform for deploying Go applications.

#### Steps:

1. **Sign up** at [railway.app](https://railway.app) (GitHub login available)

2. **Create a new project**:
   - Click "New Project"
   - Select "Deploy from GitHub repo" (or upload code)

3. **Configure environment variables**:
   - Go to your project → Variables tab
   - Add all required environment variables:
     ```
     TWILIO_ACCOUNT_SID=your_account_sid
     TWILIO_AUTH_TOKEN=your_auth_token
     TWILIO_WHATSAPP_FROM=whatsapp:+14155238886
     WHATSAPP_RECIPIENTS=whatsapp:+919876543210
     TZ=Asia/Kolkata
     ```

4. **Deploy**:
   - Railway will automatically detect the Dockerfile and deploy
   - The app will start and run continuously

5. **Monitor**:
   - Check the "Deployments" tab for logs
   - The app runs 24/7 and sends notifications daily at 8:00 AM IST

**Pricing**: Railway offers a free tier with $5 credit/month. The app should run within free tier limits.

---

### Option 2: Render

[Render](https://render.com) offers a free tier perfect for background services.

#### Steps:

1. **Sign up** at [render.com](https://render.com) (GitHub login available)

2. **Create a new Web Service**:
   - Click "New +" → "Web Service"
   - Connect your GitHub repository
   - Select the `stock-picker` repository

3. **Configure the service**:
   - **Name**: `stock-picker`
   - **Environment**: `Docker`
   - **Region**: Choose closest to India (e.g., Singapore)
   - **Branch**: `main` (or your default branch)
   - **Root Directory**: Leave empty
   - **Dockerfile Path**: `Dockerfile`
   - **Docker Context**: `.`

4. **Set environment variables**:
   - Scroll to "Environment Variables" section
   - Add:
     ```
     TWILIO_ACCOUNT_SID=your_account_sid
     TWILIO_AUTH_TOKEN=your_auth_token
     TWILIO_WHATSAPP_FROM=whatsapp:+14155238886
     WHATSAPP_RECIPIENTS=whatsapp:+919876543210
     TZ=Asia/Kolkata
     ```

5. **Deploy**:
   - Click "Create Web Service"
   - Render will build and deploy automatically
   - Check logs to ensure it's running

**Note**: On free tier, Render spins down services after 15 minutes of inactivity. For a cron job, consider upgrading to a paid plan or use a service like [cron-job.org](https://cron-job.org) to ping your service every 10 minutes.

**Pricing**: Free tier available, but may spin down. Paid plans start at $7/month.

---

### Option 3: Fly.io

[Fly.io](https://fly.io) is excellent for Go applications with global deployment.

#### Steps:

1. **Install Fly CLI**:
   ```bash
   curl -L https://fly.io/install.sh | sh
   ```

2. **Login**:
   ```bash
   fly auth login
   ```

3. **Launch the app**:
   ```bash
   fly launch
   ```
   - Follow the prompts
   - Choose a region close to India (e.g., `bom` for Mumbai)
   - Don't deploy yet (we need to set env vars first)

4. **Set environment variables**:
   ```bash
   fly secrets set TWILIO_ACCOUNT_SID=your_account_sid
   fly secrets set TWILIO_AUTH_TOKEN=your_auth_token
   fly secrets set TWILIO_WHATSAPP_FROM=whatsapp:+14155238886
   fly secrets set WHATSAPP_RECIPIENTS=whatsapp:+919876543210
   fly secrets set TZ=Asia/Kolkata
   ```

5. **Deploy**:
   ```bash
   fly deploy
   ```

6. **Monitor**:
   ```bash
   fly logs
   ```

**Pricing**: Free tier includes 3 shared-cpu VMs with 256MB RAM. Perfect for this app.

---

### Option 4: Docker (Any Platform)

You can deploy the Docker container to any platform that supports Docker:

#### Build locally:
```bash
docker build -t stock-picker .
```

#### Run locally:
```bash
docker run -d \
  --name stock-picker \
  --env-file .env \
  stock-picker
```

#### Deploy to platforms like:
- **DigitalOcean App Platform**: Upload Dockerfile, set env vars
- **AWS ECS/Fargate**: Use Dockerfile, configure task definition
- **Google Cloud Run**: Deploy container, set env vars
- **Azure Container Instances**: Deploy Docker container

---

## Post-Deployment Checklist

- [ ] Verify environment variables are set correctly
- [ ] Check application logs for startup errors
- [ ] Test by manually triggering (the app runs once on startup)
- [ ] Verify cron schedule (8:00 AM IST daily)
- [ ] Monitor for first scheduled notification
- [ ] Set up log monitoring/alerts (optional)

## Monitoring & Logs

### Railway
- View logs in the "Deployments" tab
- Real-time log streaming available

### Render
- View logs in the service dashboard
- Logs tab shows real-time output

### Fly.io
```bash
fly logs          # View logs
fly logs -a stock-picker  # View logs for specific app
```

## Troubleshooting

### App not starting
- Check environment variables are set correctly
- Verify Twilio credentials are valid
- Check logs for specific error messages

### Notifications not sending
- Verify WhatsApp recipients are in correct format: `whatsapp:+919876543210`
- Check Twilio console for message status
- Ensure Twilio account has WhatsApp API access
- For sandbox: Recipients must join the sandbox first

### Timezone issues
- Ensure `TZ=Asia/Kolkata` is set
- Verify cron schedule matches IST timezone
- Check logs for timezone-related warnings

### Service keeps restarting
- Check memory limits (should be at least 256MB)
- Review logs for panic/error messages
- Verify all dependencies are available

## Cost Estimation

- **Railway**: ~$0-5/month (free tier usually sufficient)
- **Render**: $0-7/month (free tier may spin down)
- **Fly.io**: $0/month (free tier sufficient)
- **DigitalOcean**: ~$5/month (basic droplet)

## Recommended Platform

**Railway** is recommended for:
- ✅ Easiest setup
- ✅ Automatic deployments from GitHub
- ✅ Good free tier
- ✅ No spin-down issues
- ✅ Easy environment variable management

---

## Need Help?

If you encounter issues:
1. Check the application logs
2. Verify all environment variables
3. Test Twilio credentials independently
4. Review the main README.md for local setup troubleshooting

