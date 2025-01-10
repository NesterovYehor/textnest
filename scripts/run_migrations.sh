#!/bin/sh

set -e  # Exit immediately if a command exits with a non-zero status

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
until nc -z postgres 5432; do
  echo "Waiting for PostgreSQL..."
  sleep 1
done
echo "PostgreSQL is up!"

# Run the migration tool
echo "Running migrations..."
migrate -path /migrations -database "postgres://textnest:textnest_password@postgres:5432/textnestdb?sslmode=disable" up

# Check if migrations succeeded
if [ $? -eq 0 ]; then
  echo "Migrations completed successfully!"
else
  echo "Migrations failed!"
  exit 1
fi


