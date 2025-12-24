# AnonChat 💬

**AnonChat** — полнофункциональная система анонимного чата с микросервисной архитектурой, написанная на **Golang** и **React**.

Она обеспечивает **регистрацию и авторизацию пользователей**, **создание и управление комнатами чата**, а также **обмен сообщениями в реальном времени** через WebSocket.

Использует **микросервисную архитектуру** с API Gateway, **PostgreSQL** для хранения данных и **React** для веб-интерфейса.

---

**AnonChat** is a full-featured anonymous chat system with microservices architecture, built with **Golang** and **React**.

It supports **user registration and authentication**, **room creation and management**, and **real-time messaging** via WebSocket.

Built with **microservices architecture** using API Gateway, **PostgreSQL** for data storage, and **React** for web interface.

---

## 🧩 Core Features / Основные возможности

| Feature / Функция              | Description (EN)                                          | Описание (RU)                                              |
| ------------------------------- | --------------------------------------------------------- | ---------------------------------------------------------- |
| **User Registration**           | Create account with login, email, and password             | Регистрация пользователя с логином, email и паролем       |
| **User Authentication**         | Login and receive JWT access/refresh tokens               | Вход и получение JWT токенов (access/refresh)             |
| **Token Management**            | Automatic token refresh and validation                    | Автоматическое обновление и валидация токенов              |
| **Room Creation**               | Create chat rooms with custom name and user limit         | Создание комнат чата с именем и лимитом пользователей      |
| **Room Joining**                | Join existing rooms by ID                                  | Присоединение к существующим комнатам по ID               |
| **Real-time Messaging**        | WebSocket-based instant messaging                         | Мгновенный обмен сообщениями через WebSocket              |
| **Room Management**            | Admin can remove created rooms                            | Администратор может удалять созданные комнаты              |
| **Rate Limiting**               | Protection against spam and abuse                          | Защита от спама и злоупотреблений                          |
| **Structured Logging**         | Request logging with custom logger                         | Логирование запросов через собственный логгер             |
| **Auto Token Refresh**          | Automatic access token refresh on frontend                | Автоматическое обновление токенов на фронтенде             |

---

## ⚙️ Technologies / Технологии

### Backend / Бэкенд

| Layer              | Stack                                      |
| ------------------ | ------------------------------------------ |
| **Language**       | Go 1.25                                    |
| **Framework**      | Gin                                        |
| **WebSocket**      | Gorilla WebSocket                          |
| **Database**       | PostgreSQL 17                              |
| **Auth**           | JWT (access / refresh tokens)              |
| **Middleware**     | RateLimiter, Logger, Auth                  |
| **Config**         | YAML configuration files                   |
| **Architecture**   | Microservices (Gateway + Auth Service)     |

### Frontend / Фронтенд

| Layer              | Stack                                      |
| ------------------ | ------------------------------------------ |
| **Framework**      | React 18                                   |
| **Build Tool**     | Vite                                       |
| **Routing**        | React Router DOM                           |
| **WebSocket**      | Native WebSocket API                       |
| **Styling**        | CSS Modules                                |

---

## 🚀 API Endpoints / Эндпойнты

### Authentication / Авторизация

| Method | Path                    | Description (EN)                         | Описание (RU)                      |
| ------ | ----------------------- | ---------------------------------------- | ----------------------------------- |
| `GET`  | `/api/v1/status`        | Check service health                     | Проверка статуса сервиса            |
| `POST` | `/api/v1/auth/register` | Register new user                        | Регистрация нового пользователя     |
| `POST` | `/api/v1/auth/login`    | Authenticate user and receive tokens     | Авторизация и получение токенов     |
| `POST` | `/api/v1/auth/refresh`  | Refresh access token                     | Обновление access токена             |
| `POST` | `/api/v1/auth/logout`   | Logout and revoke session                | Выход и завершение сессии            |
| `GET`  | `/api/v1/auth/validate`| Validate access token                    | Валидация access токена             |

### Rooms / Комнаты

| Method | Path                    | Description (EN)                         | Описание (RU)                      |
| ------ | ----------------------- | ---------------------------------------- | ----------------------------------- |
| `GET`  | `/api/v1/rooms/ws`      | WebSocket connection for chat            | WebSocket подключение для чата       |
| `POST` | `/api/v1/rooms/create`  | Create new chat room                     | Создание новой комнаты чата         |
| `DELETE` | `/api/v1/rooms/remove`  | Remove chat room (admin only)            | Удаление комнаты (только админ)      |

---

## 🔒 Authentication Flow / Процесс авторизации

1. **Register** — пользователь отправляет `login`, `email`, `password`.
   → получает `access_token` (15 мин) и `refresh_token` (7 дней).

2. **Login** — вводит `email` и `password`.
   → получает пару токенов.

3. **API Requests** — все запросы (кроме register/login) требуют `Authorization: Bearer <access_token>`.

4. **Token Refresh** — при истечении access токена:
   - Frontend автоматически обновляет токен через `/auth/refresh`
   - Получает новую пару токенов
   - Повторяет оригинальный запрос

5. **WebSocket** — подключение требует `room` ID и `token` (access token) в query параметрах.

6. **Logout** — удаляет refresh токен из БД, завершая сессию.

---

## 💬 Chat Flow / Процесс чата

1. **Create Room** — авторизованный пользователь создает комнату с именем и лимитом пользователей.
   → получает уникальный `room_id`.

2. **Join Room** — пользователь может:
   - Присоединиться к созданной комнате из списка
   - Ввести ID комнаты для подключения

3. **WebSocket Connection** — устанавливается соединение:
   ```
   ws://host:port/api/v1/rooms/ws?room=<room_id>&token=<access_token>
   ```

4. **Messaging** — сообщения отправляются через WebSocket и транслируются всем участникам комнаты в реальном времени.

5. **Room Management** — администратор комнаты может удалить её, что закроет все активные соединения.

---

## ⚡ Quick Start / Быстрый старт

### Prerequisites / Требования

- **Go 1.25+**
- **Node.js 18+** и **npm**
- **PostgreSQL 17** (или Docker)
- **Docker** и **Docker Compose** (опционально)

### Backend / Бэкенд

#### Local Development / Локальная разработка

1. **Start PostgreSQL** (если не используете Docker):
   ```bash
   # Настройте PostgreSQL и создайте базу данных
   ```

2. **Start Authorization Service**:
   ```bash
   cd chat-server/services/authorization
   go run cmd/server.go
   ```

3. **Start Gateway**:
   ```bash
   cd chat-server/gateway
   go run cmd/server.go
   ```

#### Docker / Докер

```bash
cd chat-server
docker-compose up --build
```

**Default URLs:**
- Gateway: `http://localhost:8080`
- Authorization Service: `http://authorization:80` (internal)

### Frontend / Фронтенд

```bash
cd chat-client-web
npm install
npm run dev
```

**Default URL:** `http://localhost:3000`

### Full Stack / Полный стек

```bash
# Terminal 1: Backend
cd chat-server
docker-compose up

# Terminal 2: Frontend
cd chat-client-web
npm install
npm run dev
```

---

## 🔧 Configuration / Конфигурация

### Gateway Config / Конфигурация Gateway

`chat-server/gateway/config/config.yaml`:

```yaml
server:
  addr: 0.0.0.0:80
  timeouts:
    write: 10s
    read: 10s
  rate-limits:
    max-requests: 100
    update-in: 1m
    mode: on
  logger:
    level: info
    mode: on
  auth:
    mode: on

authorization-sercice:
  addr: authorization:80
  mode: on
```

### Authorization Service Config / Конфигурация Auth Service

`chat-server/services/authorization/config/config.yaml`:

Настройте подключение к PostgreSQL и параметры JWT токенов.

---

## 🐳 Docker / Докер

Проект включает `docker-compose.yaml` для запуска всех сервисов:

```yaml
services:
  gateway:        # API Gateway (Port 8080)
  authorization:  # Auth Service
  postgres:       # PostgreSQL Database
```

**Запуск:**
```bash
docker-compose up --build
```

**Остановка:**
```bash
docker-compose down
```

---

## 📝 Development / Разработка

### Backend / Бэкенд

- **Gateway** обрабатывает все HTTP/WebSocket запросы
- **Authorization Service** управляет пользователями и сессиями
- **WebSocket Hub** управляет комнатами и клиентами в памяти

### Frontend / Фронтенд

- **React Router** для навигации
- **API Client** с автоматическим обновлением токенов
- **WebSocket Client** для реального времени
- **Responsive UI** с темной темой

---

## 🔐 Security / Безопасность

- **JWT Tokens** для аутентификации
- **Rate Limiting** для защиты от спама
- **Token Refresh** для безопасного обновления сессий
- **Input Validation** на сервере

---

## 📄 License / Лицензия

**MIT License** — свободное использование и модификация проекта.

---

## 🤝 Contributing / Вклад в проект

Приветствуются pull requests и issues! Пожалуйста, убедитесь, что код соответствует стилю проекта.

---

## 📧 Support / Поддержка

Для вопросов и предложений создавайте issues в репозитории проекта.



