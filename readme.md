## Ads Service / Cервис объявлений

Что есть в проекте:
- Чистая архитектура:
  - Entities: ads, user
  - Use cases: app
  - Interfaces / adapters: adrepo, userrepo
  - Infrastructure: httpgin, grpc
- Фреймворк: Gin,
- APIs: REST, gRPC (by protoc-gen-go-grpc),
- Tests: unit, integration, fuzz, benchmark, mock (by mockery v2.20.0), coverage проекта более 80%,
- Graceful shutdown,
- Подключен собственный модуль валидации данных: https://github.com/AlexeyNikitin01/validate/tree/v1.2.3
- Добавлен docker, docker-compose
- Добавлена БД: postgres
