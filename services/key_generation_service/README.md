# Key Generation Service for TextNest

## Overview

The **Key Generation Service** is a critical component of the TextNest project, responsible for securely generating and managing unique keys. These keys are stored in an in-memory database, enabling fast access and high performance when serving requests from other services.

## Features

- **Key Generation**: Generates secure, collision-resistant keys for identifying pastes.
- **Key Regeneration**: Ensures keys are pre-generated and stored in an in-memory database to optimize performance during high-demand scenarios.
- **High Performance**: Leverages in-memory storage for rapid key retrieval.

## Architecture

### Service Communication

The Key Generation Service interacts with other components of the TextNest ecosystem to:

- **Provide Keys**: Supplies unique keys to services that require them for identifying resources.
- **Optimize Performance**: Pre-generates and stores keys in memory, minimizing latency during operations.

### Workflow

1. Generate a batch of unique keys and store them in an in-memory database.
2. Handle requests for keys by retrieving them directly from the in-memory database.
3. Regenerate keys as needed to maintain availability.

## Logging

The service utilizes [slog](https://pkg.go.dev/log/slog) for robust logging. Logged details include:

- Key generation and regeneration events
- Key retrieval requests
- Errors and anomalies in the key generation or storage process

## Future Enhancements

- Introduce monitoring tools to track key generation and retrieval performance.
- Enhance fault tolerance to ensure availability during in-memory database failures.

## License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/NesterovYehor/textnest/blob/main/LICENSE) file for details.
