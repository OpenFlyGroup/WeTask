#!/bin/sh
set +e

echo "🚀 Starting service..."

# Wait for database to be ready
echo "⏳ Waiting for database connection..."
sleep 5

# Execute the main command
echo "▶️  Executing: $@"
exec "$@"
