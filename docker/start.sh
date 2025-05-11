#!/bin/sh

echo "🚀 Starting Daxa runtime on :36365 & :8080"
daxagrid-runtime &

echo "🌐 Starting Nginx on :443 and :8080"
nginx -g "daemon off;"
