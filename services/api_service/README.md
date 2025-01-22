# API Service for TextNest

## Overview

The **API Service** is a core component of the TextNest project, acting as the main entry point for client requests. It is responsible for managing interactions between the frontend and backend services. This service provides RESTful APIs and communicates with other microservices using gRPC.

## Features

- **Key generation**: Connects with the Hash Generator service to create and validate unique identifiers for pastes.
- **Content Management and Metadata Storage**: Interfaces with the Upload and Download services to store and retrieve paste content along with associated metadata.
- **Scalable Architecture**: Designed to scale independently within the TextNest microservices ecosystem.

## Architecture

### Service Communication

The API Service communicates with other services using gRPC protocols to ensure efficient data transfer and service orchestration:

- **Hash Generator**: Generates and validates unique identifiers for pastes.
- **Upload and Download Services**: Handles both storing and retrieving paste content and metadata.
- **Expiration Service**: Subscribes to Kafka notifications for handling expiring data.

## Logging

The service leverages [slog](https://pkg.go.dev/log/slog) for structured and comprehensive logging. Logged details include:

- API request and response details
- gRPC communication activities
- Errors and exceptions for debugging

## Future Enhancements

- Implement authentication and authorization mechanisms.
- Add rate limiting to prevent abuse.
- Enhance monitoring capabilities using AWS CloudWatch or Prometheus.

## License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/NesterovYehor/textnest/blob/main/LICENSE) file for details.

