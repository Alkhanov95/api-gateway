# 🚀 API Gateway – User Service

**API Gateway – User Service** — это бэкенд-микросервис на **Go** с использованием фреймворка **Fiber**.  
Он реализует подход **Чистой архитектуры**, где приложение разделено на независимые слои —  
`handler → usecase → repository → storage`, что обеспечивает поддерживаемость, масштабируемость и лёгкое расширение.

---

## 📖 Описание

Сервис предоставляет простой **REST API для управления пользователями** и может использоваться как основа для более крупных бэкенд-систем.  

### ✨ Возможности
- **Создание** новых пользователей с автоматически генерируемым UUID  
- **Получение** пользователя по ID или списка всех пользователей  
- **Обновление** данных существующих пользователей  
- **Удаление** пользователей  
- **Управление конфигурацией** через `config.yaml`  
- **Структурированное логирование** с использованием `slog`  
- **Кастомная обработка ошибок** (`apperr.ErrNotFound`)  
- **Чистое разделение слоёв** для удобного масштабирования и расширения  

Планируемые улучшения:  
- [ ] Кэширование (Redis / in-memory)  
- [ ] CI/CD пайплайны (GitHub Actions + Docker)  
- [ ] Интеграционные и модульные тесты  
- [ ] Расширение модульности за пределы User-сервиса  

---

## 📂 Структура проекта

```bash
.
├── app.go             # Точка входа
├── config.go          # Загрузчик конфигурации
├── config.yaml        # Пример конфигурации
├── router.go          # Настройка HTTP-маршрутизатора
├── internal/          
│   ├── handler/       # HTTP-хендлеры (Fiber)
│   ├── usecase/       # Бизнес-логика (UserProvider)
│   ├── repository/    # Слой доступа к данным
│   ├── models/        # Доменные модели (User и др.)
│   └── apperr/        # Кастомные ошибки приложения
├── storage.go         # Слой хранения
├── cache.go           # (Планируется) кэширование
└── tests/             # (Планируется) интеграционные и модульные тесты


---


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


