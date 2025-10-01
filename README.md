-Установка

    Для корректной работы приложения установите Go версии 1.23 или выше и добавте репозиторий данного проекта: 
        
        go get github.com/zenmaster911/L0

    Также необходимо установить сторонние пакеты вручную:

        go get -u github.com/go-chi/chi@v1.5.5 \
        go get -u github.com/go-chi/chi/v5@v5.2.2 \
        go get -u github.com/google/uuid@v1.6.0 \
        go get -u github.com/jackc/pgx/v5@v5.7.5 \
        go get -u github.com/jmoiron/sqlx@v1.4.0 \
        go get -u github.com/segmentio/kafka-go@v0.4.48 \
        go get -u github.com/spf13/viper@v1.20.1
		go get -u github.com/redis/go-redis/v9
		go get -u github.com/gojuno/minimock/v3
    
  	
	При необходимости измените значения полей в local.yaml

-Эксплуатация:

	Описание:

	 	Приложение предназначено для приёма сообщений от Kafka и записи их в базу данных PostgreSQL. Дополнительно оно поддерживает HTTP-эндпоинты для извлечения JSON-данных по уникальному идентификатору заказа (order_uid)

	Запуск приложения:

		Для запуска приложения выполните следующие шаги:
		
			-Проверьте наличие файла конфигурации local.yaml

			-Установите зависимости если ранее не сделали этого:

				go mod download
		
			-Запустите приложения с помощью команды 

				go run cmd/main.go

	Эндпоинты:

		Приложение  получает данные по заказу через GET-запрос на адрес:

			"/orders/{order_uid}"
		
		Где order_uid - уникальный идентификатор заказа.
		
		Возвращаемый ответ будет представлен в виде JSON-объекта с информацией о заказе. 
		Ниже представлен пример успешного ответа:

			{
			"order_uid": "b563feb7b2b84b6test",
			"track_number": "WBILMTESTTRACK",
			"entry": "WBIL",
			"delivery": {
				"name": "Test Testov",
				"phone": "+9720000000",
				"zip": "2639809",
				"city": "Kiryat Mozkin",
				"address": "Ploshad Mira 15",
				"region": "Kraiot",
				"email": "test@gmail.com"
			},
			"payment": {
				"transaction": "b563feb7b2b84b6test",
				"request_id": "",
				"currency": "USD",
				"provider": "wbpay",
				"amount": 1817,
				"payment_dt": 1637907727,
				"bank": "alpha",
				"delivery_cost": 1500,
				"goods_total": 317,
				"custom_fee": 0
			},
			"items": [
				{
					"chrt_id": 9934930,
					"track_number": "WBILMTESTTRACK",
					"price": 453,
					"rid": "ab4219087a764ae0btest",
					"name": "Mascaras",
					"sale": 30,
					"size": "0",
					"total_price": 317,
					"nm_id": 2389212,
					"brand": "Vivienne Sabo",
					"status": 202
				}
			],
			"locale": "en",
			"internal_signature": "",
			"customer_id": "test",
			"delivery_service": "meest",
			"shardkey": "9",
			"sm_id": 99,
			"date_created": "2021-11-26T06:22:19Z",
			"oof_shard": "1"
			}

-Структура проекта:
		
		├── cmd/                  			# Точка входа для бинарника
		│   └── main.go            			# Основной исполняемый файл
		├── frontend/              			# Фронтенд часть проекта
		│   └── index.html         			# HTML-файл для фронтенда
		├── internal/             			# Внутренний пакет с бизнес-логикой
		│   ├── config/           			# Конфигурационные файлы
		│   │   └── config.go     			# Логика чтения конфигурационных файлов
		│   ├── db/               			# Работа с базой данных
		│   │   └── db.go         			# Логика работы с БД
		│   ├── server/           			# Логика сервера
		│   │   └── server.go     			# Создание и настройка HTTP сервера
		├── pkg/                  			# Внешние публичные пакеты
		│   ├── cache/            			# Кэширование
		│   │   └── cache.go      			# Логика кэширования
		│   ├── handler/          			# Обработка запросов HTTP
		│   │   ├── handler.go     			# Основной обработчик
		│   │   ├── order_handler.go 		# Обработчик заказов
		│   │   └── middleware.go  			# Логика middleware
		│   ├── kafka_consumer/    			# Потребитель Kafka
		│   │   └── kafka_consumer.go 		# Логика потребителя Kafka
		│   ├── model/            			# Модели данных
		│   │   └── model.go      			# Определение моделей
		│   ├── repository/       			# Репозитории для работы с данными
		│   │   ├── cache_repo.go  			# Репозиторий для кэша
		│   │   ├── customer_repo.go 		# Репозиторий для клиентов
		│   │   ├── deliveries_repo.go 		# Репозиторий для доставок
		│   │   ├── items_repo.go  			# Репозиторий для товаров
		│   │   ├── order_repo.go  			# Репозиторий для заказов
		│   │   ├── repository.go  			# Основной репозиторий
		│   │   └── Statuscheck.go 			# Логика проверки статуса
		│   ├── service/          			# Бизнес-сервисы приложения
		│   │   ├── cache_service.go 		# Сервис для кэша
		│   │   ├── customer_service.go 	# Сервис для клиентов
		│   │   ├── deliveries_service.go 	# Сервис для доставок
		│   │   ├── items_service.go 		# Сервис для товаров
		│   │   ├── order_service.go 		# Сервис для заказов
		│   │   └── service.go    			# Основной сервис
		│   └── worker/           			# Логика воркера
		│       └── worker.go      			# Основной воркер
		├── schema/               			# Схемы базы данных
		│   ├── 20250809080419_init.down.sql # Миграция базы данных вниз
		│   └── 20250809080419_init.up.sql 	# Миграция базы данных вверх
		├── .gitignore            			# Игнорируемые файлы для Git
		├── docker-compose.yaml    			# Конфигурация Docker Compose
		├── docker-compose.yamlZoneId... 	# Дополнительная конфигурация Docker Compose
		├── go.mod                			# Менеджмент зависимостей Go
		├── go.sum                			# Хеш-зависимости Go
		├── local.yaml            			# Локальная конфигурация
		└── README.md             			# Текущий документ



