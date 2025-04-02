## Digitalclock

Модуль создания png изображений и отправки их клиенту.

### Query paramteters
 - time - время в формате hh:mm:ss
 - k - scale factor (0 $\lt$ k $\le$ 30)

### Examples
Команда для создания изображения текущего времени:
```bash
curl --request GET --output output.png http://localhost:21000/digitalclock?k=10
```