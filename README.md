
# OZON - URL - SHORTENER
## Структура проекта

```
url-shortener-ozon/
├── api/
│   ├── grpcserver/
│   │   └── server.go                   # Реализация gRPC сервера
│   └── restserver/
│       └── server.go                   # Реализация REST сервера
├── cmd/
│   └── url-shortener/
│       └── main.go                     # Точка входа в приложение; выбор хранилища
├── internal/
│   ├── config/
│   │   └── config.go                   # Конфигурация приложения, загрузка переменных окружения
│   ├── logger/
│   │   └── logger.go                   # Логгирование (запись в файл и консоль)
│   ├── service/
│   │   └── service.go                  # Логика сокращения URL
│   └── storage/
│       ├── models/
│       │   └── models.go               # Модели данных (URL, DeviceURL)
│       ├── in_memory_storage.go        # Реализация in-memory хранилища
│       └── postgres_storage.go         # Реализация хранилища на PostgreSQL
├── test/
│   └── service_test.go                 # Unit-тесты для сервиса
├── .env                                # Файл конфигурации окружения
├── .gitignore                          # Git ignore (игнорирует .env, app.log и т.п.)
├── Dockerfile                          # Dockerfile для сборки Docker-образа
├── go.mod                              # Go-модуль и зависимости
├── go.sum                              # Контрольные суммы зависимостей
└── README.md                           # Описание проекта и инструкция по запуску
```

## Описание проекта

Этот проект представляет собой сервис для сокращения URL. Он позволяет пользователям сокращать длинные URL и перенаправлять на оригинальные URL, используя сокращенную версию.

### Компоненты

- **api/grpcserver/server.go**: Реализация gRPC сервера.
- **api/restserver/server.go**: Реализация REST сервера.
- **cmd/url-shortener/main.go**: Точка входа в приложение.
- **internal/service/service.go**: Логика сокращения URL.
- **internal/storage/postgres_storage.go**: Реализация хранилища на PostgreSQL.
- **internal/storage/in_memory_storage.go**: Реализация in-memory хранилища.
- **test/service_test.go**: Unit-тесты для сервиса.

## Как работает?

1. **Сокращение URL**: Пользователи отправляют POST-запрос с длинным URL в сервис. Сервис генерирует короткий URL и сохраняет соответствие в репозитории.
2. **Перенаправление**: Пользователи отправляют GET-запрос с коротким URL. Сервис ищет оригинальный URL в репозитории и перенаправляет пользователя на него.

## Инструкции по запуску

### Запуск через Docker

1. **Клонируйте репозиторий**:
    ```sh
    git clone https://github.com/yourusername/url-shortener-ozon.git
    cd url-shortener-ozon
    ```

2. **Соберите Docker-образ**:
    ```sh
    docker build -t url-shortener .
    ```

3. **Запустите контейнер**:
    ```sh
    docker run -p 8080:8080 -p 50051:50051 --env-file .env url-shortener
    ```

### Переменные окружения

Создайте файл [.env](http://_vscodecontentref_/4) в корне проекта со следующим содержимым:
```properties
PORT=8080
BASE_URL=http://localhost:8080
DATABASE_URL=postgres://user:password@localhost:5432/url_shortener?sslmode=disable
STORAGE_MODE=postgres
SECRET_KEY=your_secret_key_here
