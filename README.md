# Message & Notification Service — FreelanceX

## Overview
Handles internal user messaging and platform notifications for proposal and websocket messaging using SMTP packages and Kafka consumer events.

## Tech Stack
- Go (Golang)
- gRPC
- Redis (for pub-sub or notifications)
- Protocol Buffers
- MongoDB
- Kafka
- SMTP
- CRON 

## Setup

### 1. Clone & Navigate
```bash
git clone https://github.com/Prototype-1/freelancex_timeTrancker_service.git
cd freelancex_project/message.notification_service
```

### Install Dependencies

go mod tidy

## Create .env

PORT=50055
REDIS_ADDR=localhost:6379

## Start the Service

go run main.go

### Proto Definitions

    proto/message/message.proto

## Notes

    Publishes notifications to Redis channels.

    Future support for email/SMS integrations.

#### Maintainers

aswin100396@gmail.com
