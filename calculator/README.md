## Сервис арифметических операций

Тебе надо создать API для арифметических операций с цифрами.
API должно быть доступно только если в HTTP Header есть “User-Access”
со значением “superuser”. В случае отказа в доступе, нужно вывести
сообщение в консоли и отправить соответствующий ответ клиенту.
Должны поддерживаться только операции сложения и вычитания.
Например, с Frontend-а к тебе приходит строка “2+2-3-5+1”. В ответ ты должен
отправить JSON со статусом 200, в котором будет поле с ответом (в данном
случае это -3)