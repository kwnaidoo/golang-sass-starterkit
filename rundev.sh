#!/bin/bash
set -e

export MYSQL_HOST=127.0.0.1
export MYSQL_PORT=3306
export MYSQL_DATABASE=gosass
export MYSQL_USER=kevin
export MYSQL_PASSWORD=1234
export gosass_URL=http://127.0.0.1:3001
export ALLOWED_IPS="127.0.0.1,116.202.156.187"
export REDIS_DSN="127.0.0.1:6379"
export ENCRYPTION_KEY="*~#^2^#s0^=)^^7%b34@#$%1"
export SMTP_HOST="sandbox.smtp.mailtrap.io"
export SMTP_USERNAME="xxx"
export SMTP_PASSWORD="xxx"
export SMTP_PORT="2525"
export SMTP_FROM_EMAIL="gosass <test@mailtrap.io>"
export TZ="UTC"
export ALLOW_REGISTER=true

go run main.go