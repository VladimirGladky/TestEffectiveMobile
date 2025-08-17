# TestEffectiveMobile
## О проекте
REST-сервис для агрегации данных об онлайн-подписках пользователей.

### 🌐 Доступные API эндпоинты

| Метод | Путь | Описание |
| :--- | :--- | :--- |
| POST | /api/v1/create | Создание подписки |
| GET | /api/v1/list/{user_id} | Получение списка подписок пользователя |
| GET | /api/v1/sum | Расчет суммы подписок |
| GET | /api/v1/read/{id} | Получение подписки по id |
| PUT | /api/v1/update/{id} | Обновление подписки по id |
| DELETE | /api/v1/delete/{id} | Удаление подписки по id |

## 🗄️ База данных

В качестве базы данных используется **PostgreSQL**.

### 🛠️ Миграции

Для создания и управления схемой базы данных применяются миграции, которые находятся в папке [`migrations/`](./migrations).

### 🗃️ Структура базы данных

В базе данных предусмотрена одна таблица `subscriptions`:

| Поле | Тип | Описание |
| :--- | :--- | :--- |
| id | VARCHAR(255) | Уникальный идентификатор подписки |
| service_name | VARCHAR(255) | Название сервиса |
| price | INT | Цена подписки |
| user_id | VARCHAR(255) | Идентификатор пользователя |
| start_date | DATE | Дата начала подписки |
| end_date | DATE | Дата окончания подписки |

## Технологии и библиотеки
Вычислитель написан на языке **Go** и использует следующие библиотеки и инструменты:

#### Язык программирования:
- **Go** (версия 1.24.0)

## 📚 Используемые библиотеки

В проекте используются следующие ключевые Go-библиотеки:

### 🔑 Основные утилиты
| Библиотека | Назначение | Документация |
|------------|------------|--------------|
| `github.com/google/uuid` | Генерация и парсинг UUID (уникальных идентификаторов) | [ссылка](https://github.com/google/uuid) |
| `github.com/gorilla/mux` | Маршрутизация HTTP-запросов с поддержкой переменных и middleware | [ссылка](https://github.com/gorilla/mux) |
| `github.com/joho/godotenv` | Загрузка конфигурации из `.env` файлов в переменные окружения | [ссылка](https://github.com/joho/godotenv) |
| `github.com/ilyakaznacheev/cleanenv` | Чтение и валидация конфигурации из окружения и файлов | [ссылка](https://github.com/ilyakaznacheev/cleanenv) |

### 🗃️ Работа с данными
| Библиотека | Назначение | Документация |
|------------|------------|--------------|
| `github.com/golang-migrate/migrate/v4` | Управление миграциями базы данных (создание, применение, откат) | [ссылка](https://github.com/golang-migrate/migrate) |
| `github.com/jackc/pgx/v5` | Высокопроизводительный драйвер PostgreSQL и инструменты работы с ним | [ссылка](https://github.com/jackc/pgx) |

### 📝 Логирование
| Библиотека | Назначение | Документация |
|------------|------------|--------------|
| `go.uber.org/zap` | Быстрое структурированное логирование с минимальным оверхедом | [ссылка](https://go.uber.org/zap) |

## 📚 Структура проекта

```bash
├── cmd/ #Основной исполняемый файл проекта
├── config/ #Конфигурационные файлы
├── docs/ # Swagger документация
├── internal/ # Внутренняя бизнес-логика (не предназначена для внешнего использования)
│   ├── app/ # Инициализация приложения 
│   ├── config/ # Конфигурация приложения
│   ├── models/ # Модели данных
│   ├── repository/ # Слой взаимодействия с базой данных
│   ├── service/ # Слой бизнес-логики
│   ├── transport/ # Слой взаимодействия с пользователем
├── migrations/ # Скрипты миграций базы данных
├── pkg/ # Публичные пакеты, доступные извне (reusable)
│   ├── logger/ # Пакет логирования
│   ├── postgres/ # Пакет работы с базой данных PostgreSQL
```

## 📚 Запуск проекта

Для запуска проекта необходимо выполнить следующие шаги:

```
git clone https://github.com/VladimirGladky/TestEffectiveMobile.git
```

```
cd TestEffectiveMobile
```

Теперь вы можете запустить проект , но для этогт нужно чтобы был установлен Go версии 1.24.0
Ссылка для скачивания: [Go Download](https://go.dev/doc/install)

Перед запуском агента и оркестратора , воспользуйтесь командой

```bash
go mod download

```

### Запуск сервера

```bash
docker compose up
```

Сервер будет доступен по адресу http://localhost:4047/api/v1

Чтобы остановить сервер, выполните команду:

```bash
docker compose down
```

## Примеры использования со всеми возможными сценариями

1. Создание подписки

```bash
curl -X POST -H "Content-Type: application/json" -d '{"service_name":"Netflix","price":100,"user_id":"user123","start_date":"06-2025","end_date":"11-2025"}' http://localhost:4047/api/v1/create
```

В ответ мы получаем id созданной подписки:

```bash
{"id":"1bccafaa-9a3a-4f0c-9223-c1863c2c2b5b"}
```

2. Получение подписки по id 

```bash
curl -X GET http://localhost:4047/api/v1/read/1bccafaa-9a3a-4f0c-9223-c1863c2c2b5b
```

В ответ мы получаем подписку:

```bash
{
    "service_name":"Netflix",
    "price":100,
    "id":"1bccafaa-9a3a-4f0c-9223-c1863c2c2b5b",
    "user_id":"user123",
    "start_date":"06-2025",
    "end_date":"11-2025"
}
```

3. Обновление подписки по id 

```bash
curl -X PUT -H "Content-Type: application/json" -d '{"service_name":"Netflix3","price":500,"start_date":"08-2025","end_date":"12-2025"}' http://localhost:4047/api/v1/update/1bccafaa-9a3a-4f0c-9223-c1863c2c2b5b
```

В ответ мы получаем , что подписка была обновлена:

```bash
{"message":"Updated"}
```

4. Удаление подписки по id 

```bash
curl -X DELETE http://localhost:4047/api/v1/delete/1bccafaa-9a3a-4f0c-9223-c1863c2c2b5b
```

В ответ мы получаем , что подписка была удалена:

```bash
{"message":"Deleted"}
```

5. Получение всех подписок пользователя по id 

```bash
curl -X GET http://localhost:4047/api/v1/list/user123
```

В ответ мы получаем все подписки пользователя:

```bash
{
    "subscriptions":
        [
            {
              "service_name":"Netflix","price":100,"id":"d0b9c37b-6566-4792-ab98-71e9b5fd251e",
              "user_id":"user123","start_date":"06-2025","end_date":"11-2025"
            },
            {
              "service_name":"Netflix","price":100,"id":"3e33a070-4545-465f-b5f3-8ee1c216f921",
              "user_id":"user123","start_date":"07-2025","end_date":"12-2025"
            }
        ]
}
```

6. Расчет суммы подписок , тут можно выставлять query параметры start_date , end_date , user_id , service_name
    параметры start_date и end_date должны быть в формате MM-YYYY - например 06-2025

```bash
curl -X GET http://localhost:4047/api/v1/sum?start_date=05-2025&end_date=12-2025&service_name=Netflix&user_id=user123
```

В ответ мы получаем сумму подписок:

```bash
{"sum":200}
```

## Swagger документация

После запуска сервера можно посмотреть документацию в браузере по адресу http://localhost:4047/swagger/index.html



