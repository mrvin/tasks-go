# **Магазин мерча**

В Авито существует внутренний магазин мерча, где сотрудники могут приобретать товары за монеты (coin). Каждому новому сотруднику выделяется 1000 монет, которые можно использовать для покупки товаров. Кроме того, монеты можно передавать другим сотрудникам в знак благодарности или как подарок.

## Описание задачи

Необходимо реализовать сервис, который позволит сотрудникам обмениваться монетками и приобретать на них мерч. Каждый сотрудник должен иметь возможность видеть:

- Список купленных им мерчовых товаров  
- Сгруппированную информацию о перемещении монеток в его кошельке, включая:  
  - Кто ему передавал монетки и в каком количестве  
  - Кому сотрудник передавал монетки и в каком количестве

Количество монеток не может быть отрицательным, запрещено уходить в минус при операциях с монетками.

## **Общие вводные**

**Мерч** — это продукт, который можно купить за монетки. Всего в магазине доступно 10 видов мерча. Каждый товар имеет уникальное название и цену. Ниже приведён список наименований и их цены.

| Название     | Цена |
|--------------|------|
| t-shirt      | 80   |
| cup          | 20   |
| book         | 50   |
| pen          | 10   |
| powerbank    | 200  |
| hoody        | 300  |
| umbrella     | 200  |
| socks        | 10   |
| wallet       | 50   |
| pink-hoody   | 500  |

Предполагается, что в магазине бесконечный запас каждого вида мерча.


## **Условия**

* Используйте этот [API](https://github.com/avito-tech/tech-internship/blob/main/Tech%20Internships/Backend/Backend-trainee-assignment-winter-2025/schema.json) 
* Сотрудников может быть до 100к, RPS — 1k, SLI времени ответа — 50 мс, SLI успешности ответа — 99.99%   
* Для авторизации доступов должен использоваться JWT. Пользовательский токен доступа к API  выдается после авторизации/регистрации пользователя. При первой авторизации пользователь должен создаваться автоматически.
* Реализуйте покрытие бизнес сценариев юнит-тестами. Общее тестовое покрытие проекта должно превышать 40%
* Реализуйте интеграционный или E2E-тест на сценарий покупки мерча  
* Реализуйте интеграционный или E2E-тест на сценарий передачи монеток другим сотрудникам


## **Дополнительные задания**

Эти задания не являются обязательными, но выполнение всех или части из них даст вам преимущество перед другими кандидатами.  

* Провести нагрузочное тестирование полученного решения и приложить результаты тестирования 
* Реализовать интеграционное или E2E-тестирование для остальных сценариев  
* Описать конфигурацию линтера (.golangci.yaml в корне проекта для go, phpstan.neon для PHP или ориентируйтесь на свои, если используете другие ЯП для выполнения тестового)

## **Требования по стеку**

**Язык сервиса:** предпочтительным является Go, но также допустимы следующие языки: PHP, Java, Python, C#.
 
**База данных:** рекомендуется использовать PostgreSQL, но также допустимо использовать MySQL.
 
Для **деплоя зависимостей и самого сервиса** используйте Docker Compose. Порт доступа к сервису должен быть 8080 и быть доступен снаружи как `localhost:8080`.

## **Сборка и запуск приложения в Docker Compose**
```shell script
$ make run
...............
```

## **Пример использования HTTP API**
```shell script
$ curl -X 'POST' 'http://localhost:8080/api/auth' \
  -H 'Content-Type: application/json' \
  -d '{
	"username": "Bob",
  	"password": "qwerty"
}'
{
	"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk2MzcwODEsImlhdCI6MTczOTYzNjc4MSwidXNlcm5hbWUiOiJCb2IifQ.VTlSCG-dPg15S_F6AQYrMm6iUL4YGendt48UBmfY38s",
	"status": "OK"
}

$ curl -X 'POST' 'http://localhost:8080/api/auth' \
  -H 'Content-Type: application/json' \
  -d '{
	"username": "Alice",
  	"password": "qwerty"
}'
{
	"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk2MzcxMDIsImlhdCI6MTczOTYzNjgwMiwidXNlcm5hbWUiOiJBbGljZSJ9.8Q6tvQOzggdBstorK-2Y-cQVMDk7_tLNly33UmQkIME",
	"status":"OK"
}

$ curl -X 'POST' 'http://localhost:8080/api/sendCoin' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk2MzY2NjQsImlhdCI6MTczOTYzNjM2NCwidXNlcm5hbWUiOiJCb2IifQ.-k5e1zikQVW9H5GvjMv0lgG00aRNHRM8FcyBfKsJ7RY' \
  -H 'Content-Type: application/json' \
  -d '{
  "toUser": "Alice",
  "amount": 500
}'
{
	"status":"OK"
}

$ curl -X 'GET' 'http://localhost:8080/api/buy/pink-hoody' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk2MzcxMDIsImlhdCI6MTczOTYzNjgwMiwidXNlcm5hbWUiOiJBbGljZSJ9.8Q6tvQOzggdBstorK-2Y-cQVMDk7_tLNly33UmQkIME'
{
	"status":"OK"
}

$ curl -X 'GET' 'http://localhost:8080/api/info' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk2MzcwODEsImlhdCI6MTczOTYzNjc4MSwidXNlcm5hbWUiOiJCb2IifQ.VTlSCG-dPg15S_F6AQYrMm6iUL4YGendt48UBmfY38s'
{
	"coins": 500,
	"inventory": [],
	"coinHistory": {
		"received": [],
		"sent": [
			{
				"toUser": "Alice",
				"amount": 500
			}
		]
	},
	"status": "OK"
}

$ curl -X 'GET' 'http://localhost:8080/api/info' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk2MzcxMDIsImlhdCI6MTczOTYzNjgwMiwidXNlcm5hbWUiOiJBbGljZSJ9.8Q6tvQOzggdBstorK-2Y-cQVMDk7_tLNly33UmQkIME'
{
	"coins": 1000,
	"inventory": [
		{
			"type": "pink-hoody",
			"quantity": 1
		}
	],
	"coinHistory": {
		"received": [
			{
				"fromUser": "Bob",
				"amount": 500
			}
		],
		"sent": []
	},
  	"status": "OK"
}

```

### Полезные ссылки
- [Полезные материалы](https://github.com/avito-tech/tech-internship/tree/main/Tech%20Internships/Backend)
