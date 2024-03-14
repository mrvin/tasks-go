## Сервис "фотогалерея"

Реализовать HTTP API фотогалереи

### Возможности API: 
 - Загрузка фото.
 	- Сохранять данные о новых файлах в БД
	- Генерировать preview
 - Просмотр списка фото
 - Удаление фото

### Дополнительные требования:
 - Формат ответа - json
 - В качестве БД использовать sqlite

### Сборка и запуск приложения
```shell script
$ cd cmd/photo-gallery-server/
$ make build
go build -o ../../bin/photo-gallery-server -ldflags '-w -s'
$ cd ../../bin
$ ./photo-gallery-server -config ../configs/photo-gallery.yml
```

### Пример использования HTTP API
```shell script
$ curl -i -X POST 'http://localhost:8088/api/v1/photo?name=Poster_The_Red_Turtle_2016.jpg' --data-binary @../testdata/Poster_The_Red_Turtle_2016.jpg -H "Content-Type: image/jpeg"
HTTP/1.1 201 Created
Content-Type: application/json
Date: Thu, 14 Mar 2024 12:16:49 GMT
Content-Length: 15

{"status":"OK"}

$ curl -i -X POST 'http://localhost:8088/api/v1/photo?name=Poster_An_Interview_with_God_2018.jpg' --data-binary @../testdata/Poster_An_Interview_with_God_2018.jpg -H "Content-Type: image/jpeg"
HTTP/1.1 201 Created
Content-Type: application/json
Date: Thu, 14 Mar 2024 12:17:01 GMT
Content-Length: 15

{"status":"OK"}

$ curl -i -X GET 'http://localhost:8088/api/v1/listphotos'
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 14 Mar 2024 12:21:13 GMT
Content-Length: 471

{"listPhotos":[{"name":"Poster_The_Red_Turtle_2016.jpg","urlPhoto":"http://localhost:8088/api/v1/photo/Poster_The_Red_Turtle_2016.jpg","urlThumbnail":"http://localhost:8088/api/v1/photo/Poster_The_Red_Turtle_2016.thumb.jpg"},{"name":"Poster_An_Interview_with_God_2018.jpg","urlPhoto":"http://localhost:8088/api/v1/photo/Poster_An_Interview_with_God_2018.jpg","urlThumbnail":"http://localhost:8088/api/v1/photo/Poster_An_Interview_with_God_2018.thumb.jpg"}],"status":"OK"}

$ curl 'http://localhost:8088/api/v1/photo/Poster_An_Interview_with_God_2018.jpg' > Poster_An_Interview_with_God_2018.jpg
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  213k  100  213k    0     0  10.5M      0 --:--:-- --:--:-- --:--:-- 10.9M

$ curl 'http://localhost:8088/api/v1/photo/Poster_An_Interview_with_God_2018.thumb.jpg' > Poster_An_Interview_with_God_2018.thumb.jpg
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  4998  100  4998    0     0  3420k      0 --:--:-- --:--:-- --:--:-- 4880k

$ curl 'http://localhost:8088/api/v1/photo/Poster_The_Red_Turtle_2016.jpg' > Poster_The_Red_Turtle_2016.jpg
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  149k  100  149k    0     0  73.9M      0 --:--:-- --:--:-- --:--:--  146M

$ curl 'http://localhost:8088/api/v1/photo/Poster_The_Red_Turtle_2016.thumb.jpg' > Poster_The_Red_Turtle_2016.thumb.jpg
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  4110  100  4110    0     0  3463k      0 --:--:-- --:--:-- --:--:-- 4013k

$ curl -i -X DELETE 'http://localhost:8088/api/v1/photo?name=Poster_An_Interview_with_God_2018.jpg'
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 14 Mar 2024 12:21:55 GMT
Content-Length: 15

{"status":"OK"}

$ curl -i -X DELETE 'http://localhost:8088/api/v1/photo?name=Poster_The_Red_Turtle_2016.jpg'
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 14 Mar 2024 12:22:39 GMT
Content-Length: 15

{"status":"OK"}

$ curl -i -X GET 'http://localhost:8088/api/v1/listphotos'
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 14 Mar 2024 12:23:18 GMT
Content-Length: 31

{"listPhotos":[],"status":"OK"}
```
