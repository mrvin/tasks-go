## Лучший продукт

Есть 2 файла с данными о продуктах (наименование, цена, рейтинг) в 2-х
форматах - CSV и JSON. Необходимо написать программу, которая считывает
данные из переданного в параметре файла, и выводит  «самый дорогой
продукт» и «с самым высоким рейтингом». Предусмотреть, что файлы могут
быть огромными. Репозиторий должен содержать Dockerfile для сборки
готового приложения в docker среде.

#### Сборка и запуск приложения в Docker

```shell script
$ docker build -t top-product .
$ docker run --rm --name top-product top-product

/topProduct -f testdata/db.json
Top Price:
Name     Price  Rating
----     -----  ------
Варенье  200    5

Top Rating:
Name     Price  Rating
----     -----  ------
Варенье  200    5
/topProduct -f testdata/db.csv
Top Price:
Name     Price  Rating
----     -----  ------
Печенье  3      5

Top Rating:
Name     Price  Rating
----     -----  ------
Печенье  3      5
```
