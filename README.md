# crypto-service'


План действий по доработке проекта

Вы уже реализовали всю основную логику. Теперь мы будем ее "шлифовать" и делать production-ready.
1. Конфигурация и окружение (быстрые и важные исправления)

    [ ] Задача 1.1: Унифицировать переменные портов.

        Действие: В config/config.go изменить ожидание APP_PORT на HTTP_PORT. Или наоборот, в .env.example и docker-compose.yml заменить HTTP_PORT на APP_PORT. Давайте выберем APP_PORT как единый стандарт.

    [ ] Задача 1.2: Дополнить .env.example.

        Действие: Добавить в .env.example строки COLLECTOR_INTERVAL_SECONDS=60 и COINGECKO_API_URL=....

    [ ] Задача 1.3 (Опционально, но хорошо): Сделать Makefile более надежным.

        Действие: Заменить include .env на -include .env. Знак - говорит make не падать, если файл отсутствует.

    [ ] Задача 1.4: Убрать хардкод порта из Dockerfile.

        Действие: В Dockerfile можно либо убрать строку EXPOSE 8080 совсем (так как ports в docker-compose.yml главнее), либо сделать ее динамической: ARG APP_PORT=8080, ENV APP_PORT=$APP_PORT, EXPOSE $APP_PORT.

2. Работа с БД (критически важные улучшения)

    [ ] Задача 2.1: Использовать тип decimal для денег.

        Действие: Установить библиотеку shopspring/decimal (go get github.com/shopspring/decimal). Везде, где используется price (DTO, модели DAO, таблицы БД), заменить float64 на decimal.Decimal и NUMERIC(20, 8). Это самое важное исправление.

    [ ] Задача 2.2: Нормализовать symbol в одном месте.

        Действие: В CurrencyService перед вызовом репозитория всегда делать cleanSymbol = strings.ToUpper(strings.TrimSpace(symbol)). В репозитории эту логику убрать.

    [ ] Задача 2.3: Исправить логику ConnectDB.

        Действие: В ConnectDB добавить условие: if cfg.PostgresDSN != "" { connStr = cfg.PostgresDSN } else { ...собираем по частям... }.

3. Сервис цен (Price Collector) — Повышение надежности

    [ ] Задача 3.1 (Самое важное): Добавить таймауты и контекст в HTTP-клиент.

        Действие: В PriceCollector создать кастомный http.Client с таймаутами.
        code Go

        IGNORE_WHEN_COPYING_START
        IGNORE_WHEN_COPYING_END

              
        pc.httpClient = &http.Client{
            Timeout: 30 * time.Second, // Общий таймаут
        }

            

        При запросе использовать http.NewRequestWithContext(ctx, "GET", url, nil) и pc.httpClient.Do(req). Это автоматически прервет запрос, если ctx будет отменен.

    [ ] Задача 3.2: Проверять ошибки при сохранении цены.

        Действие: В цикле сохранения цен в price_collector.go проверять appErr от pc.priceRepo.Add и логировать его: if appErr != nil { l.Error(...) }.

    [ ] Задача 3.3 (Бонусный балл): Реализовать маппинг символов.

        Действие: Создать в сервисе map[string]string ("BTC": "bitcoin", "ETH": "ethereum") и использовать его для формирования запроса к CoinGecko.

4. HTTP-ручки и DTO

    [ ] Задача 4.1: Привести /currency/price в соответствие с ТЗ.

        Действие: Изменить r.Get("/price", ...) на r.Post("/price", ...). В хендлере GetPrice изменить логику парсинга с r.URL.Query() на json.NewDecoder(r.Body).Decode(&req).

    [ ] Задача 4.2: Исправить код ответа 204.

        Действие: В CurrencyHandler.RemoveCurrency изменить w.WriteHeader(http.StatusOK) (200) на w.WriteHeader(http.StatusNoContent) (204) и полностью убрать строку json.NewEncoder(...).Encode(...).

    [ ] Задача 4.3 (Опционально): Заменить chi/middleware.Logger.

        Действие: Написать свою middleware, которая использует ваш zap-логгер. Это покажет глубокое понимание.

5. Документация и Тесты (Обязательно для Middle)

    [ ] Задача 5.1: Написать README.md.

        Действие: Пройтись по плану из предыдущего ответа и заполнить все разделы.

    [ ] Задача 5.2: Добавить Swagger.

        Действие: Пройтись по плану из предыдущего ответа: установить swag, аннотировать хендлеры, сгенерировать docs, добавить роут.

    [ ] Задача 5.3: Написать тесты.

        Действие: Восстановить юнит-тесты для service и repository, которые мы уже писали. Они были почти готовы.

Приоритеты:

    Критичные исправления: Таймауты (3.1), тип decimal для денег (2.1), соответствие ТЗ (4.1), исправление кода 204 (4.2).

    Обязательно для Middle: README.md (5.1) и Тесты (5.3).

    Сильно рекомендуется: Swagger (5.2), унификация конфигов (1.1, 1.2).

    Хорошие улучшения: Все остальное.