# Use a lightweight base image
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy the migration script into the container
COPY ./run_migrations.sh /scripts/run_migrations.sh
RUN chmod +x /scripts/run_migrations.sh

# Install necessary tools for migrations
RUN apt-get update && apt-get install -y \
    netcat-openbsd curl \
    && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz \
    && mv migrate /usr/local/bin \
    && apt-get clean

# Default entrypoint for the migration container
ENTRYPOINT ["/scripts/run_migrations.sh"]

