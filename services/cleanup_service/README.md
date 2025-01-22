# Cleanup Service for TextNest

## Overview

The **Cleanup Service** is a crucial part of the TextNest project, responsible for managing the expiration of pastes. It ensures that expired pastes are promptly identified and removed, maintaining the efficiency and reliability of the system.

## Features

- **Expiration Check**: Periodically checks for expired pastes in the system.
- **Notification System**: Sends notifications to other services about expiring data via Kafka.
- **System Optimization**: Enhances performance by removing stale data.

## Architecture

### Service Communication

The Cleanup Service interacts with other components of the TextNest ecosystem:

- **Expiration Notifications**: Publishes messages to Kafka about expiring data.
- **Dependent Services**: Works in conjunction with other services like Metadata and Content Management to trigger deletion of expired pastes.

### Workflow

1. Periodic scan for expired pastes.
2. Notify other services about expiring data via Kafka.
3. Coordinate deletion across the system.

## Logging

The service utilizes [slog](https://pkg.go.dev/log/slog) for robust logging. Logged details include:

- Scheduled expiration checks
- Kafka notification logs
- Errors and anomalies for debugging

## Future Enhancements

- Introduce configurable expiration policies.
- Add monitoring for expiration performance metrics.
- Improve system resilience with retry mechanisms for Kafka notifications.

## License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/NesterovYehor/textnest/blob/main/LICENSE) file for details.
