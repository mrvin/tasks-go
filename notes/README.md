## Сервис заметки

Необходимо спроектировать и реализовать сервис, предоставляющий REST API интерфейс с методами:

#### Создание заметки
- Эндпоинт - POST /notes
- Параметры запроса:
   - JSON-объект в теле запроса с параметрами:
        - title – заголовок заметки
        - description – описание заметки
 - Статус ответа 201 если заметка создана успешно
 - 
#### Получение списка заметок
- Эндпоинт - GET /notes
- Параметры запроса отсутствуют

### Условия
 - Необходимо реализовать аутентификацию и авторизацию. Пользователи должны иметь доступ только к своим заметкам;
 - Для реализации сервиса использовать язык программирования Golang;
 - Сервис должен работать через REST API, для передачи данных использовать формат json;
 - Логирование событий в едином формат;
 - Запуск сервиса и требуемой ему инфраструктуры должен производиться в Docker контейнерах.

### Сборка и запуск приложения в Docker Compose

```shell script
$ make build
...............
$ make up
...............
```

### Пример использования http API
```bash
$ curl -i -X GET http://localhost:8088/health
$ curl -i -X POST 'http://localhost:8088/users' -H "Content-Type: application/json" -d '{"user_name":"Bob","password":"qwerty"}'
$ curl -i -X POST 'http://localhost:8088/users' -H "Content-Type: application/json" -d '{"user_name":"Alice","password":"password123"}'
$ curl --user Bob:qwerty -i -X POST 'http://localhost:8088/notes' -H "Content-Type: application/json" -d '{"title":"title 1","description":"description 1"}'
$ curl --user Bob:qwerty -i -X POST 'http://localhost:8088/notes' -H "Content-Type: application/json" -d '{"title":"title 2","description":"description 2"}'
$ curl --user Alice:password123 -i -X POST 'http://localhost:8088/notes' -H "Content-Type: application/json" -d '{"title":"title 3","description":"description 3"}'
$ curl --user Bob:qwerty -i -X GET 'http://localhost:8088/notes'
$ curl --user Alice:password123 -i -X GET 'http://localhost:8088/notes'
```
