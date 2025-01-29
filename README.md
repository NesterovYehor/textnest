# 🚀 TextNest – The Terminal-Powered Code Sharing Tool

## 🧐 The Story Behind TextNest  

It all started when I stumbled upon [terminal.shop](https://terminal.shop) (by the way, they serve the **best coffee in your terminal!** ☕). It was such a cool and effortless experience that it sparked an idea:

> *What if there was a tool for developers who are tired of building UIs but still need a fast, secure way to share code snippets?*

That’s how **TextNest** was born – a **developer-first, terminal-friendly tool** designed to make text and code sharing **blazing fast, secure, and scalable** without the clutter of traditional interfaces.

---

## 📌 Why TextNest?

### 🔥 Designed for Speed, Security & Simplicity

✅ **Fast & Scalable** – Built with microservices for maximum performance.  
✅ **Secure** – JWT authentication ensures safe access.  
✅ **Developer-Friendly** – API-first approach, with an **SSH interface coming soon!**  
✅ **Reliable Storage** – Metadata in PostgreSQL, caching with Redis, and files in Amazon S3.  
✅ **Minimalist, Terminal-First Approach** – No need for UIs. Just focus on **getting things done.**  

---

## 🛠️ Tech Stack – What’s Under the Hood?

| **Technology**  | **Purpose** |
|---------------|------------|
| **Go** | High-performance language, ideal for microservices. |
| **gRPC** | Efficient, lightweight service-to-service communication. |
| **PostgreSQL** | Secure storage for metadata. |
| **Redis** | In-memory caching for speed. |
| **Amazon S3** | Reliable, scalable file storage. |
| **Kafka** | Event-driven architecture for automated paste expiration. |
| **Docker** | Simplifies deployment and scalability. |
| **Nginx** | Load balancing for high availability. |
| **JWT (Auth Service)** | Ensures secure authentication. |

With this tech stack, **TextNest is built to perform at scale while ensuring security and developer efficiency.**

---

## 🔧 System Architecture – How It Works

TextNest is a **microservices-based** application with distinct services:

🛜 **API Service** – Handles all HTTP requests.  
🔐 **Auth Service** – Manages user authentication via JWT.  
📤 **Upload Service** – Securely stores pastes in S3.  
📥 **Download Service** – Retrieves stored pastes securely.  
🔑 **Hash Generator** – Generates unique IDs for pastes.  
⏳ **Expiration Service** – Deletes expired pastes using Kafka.  
💾 **Metadata Storage** – PostgreSQL (structured data) + Redis (caching).  

Everything runs in **Docker**, making deployment a breeze! 🐳

---

## 🚀 Getting Started – Set Up in Minutes

### 📂 Clone the Repository

```bash
git clone https://github.com/NesterovYehor/textnest.git
cd textnest
```

### ⚙️ Configure the Services

Most services use `.yaml` configuration files. Example API config:

```yaml
database:
  user: youruser
  password: yourpassword
s3:
  bucket: textnest-bucket
redis:
  url: redis://localhost:6379
auth:
  service_url: http://auth-service:8081
jwt:
  secret: your_secret_key
```

### 🏃‍♂️ Run Everything with Docker

```bash
docker-compose up -d
```

Once running, your API should be available at:

```
http://localhost:8080
```

For authentication, the Auth Service runs on:

```
http://localhost:8081
```

---

## 📡 Using the API – Quick Example

**Create a new snippet:**
```bash
curl -X POST http://localhost:8080/api/snippet \
     -H "Authorization: Bearer <your_token>" \
     -d '{"text": "Hello, World!", "expires_in": 3600}'
```

**Retrieve a snippet:**
```bash
curl -X GET http://localhost:8080/api/snippet/<paste_id>
```

**Delete a snippet:**
```bash
curl -X DELETE http://localhost:8080/api/snippet/<paste_id>
```

Simple. Fast. Secure. ✅

---

## 🔮 What’s Next?

🚀 **SSH Interface** – Share and retrieve pastes directly from the terminal.
🎛️ **Admin Dashboard** – For advanced management and monitoring.
🔄 **Public & Private Snippets** – More control over who sees what.

---

## 🐝 License

TextNest is open-source and licensed under the **MIT License**. See [LICENSE](./LICENSE) for details.

---

## 📬 Contact

📧 **Email:** [yehor.nesterov@example.com](mailto:yehor.nesterov@example.com)  
🤖 **GitHub Issues:** [Open an Issue](https://github.com/NesterovYehor/textnest/issues)  

---

🔥 **TextNest – Fast. Secure. Developer-Friendly.** 🔥

