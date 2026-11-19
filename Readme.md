
# System Design Handbook

A collection of Low-Level Design (LLD), High-Level Design (HLD),
and distributed systems concepts with practical explanations,
architecture diagrams, and Go/Java implementations.

## About This Repository

This repository documents my journey of exploring and
implementing scalable backend systems.

My focus areas include:
- Distributed Systems
- Backend Architecture
- Concurrency and Thread Safety
- Scalability and Performance
- Database Design
- Fault Tolerance

## 📚 Core Concepts

| Topic | Link | Core Implementation|
|---|---|---|
| Caching | [Read](https://redis.io/blog/cache-layer-architecture-guide) | [LRU cache diagram](lld/caching/LRU/LRU.excalidraw.svg)
| Consistent Hashing | [Read](https://bytebytego.com/courses/system-design-interview/design-consistent-hashing) | [Solution](lld/consistant-hashing/BST-implementation.go)
| Database Indexing | [Read](concepts/database-indexing/README.md) |
| Concurrency | [Read](concepts/concurrency/README.md) |
| Message Queues |  | [Solution](https://github.com/shubhamjsh/Message-Queue/)
| Idempotency | [Read](https://algomaster.io/learn/system-design/idempotency) | [Solution](lld/idempotency/payment.go)
| Merkle Tree | [Read](lld/tamper-detection/Merkle-tree.excalidraw.svg) | [Hash Chain](lld/tamper-detection/hash-chain/main.go) [Merkle Tree](lld/tamper-detection/merkle-tree/main.go)


## 🧩 Low-Level Design (LLD)

| Problem | Key Concepts | Solution |
|---|---|---|
| LRU Cache | HashMap, Linked List, Concurrency | [View](lld/caching/LRU/LRU.go) |
| Chat app | gRPC stream | [View](lld/grpc-chat)
| Rate Limiter | Token Bucket, Thread Safety | [View](lld/rate-limiter/README.md) |
| Task Scheduler | Worker Pool, Retry, Cancellation | [View](lld/task-scheduler/README.md) |
| In-Memory KV Store | TTL, Concurrency, Storage | [View](lld/in-memory-kv-store/README.md) |

## 🏗️ High-Level Design (HLD)

| Problem | Key Concepts | Design |
|---|---|---|
| Distributed Cache | Sharding, Replication | [View](hld/distributed-cache/README.md) |
| Notification System | Kafka, Retry, Fan-out | [View](hld/notification-system/README.md) |
| Payment System | Idempotency, Saga | [View](hld/payment-system/README.md) |
| URL Shortener | Hashing, Caching | [View](hld/url-shortener/README.md) |

## 🛠️ Technologies

- Languages: Go, Java
- Databases: PostgreSQL, Redis
- Messaging: Kafka
- Infrastructure: Docker, Kubernetes
- Diagrams: Excalidraw, Mermaid

## Design Approach

Each problem follows a consistent structure:

1. Requirements
2. Basic approach
3. Optimized design
4. Architecture diagrams
5. Data models and APIs
6. Concurrency and failure handling
7. Performance and scalability
8. Trade-offs and future improvements

## Author

Backend and Distributed Systems Engineer with 9+ years
of experience building scalable production systems.

[LinkedIn](https://www.linkedin.com/in/shubham-joshi-0b1b21102/)
