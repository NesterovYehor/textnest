# TextNest

TextNest is a backend-focused project designed to handle the storage and management of text-based data efficiently. It allows users to interact with the system via SSH and provides features like anonymous usage and future support for organized data management. Built with a modern microservices architecture, TextNest leverages technologies like Go, gRPC, Protobufs, Redis, PostgreSQL, and Kafka to ensure scalability, reliability, and performance.

---

## Features

- **Anonymous Access**: Users can interact with the system anonymously via SSH, making it simple and fast to use without requiring authentication.
- **Microservices Architecture**: Modular design with specialized services for handling different functionalities such as uploads, downloads, and metadata management.
- **Fast and Reliable Communication**: 
  - **gRPC and Protobufs** are used for efficient and high-performance communication between microservices.
  - **Kafka** is selectively employed to manage interactions requiring coordination across multiple services or where data consistency and safety are prioritized.
- **Efficient Storage**: Uses Amazon S3 for blob storage and PostgreSQL for metadata storage.
- **In-Memory Caching**: Redis accelerates access to frequently used data.
- **Scalability**: Designed with scalability in mind, leveraging Docker and Kubernetes (optional) for containerized deployment.

---

## Technologies Used

TextNest is built with a modern stack of technologies to deliver high performance and maintainability:

### Programming Language
- **Go**: The main language for implementing the backend, chosen for its strong concurrency model, simplicity, and performance.

### Communication
- **gRPC and Protobufs**: The primary means of communication between microservices, offering high performance, type safety, and lightweight message serialization.
- **Kafka**: Used strategically for asynchronous messaging in scenarios that involve multiple services or require data safety, ensuring consistent and reliable processing.

### Data Storage
- **PostgreSQL**: Manages structured metadata storage, offering robustness and ACID compliance.
- **Amazon S3**: Provides scalable and durable blob storage for text files.

### Caching
- **Redis**: Improves performance by caching frequently accessed data.

### Configuration Management
- **YAML**: Configuration files manage service settings, including Redis and PostgreSQL connections.

---

## Architecture Overview

TextNest's backend is organized into multiple microservices to ensure modularity and scalability. Key components include:

1. **API Service**: Central gateway for user interactions, delegating tasks to other services via gRPC.
2. **Key Generator Service(KGS)**: Generates unique identifiers for each paste files.
3. **Upload Service**: Handles the storage of user data to Amazon S3 and updates metadata in PostgreSQL.
4. **Download Service**: Manages the retrieval of stored data and delivers it to users.
5. **Clean up Service**: Periodically identifies expired content and triggers clean-up operations across services.
6. **gRPC and Protobufs**: Facilitate seamless and high-performance communication between all microservices.
7. **Kafka**: Used to orchestrate complex operations involving multiple services, such as clean-up tasks triggered by the Expiration Service, ensuring data safety and consistency during concurrent operations.

This architecture ensures a balance of high-speed communication through gRPC and Protobufs while leveraging Kafka for critical operations requiring robust message handling and service coordination.

---

## Usage

Users can connect to TextNest via SSH for a simple and efficient interface to create and manage pastes anonymously. A sample SSH client is available [here](https://github.com/NesterovYehor/txtnest-cli).

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Contact

For inquiries or feedback, contact via GitHub: [NesterovYehor](https://github.com/NesterovYehor) or [Discord](https://discordapp.com/users/591678870973841428).
