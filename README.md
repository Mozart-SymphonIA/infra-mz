# Infra-mz

Infra-mz is a lightweight, high-performance toolkit designed for building robust Go microservices. It focuses on simplifying the common infrastructure needs of modern distributed systems while adhering to Clean Architecture and Domain-Driven Design (DDD) principles.

## Features

Infra-mz is composed of three main architectural pillars:

- **`busx`**: A clean abstraction for RabbitMQ that facilitates asynchronous communication and event-driven patterns. It is designed to work seamlessly within Clean Architecture and DDD contexts.
- **`dbx`**: A powerful wrapper around `sqlx` that simplifies database operations, providing a more ergonomic and safer way to interact with SQL databases.
- **`bootx`**: A standardized service bootstrapping library that handles the heavy lifting of setting up HTTP servers, gRPC servers, and worker processes, ensuring consistent behavior across all your microservices.

## Installation

You can easily add Infra-mz to your project using `go get`:

```bash
go get github.com/Mozart-SymphonIA/infra-mz
```

## Why Infra-mz?

Building microservices often involves repetitive boilerplate code for infrastructure. Infra-mz aims to eliminate this redundancy, allowing developers to focus on what matters most: the core business logic. Whether you're building a simple API or a complex event-driven system, Infra-mz provides the building blocks you need to move fast and maintain high standards of code quality.

---

We believe in the power of the open-source community. If you find this project useful, feel free to contribute or give it a star!
