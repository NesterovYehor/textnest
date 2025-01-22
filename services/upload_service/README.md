# Upload Service for TextNest

## Overview

The **Upload Service** is a core component of the TextNest project, designed to handle the ingestion and storage of user-submitted paste content. It ensures efficient, secure, and reliable processing of uploads while integrating seamlessly with other system components.

## Features

- **Content Uploading**: Facilitates the secure and efficient upload of paste content.
- **Metadata Integration**: Automatically associates metadata with uploaded content for better organization and access control.
- **Validation**: Validates uploaded data to ensure compliance with system rules and limits.

## Architecture

### Service Communication

The Upload Service integrates with other components of the TextNest ecosystem:

- **Metadata Service**: Sends metadata about uploaded content for storage and indexing.
- **Storage Layer**: Transfers and stores uploaded paste content in the designated storage system.

### Workflow

1. Receive paste content and user metadata via the API.
2. Validate the uploaded content to ensure it meets system requirements (e.g., size, format).
3. Transfer the content to the storage layer and associate it with the generated key.
4. Notify the Metadata Service to index the paste for retrieval.

## Logging

The service utilizes [slog](https://pkg.go.dev/log/slog) for comprehensive logging. Logged details include:

- Incoming upload requests
- Content validation results
- Interactions with the Key Generation Service
- Storage transfer statuses
- Errors and anomalies for debugging

## Future Enhancements

- Add support for resumable uploads to handle large paste content efficiently.


- Introduce upload rate-limiting to prevent abuse.
- Add monitoring and analytics for upload performance metrics.

## License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/NesterovYehor/textnest/blob/main/LICENSE) file for details.

