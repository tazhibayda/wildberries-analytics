# Что сейчас есть
Пока что поднимается kafka и zookeeper, подключается к wbapi и получает данные. Отправляет их в топик kafka.

# Как запустить 
example.env - создать по типу
```.env
WB_TOKEN="Bearer your_token_here"
KAFKA_BROKERS="kafka:9092"
KAFKA_TOPIC="wb.raw"
```
дальше:
```terminal
docker compose -f docker-compose.yml up -d zookeeper kafka
docker compose -f docker-compose.yml up wildberriesapi
```
После успешного запуска можно проверить, что данные приходят в топик:
```
docker exec -it kafka /usr/bin/kafka-console-consumer   --bootstrap-server kafka:9092   --topic wb.raw   --from-beginning
```


# Что дальше
Go:
- Добавить обработку ошибок
- Добавить логирование
- Добавить тесты 

Архитектура:
- Добавить ClickHouse
- Принять данные из Kafka и положить в ClickHouse
- Соединить с главным docker-compose 

Функционал:
- Добавить получение данных по фильрам
- 


Это только пример использование wb api с kafka, ткпкрь нужно правильно обработать данных через ETL-service



internal/api/
├── client.go               ← WBClient struct + DoRequest() + retry logic
├── endpoints.go            ← все URL WB API
├── paid_storage.go         ← методы start/status/download
├── prices.go               ← методы get_prices
├── tariffs.go              ← методы get_tariffs
├── adverts.go              ← методы для рекламных API
├── reports.go              ← nm_report_history/detail и т.д.
├── finance.go              ← финансы, возвраты, supplies
└── search_texts.go         ← POST отчёт по запросам
