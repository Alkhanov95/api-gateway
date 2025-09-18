# 🚀 API Gateway – User Service

**API Gateway – User Service** is a backend microservice written in **Go** with the **Fiber** framework.  
It follows the principles of **Clean Architecture**, where the application is divided into independent layers —  
`handler → usecase → repository → storage` — ensuring maintainability, scalability, and easy extension.

---

## 📖 Description

This project provides a simple **REST API for user management** and is designed as a building block for larger backend systems.  

### ✨ Features
- **Create** new users with auto-generated UUIDs  
- **Retrieve** users by ID or list all users  
- **Update** existing user data  
- **Delete** users  
- **Config management** via `config.yaml`  
- **Structured logging** using `slog`  
- **Custom error handling** (`apperr.ErrNotFound`)  
- **Clean separation of concerns** for easy scaling and extension  

Planned improvements:  
- [ ] Caching (Redis / in-memory)  
- [ ] CI/CD pipelines (GitHub Actions + Docker)  
- [ ] Integration & unit testing  
- [ ] Extended modular services  

---

## 📂 Project Structure

```bash
.
├── app.go             # Entry point
├── config.go          # Config loader
├── config.yaml        # Example configuration
├── router.go          # HTTP router setup
├── internal/          
│   ├── handler/       # HTTP handlers (Fiber)
│   ├── usecase/       # Business logic (UserProvider)
│   ├── repository/    # Data access layer
│   ├── models/        # Domain models (User, etc.)
│   └── apperr/        # Custom application errors
├── storage.go         # Storage layer
├── cache.go           # (Planned) caching logic
└── tests/             # (Planned) integration & unit tests
