# 🔐 Auth Service — MPB Blog Platform

Authentication and authorization module for the MPB Blog Platform, built with **Go + Fiber + PostgreSQL + JWT**.  
Implements secure user registration, login, and JWT-based authentication.

---

## 🚀 Features

- User registration (`POST /api/auth/register`)
- User login with JWT token generation (`POST /api/auth/login`)
- JWT authentication middleware
- Input validation via generics + go-playground/validator
- Protected routes example
- Integrated Swagger documentation

---

## ⚙️ Tech Stack

| Component       | Library |
|-----------------|----------|
| Web Framework   | [Fiber v2](https://github.com/gofiber/fiber) |
| DB Access | [sqlx](https://github.com/jmoiron/sqlx) |
| Database        | PostgreSQL |
| JWT Auth        | [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt) |
| Validation      | [go-playground/validator/v10](https://github.com/go-playground/validator) |
| Documentation   | [fiber-swagger](https://github.com/swaggo/fiber-swagger) |

---
