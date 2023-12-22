## throttler-обёртка

Реализовать throttler-обёртку для типа Transport из стандартной библиотеки
[golang.org/pkg/net/http/#Transport](https://golang.org/pkg/net/http/#Transport).
Обёртка должна реализовывать интерфейс RoundTripper
[golang.org/pkg/net/http/#RoundTripper](https://golang.org/pkg/net/http/#RoundTripper) и
инициализироваться следующими параметрами:
- RoundTripper, который будет оборачиваться
- Лимит запросов в единицу времени (целое число, если равно 0 то throttling не применяется)
- Единица времени учёта (тип time.Duration (https://golang.org/pkg/time/#Duration))
- Список HTTP методов для которых будет задействован throttling (если список пуст или nil - то
должно работать для любых методов)
- Список префиксов URL для которых будет задействован throttling (если список пуст или nil - то
должно работать для любых URL)
- Список префиксов исключений URL для которых throttling не будет задействован, даже если они
удовлетворяют условию из предыдущего фильтра (если список пуст или nil - то исключений нет)
- Флаг быстрого возврата ошибки

Если частота запросов превышает лимит, то запрос должен быть отложен до момента когда его
выполнение не вызовет превышение лимита либо завершён со специальной ошибкой (в зависимости от
флага быстрого возврата ошибки).
Для учёта частоты запросов можно считать что они выполняются мгновенно, коды возврата не имеют
значения, запросы не подпадающие под условия фильтров не учитываются.
Списки префиксов URL могут содержать * в любой части пути.
Нужно помнить что обёртка может использоваться из многих параллельных горутин, а так же может
быть использована в цепочке из нескольких обёрток.

#### Пример использования:

```go
throttled := NewThrottler(
	http.DefaultTransport,
	60,
	time.Minute, // 60 rpm
	[]string{"POST", "PUT", "DELETE"}, // limit only POST, PUT, DELETE requests
	nil, // use for all URLs
	[]string{"/servers/*/status", "/network/"}, // except servers status and network operations
	false, // wait on limit
)

client := http.Client{
	Transport: throttled,
}

// ...
resp, err:= client.GET("http://apidomain.com/network/routes") // no throttling
// ...
req := http.NewRequest("PUT", "http://apidomain.com/images/reload", nil)
resp, err:= client.Do(req) // throttling might be used
// ...
resp, err:= client.GET("http://apidomain.com/servers/1337/status?simple=true") // no throttling
// ...
```
