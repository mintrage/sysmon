# 📊 Sysmon - Distributed System Monitoring

Sysmon — это легковесная, селф-хост система мониторинга серверов, написанная на Go. Она состоит из агентов, которые собирают системные метрики (CPU, RAM), и центрального сервера с REST API, который сохраняет данные в PostgreSQL, оповещает о падениях в Telegram и выводит графики в Grafana.

## 🚀 Возможности (Features)

* **Легковесные агенты:** Написаны на Go с использованием `gopsutil`. Потребляют минимум ресурсов, собирают реальную загрузку CPU (%) и использование RAM.
* **Централизованный сервер:** REST API бэкенд на Go с чистой архитектурой, принимающий метрики в формате JSON.
* **Telegram Watchdog:** Фоновый многопоточный надзиратель (на базе RWMutex), который следит за пульсом серверов. Автоматически пришлет алерт 🚨 при падении узла и ✅ при его восстановлении. Защищен от спама.
* **Надежное хранилище:** PostgreSQL для долговременного хранения Time-Series данных.
* **Визуализация:** Готовая интеграция с Grafana для построения красивых дашбордов в реальном времени.
* **Docker-Ready:** Вся инфраструктура (БД, Сервер, Агенты, Grafana) упакована в Docker-контейнеры для деплоя в одну команду.

## 🏗 Архитектура

1.  **sysmon-agent:** Запускается на целевых серверах (нодах). Каждые пару секунд читает `/proc` ОС и делает POST-запрос на главный сервер.
2.  **sysmon-server:** Принимает данные, валидирует их и пишет в БД. Параллельно горутина Watchdog проверяет время последней активности каждого агента.
3.  **PostgreSQL:** Хранит метрики с автоматической простановкой `created_at`.
4.  **Grafana:** Подключается к БД и рисует графики потребления ресурсов.

## 🛠 Быстрый старт (Деплой)

### 1. Запуск серверной части (Главный VDS)

Создайте файл `docker-compose.yml` и укажите ваши токены Telegram для алертов:

```yaml
services:
  db:
    image: postgres:15-alpine
    restart: always
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: 12345678     
      POSTGRES_DB: sysmon   
    volumes:
      - pgdata:/var/lib/postgresql/data
  
  server:
    image: mintrage/sysmon-server:v4
    restart: always
    ports:
      - "8080:8080"
    environment:
      - TG_TOKEN=your_telegram_bot_token
      - TG_CHAT_ID=your_telegram_chat_id
    depends_on:
      - db

  grafana:
    image: grafana/grafana-oss:latest
    restart: always
    ports:
      - "0.0.0.0:3030:3000"
    depends_on:
      - db

volumes:
  pgdata:
```

Запустите инфраструктуру:
```bash
docker compose up -d
```

### 2. Запуск агентов (на целевых серверах)

Выполните команду на любом сервере, который хотите мониторить. Укажите IP вашего главного сервера и придумайте имя для агента (оно появится в Grafana и Telegram):

```bash
docker run -d \
  --name sysmon-agent \
  --restart always \
  -e SYSMON_SERVER_URL="http://IP_ГЛАВНОГО_СЕРВЕРА:8080/api/metrics" \
  -e SYSMON_AGENT_NAME="vds-node-1" \
  mintrage/sysmon-agent:v4
```

## 📈 Настройка Grafana

1. Перейдите по адресу `http://IP_ГЛАВНОГО_СЕРВЕРА:3030` (логин/пароль по умолчанию: `admin` / `admin`).
2. Добавьте Data Source -> **PostgreSQL**.
   * Host: `db:5432`
   * Database: `sysmon`
   * User: `postgres`
   * Password: `12345678`
   * TLS/SSL Mode: `disable`
3. Создайте дашборд и используйте SQL-запросы для вывода графиков:

**Пример запроса для CPU (%):**
```sql
SELECT
  created_at AS "time",
  server_name AS metric,
  cpu_usage AS value
FROM metrics
WHERE
  $__timeFilter(created_at)
ORDER BY created_at ASC
```
*(Не забудьте выставить Unit -> Misc -> Percent (0-100) и Connect null values -> Always).*