## Сервис записная книга

Реализовать сервис, который будет получать по API ФИО, из открытых API обогащать ответ наиболее вероятными возрастом, полом и национальностью и сохранять данные в БД. По запросу выдавать информацию о найденных людях. Сервер должен реализовать 5 метода и их логику:

#### Создание персоны
 - Эндпоинт: POST /person
 - Параметры запроса:
    - JSON-объект в теле запроса с параметрами:
        - name - имя персоны
        - surname - фамилия персоны
        - patronymic - отчество персоны (необязательно)
 - Ответ содержит JSON-объект:
    - id – ID персоны. Генерируется сервером
 - Корректное сообщение обогатить наиболее вероятными:
    - Возрастом - [api.agify.io](https://api.agify.io/?name=Dmitriy) 
    - Полом - [api.genderize.io](https://api.genderize.io/?name=Dmitriy)
    - Национальностью - [api.nationalize.io](https://api.nationalize.io/?name=Dmitriy)
    
#### Получение информации о персоне
 - Эндпоинт - GET /person/?id=1
 - Параметры запроса:
    - id – ID персоны, указан в пути запроса
 - Ответ содержит JSON-объект:
    - id - ID персоны
    - name - имя персоны
    - surname - фамилия персоны
    - patronymic - отчество персоны (необязательно)
    - age - возраст
    - gender - пол
    - CountryID - код страны (национальность)
    
#### Обновление персоны
 - Эндпоинт – PUT /person/?id=1
 - Параметры запроса:
    - id – ID персоны, указан в пути запроса
    - JSON-объект в теле запроса с параметрами:
        - name - имя персоны
        - surname - фамилия персоны
        - patronymic - отчество персоны (необязательно)
        - age - возраст
        - gender - пол
        - CountryID - код страны (национальность)
#### Удаление персоны
 - Эндпоинт – DELETE /person/?id=1
 - Параметры запроса:
    - id – ID персоны, указан в пути запроса
    
#### Получение списка всех персон
 - Эндпоинт – GET /list-persons
 - Ответ должен содержать в теле массив JSON-объектов с всеми персонами. Каждый объект содержит параметры:
    - id - ID персоны
    - name - имя персоны
    - surname - фамилия персоны
    - patronymic - отчество персоны (необязательно)
    - age - возраст
    - gender - пол
    - CountryID - код страны (национальность)

### Сборка и запуск приложения в Docker Compose

```shell script
$ make build
...............
$ make up
...............
```
### Пример использования http API
```bash
$ curl -i -X POST http://localhost:8088/person -H "Content-Type: application/json" -d '{"name": "Dmitriy","surname": "Ushakov","patronymic": "Vasilevich"}'
{"id":1,"status":"OK"}
$ curl -i -X POST http://localhost:8088/person -H "Content-Type: application/json" -d '{"name": "Vladimir","surname": "Vinogradov"}'
{"id":2,"status":"OK"}
$ curl -i -X GET http://localhost:8088/person/?id=1
{"id":1,"name":"Dmitriy","surname":"Ushakov","patronymic":"Vasilevich","age":43,"gender":"male","countryID":"UA"}
$ curl -i -X GET http://localhost:8088/person/?id=2
{"id":2,"name":"Vladimir","surname":"Vinogradov","age":64,"gender":"male","countryID":"RS"}
$ curl -i -X PUT http://localhost:8088/person/?id=2 -H "Content-Type: application/json" -d '{"name":"Vladimir","surname":"Vinogradov","age":28,"gender":"male","countryID":"RU"}'
$ curl -i -X GET http://localhost:8088/list-persons
{"persons":[{"id":1,"name":"Dmitriy","surname":"Ushakov","patronymic":"Vasilevich","age":43,"gender":"male","countryID":"UA"},{"id":2,"name":"Vladimir","surname":"Vinogradov","age":28,"gender":"male","countryID":"RU"}],"status":"OK"}
$ curl -i -X DELETE http://localhost:8088/person/?id=1
{"status":"OK"}
```
