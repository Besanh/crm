#!/bin/sh
echo "Build go application"
GOOS=linux GOARCH=amd64 go build -o app.exe main.go
echo "Restart service"
systemctl restart crm-service
systemctl status crm-service