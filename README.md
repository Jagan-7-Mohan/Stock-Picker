# Stock Picker - IPO Notification Bot

A Go-based WhatsApp bot that sends daily notifications about open IPOs in the Indian stock market, including Grey Market Premium (GMP) and subscription details.

## Features

- 📈 Daily automated notifications at 6:00 AM IST
- 💰 IPO details including price range, lot size, and exchange
- 📊 Grey Market Premium (GMP) information
- 📈 Subscription details and status
- 📱 WhatsApp integration via Twilio
- 🔄 Automatic filtering of only open IPOs

## Project Structure

The project follows standard Go project layout:

```
stock-picker/
├── cmd/
│   └── stock-picker/      # Main application entry point
│       └── main.go
├── internal/               # Private application code
│   ├── config/            # Configuration management
│   ├── fetcher/           # IPO data fetching logic
│   ├── models/            # Data models
│   └── whatsapp/          # WhatsApp integration
├── configs/               # Configuration files
│   └── .env.example
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── .gitignore
```

## Prerequisites

1. **Go 1.21 or higher** - [Install Go](https://golang.org/doc/install)
2. **Twilio Account** - Sign up at [Twilio](https://www.twilio.com/)
   - Get your Account SID and Auth Token
   - Set up WhatsApp Sandbox or get a WhatsApp Business API number

## Setup

### 1. Clone and Install Dependencies

```bash
cd stock-picker
go mod download
```

### 2. Configure Environment Variables

Copy the example environment file and fill in your details:

```bash
cp configs/.env.example .env
```

Edit `.env` in the project root with your Twilio credentials:

```env
TWILIO_ACCOUNT_SID=your_account_sid_here
TWILIO_AUTH_TOKEN=your_auth_token_here
TWILIO_WHATSAPP_FROM=whatsapp:+14155238886
WHATSAPP_RECIPIENTS=whatsapp:+919876543210
```

### 3. Twilio WhatsApp Setup

#### For Testing (Sandbox):
1. Go to [Twilio Console](https://console.twilio.com/)
2. Navigate to Messaging > Try it out > Send a WhatsApp message
3. Follow the instructions to join the sandbox
4. Use the sandbox number as `TWILIO_WHATSAPP_FROM`
5. Add your WhatsApp number to `WHATSAPP_RECIPIENTS` (must join sandbox first)

#### For Production:
1. Apply for WhatsApp Business API access in Twilio Console
2. Get your approved WhatsApp Business number
3. Use that number in `TWILIO_WHATSAPP_FROM`

### 4. Build and Run

```bash
# Build the application
make build
# or
go build -o stock-picker ./cmd/stock-picker

# Run the application
make run
# or
./stock-picker
```

Or run directly:

```bash
make run-dev
# or
go run ./cmd/stock-picker
```

## How It Works

1. **Scheduler**: Uses cron to run daily at 6:00 AM IST
2. **Data Fetching**: Scrapes IPO data from multiple sources:
   - Chittorgarh.com for IPO listings
   - GMP share websites for Grey Market Premium
3. **Filtering**: Only includes IPOs that are currently open
4. **Notification**: Formats and sends WhatsApp messages with all IPO details

## Message Format

The bot sends formatted messages like:

```
📈 Daily IPO Update - 15 Jan 2024

🟢 Open IPOs:

━━━━━━━━━━━━━━━━━━━━
Company Name
━━━━━━━━━━━━━━━━━━━━
📅 IPO Window: 10 Jan 2024 to 15 Jan 2024
💰 Price Range: ₹100-110
📊 GMP: ₹15 (+13.6%)
📈 Subscription: 2.5x
📦 Lot Size: 13 shares
🏛️ Exchange: NSE, BSE
ℹ️ Info: Technology company
```

## Data Sources

The bot currently fetches data from:
- Chittorgarh.com - IPO listings and basic details
- GMP Share - Grey Market Premium data

**Note**: You may need to adjust the scraping logic if these websites change their structure. Consider using official APIs if available.

## Customization

### Change Schedule Time

Edit `cmd/stock-picker/main.go` and modify the cron expression:

```go
// Current: 6:00 AM IST daily
_, err := c.AddFunc("0 6 * * *", func() {
    // ...
})

// Example: 8:00 AM IST daily
_, err := c.AddFunc("0 8 * * *", func() {
    // ...
})
```

### Add More Data Sources

Extend `internal/fetcher/ipo_fetcher.go` to add more data sources:

```go
func (f *IPOFetcher) fetchFromNewSource() ([]IPO, error) {
    // Your implementation
}
```

Then add it to the `sources` slice in `FetchOpenIPOs()`.

## Troubleshooting

### Messages Not Sending
- Verify Twilio credentials are correct
- Check that recipient numbers are in correct format: `whatsapp:+919876543210`
- For sandbox: Ensure recipient has joined the sandbox
- Check Twilio console for error logs

### No IPOs Found
- Verify internet connection
- Check if data source websites are accessible
- Review scraping logic if website structure changed

### Date Parsing Issues
- IPO dates might be in different formats
- Adjust `parseDate()` function in `internal/fetcher/ipo_fetcher.go` to support more formats

## Running as a Service

### Using systemd (Linux)

Create `/etc/systemd/system/stock-picker.service`:

```ini
[Unit]
Description=Stock Picker IPO Bot
After=network.target

[Service]
Type=simple
User=your-user
WorkingDirectory=/path/to/stock-picker
ExecStart=/path/to/stock-picker/stock-picker
Restart=always
RestartSec=10
Environment="PATH=/usr/local/bin:/usr/bin:/bin"

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl enable stock-picker
sudo systemctl start stock-picker
```

### Using Docker

Create `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o stock-picker

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/stock-picker .
COPY --from=builder /app/.env .
CMD ["./stock-picker"]
```

Build and run:

```bash
docker build -t stock-picker .
docker run -d --name stock-picker-bot --env-file .env stock-picker
```

## License

MIT License - feel free to use and modify as needed.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Disclaimer

This bot is for informational purposes only. Always verify IPO information from official sources before making investment decisions. The bot scrapes data from public sources and accuracy is not guaranteed.

