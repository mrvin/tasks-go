## Эхо-сервер на выдуманном протоколе

Реализовать tcp клиент и сервер, выполняющих сетевой обмен сообщениями на описанном ниже протоколе.

#### Сообщение протокола "запрос"
msgRequest = 1

message = msgRequest (int32)
requestString (str)

#### Сообщение протокола "ответ"
msgResponse = 2

message = msgResponse (int32)
errorNo (int32)
responseBuffer (buffer)

#### Кодировки
int32 - 8 байт, кодируется, как Big Endian

bytes - кодируется, как "сырые" данные в Little Endian

buffer - кодируется, как:
int32 - длина буфера
bytes - данные буфера
выравнивается на границу 4 байта

str - кодируется, как:
int32 - длина строки
bytes - данные строки, кодировка строки - UTF 8
выравнивается на границу 4 байта

### Сборка и запуск приложения
```bash
$ cd cmd/tcp-echo-server/
$ go build main.go
$ ./main
2025/09/14 20:17:36 Server started on port 8000
2025/09/14 20:18:15 Received request: Hello!
2025/09/14 20:18:15 Response sent successfully
2025/09/14 20:19:13 Error receive request: invalid request type

```

```bash
$ cd cmd/tcp-echo-client/
$ go build main.go
$ ./main
ErrorNo: 0
Buffer: Echo: Hello!
$ ./main -type 2
ErrorNo: 1
Buffer: invalid request type
```