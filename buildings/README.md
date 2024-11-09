## Сервис для хранения информации о строениях

Необходимо спроектировать и реализовать сервис, предоставляющий REST API интерфейс с методами:

#### Создание строения
- Эндпоинт - POST /buildings
- Параметры запроса:
   - JSON-объект в теле запроса с параметрами:
        - name – название строения
        - city – город
        - year - год сдачи
        - number_floors - кол-во этажей
 - Статус ответа 201 если пользователь создан успешно

#### Получение списка строений
- Эндпоинт - GET /buildings
- Параметры запроса:
    - city – город (необязательный параметр)
    - year - год сдачи (необязательный параметр)
    - number_floors - кол-во этажей (необязательный параметр)
- Ответ должен содержать в теле массив JSON-объектов с информацией о строениях. Каждый объект содержит параметры:
    - name – название строения
    - city – город
    - year - год сдачи
    - number_floors - кол-во этажей
- Статус ответа 200 если список получен успешно

#### Получение документации
- Эндпоинт - GET /swagger/index.html.

### Условия
- Для реализации сервиса использовать язык программирования Golang;
- Сервис должен работать через REST API, для передачи данных использовать формат json;
- Для реализации REST API использовать веб-фреймворк [Gin](https://gin-gonic.com/);
- Данные необходимо хранить в PostgreSQL;
- Настройки сервиса читать из конфигурационного файла;
- Документацию генерировать из OpenApi файла при помощи [swag](https://github.com/swaggo/swag).

### Пример использования http API
```bash
$ curl -i -X GET http://localhost:8081/health
$ curl -i -X POST 'http://localhost:8081/buildings' -H "Content-Type: application/json" -d '{"name":"Building #1","city":"Saint Petersburg","year":2022,"number_floors":22}'
$ curl -i -X GET http://localhost:8081/buildings
$ curl -i -X GET 'http://localhost:8081/buildings?year=2022&city=Saint+Petersburg&number_floors=22'
```