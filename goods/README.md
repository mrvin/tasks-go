## Сервис "Товары"

## Требования
- При добавлении, редактировании или удалении записи в PostgreSQL писать события в
	ClickHouse через очередь NATS. События писать пачками.

Реализовать REST API-сервис для хранения и управления информацией о товарах. Cервис
должен реализовать 6 методов и их логику:

#### Проверка работоспособности
- Эндпоинт - GET /health
- Статус ответа 200 если сервис работает исправно

##### Пример
```bash
$ curl -i -X GET 'http://localhost:8080/health'

{
  "status": "OK"
}
```

#### Создание товара
При добавлении товара в таблицу устанавливать приоритет как максимальный приоритет в
таблице +1. Приоритеты начинаются с 1.

##### Пример
```bash
$ curl -i -X POST 'http://localhost:8080/good/create?projectID=1' \
-H "Content-Type: application/json" \
-d '{
	"name": "Sample Product"
}'

{
  "id": 1,
  "projectID": 1,
  "name": "Sample Product",
  "priority": 1,
  "removed": false,
  "createdAt": "2025-06-25T14:42:09.557392Z"
}
```

#### Обновление товара

##### Пример
```bash
$ curl -i -X PATCH 'http://localhost:8080/good/update?id=1&projectID=1' \
-H "Content-Type: application/json" \
-d '{
	"name": "New name",
	"description": "New description"
}'

{
  "id": 1,
  "projectID": 1,
  "name": "New name",
  "description": "New description",
  "priority": 1,
  "removed": true,
  "createdAt": "2025-06-25T14:05:21.68716Z"
}
```

#### Удаление товара

##### Пример
```bash
$ curl -i -X DELETE 'http://localhost:8080/good/remove?id=1&projectID=1'

{
  "id": 1,
  "projectID": 1,
  "removed": true,
  "status": "OK"
}
```

#### Получение списка всех товаров

##### Пример
```bash
$ curl -i -X GET 'http://localhost:8080/goods/list?limit=10&offset=0'

{
  "meta": {
    "total": 1,
    "removed": 1,
    "limit": 10,
    "offset": 0
  },
  "goods": [
    {
      "id": 1,
      "projectID": 1,
      "name": "Sample Product 2",
      "description": "Description 2",
      "priority": 1,
      "removed": true,
      "createdAt": "2025-06-25T14:05:21.68716Z"
    }
  ],
  "status": "OK"
}
```

#### Обновление приоритета товара
##### Пример
```bash
$ curl -i -X PATCH 'http://localhost:8080/good/reprioritize?id=2&projectID=1' \
-H "Content-Type: application/json" \
-d '{
	"newPriority": 4
}'

{
  "priorities": [
    {
      "id": 3,
      "priority": 2
    },
    {
      "id": 4,
      "priority": 3
    },
    {
      "id": 2,
      "priority": 4
    }
  ]
}
```
### Сборка и запуск приложения в Docker Compose

```bash
$ make run
...............
```

## Структура проекта
```bash
$ tree .
.
├── cmd
│   └── goods-server
│       ├── Dockerfile
│       ├── main.go
│       └── Makefile
├── configs
│   └── config.env
├── deployments
│   └── docker-compose.yaml
├── go.mod
├── go.sum
├── internal
│   ├── config
│   │   └── config.go
│   ├── httpserver
│   │   ├── handlers
│   │   │   ├── good
│   │   │   │   ├── create
│   │   │   │   │   └── create.go
│   │   │   │   ├── delete
│   │   │   │   │   └── delete.go
│   │   │   │   ├── list
│   │   │   │   │   └── list.go
│   │   │   │   ├── reprioritize
│   │   │   │   │   └── reprioritize.go
│   │   │   │   └── update
│   │   │   │       └── update.go
│   │   │   └── health
│   │   │       └── health.go
│   │   └── server.go
│   ├── logger
│   │   └── logger.go
│   ├── queue
│   │   └── nats
│   │       └── mq.go
│   └── storage
│       ├── clickhouse
│       │   └── storage.go
│       ├── sql
│       │   ├── create.go
│       │   ├── delete.go
│       │   ├── list.go
│       │   ├── reprioritize.go
│       │   ├── storage.go
│       │   └── update.go
│       └── storage.go
├── Makefile
├── migrations
│   ├── clickhouse
│   │   ├── 000001_init_schema.down.sql
│   │   └── 000001_init_schema.up.sql
│   └── postgres
│       ├── 000001_init_schema.down.sql
│       └── 000001_init_schema.up.sql
├── pkg
│   ├── http
│   │   ├── logger
│   │   │   └── logger.go
│   │   └── response
│   │       └── response.go
│   └── retry
│       └── retry.go
└── README.md
```

### Написать SQL-запросы для ClickHouse:
#### Выборки всех уникальных Description у которых более 2 событий
```sql
SELECT Description,
    count() AS descriptions_count
FROM goods_events
GROUP BY Description
HAVING count() > 2;
```

#### Выборки ProjectID которые совершили более 2 различных Description
```sql
SELECT ProjectID,
	uniqExact(Description) AS unique_descriptions_count
FROM goods_events
GROUP BY ProjectID
HAVING unique_descriptions_count > 2
ORDER BY ProjectID;
```

#### Выборки событий которые произошли в первый день каждого месяца
```sql
SELECT *
FROM goods_events
WHERE toDayOfMonth(EventTime) = 1
ORDER BY EventTime DESC;
```

#### Вывод событий по заданному Description и временному диапазону
```sql
SELECT *
FROM goods_events
WHERE 
    Description = 'Create new good'
    AND EventTime BETWEEN '2025-06-26 10:30:00' AND '2025-06-26 22:30:00'
ORDER BY EventTime DESC;
```
