# Currency Converter

Простой REST API сервис для конвертации валют.

Курсы валют получаются через CurrencyAPI и кэшируются в Redis для уменьшения количества запросов к внешнему сервису.


## Конфигурация

Создайте файл `.env` на основе `.env.example`.


### Получение API ключа

API ключ можно получить после регистрации на:

https://app.currencyapi.com/

## Запуск Redis

Если Redis не установлен локально:

```bash
docker run -d --name redis -p 6379:6379 redis
```

## Запуск приложения

```bash
go run ./cmd
```

После запуска сервис будет доступен по адресу:

```text
http://localhost:8080
```

## API

### Конвертация валют

Запрос:

```http
GET /convert?from=USD&to=RUB,CNY,EUR&amount=1
```



Пример ответа:

```json
{
  "amount": 1,
  "from": "USD", 
  "to": {
    "RUB": 78.65,
    "CNY": 7.18,
    "EUR": 0.87
  }
}
```
