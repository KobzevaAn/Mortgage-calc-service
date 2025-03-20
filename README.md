# Mortgage-Calc-Service

`Mortgage-Calc-Service` — это REST API для расчета параметров ипотеки.  
Сервис позволяет рассчитывать сумму кредита, процентную ставку и ежемесячный платеж.
Хранит рассчитанные данные в кэше.

## Разработка в локальной среде
Для удобства в Makefile добавлены команды:

`make dev` запускает локальное окружения.

`make test` запускает тесты.

`make lint` запускает линтеры.

`make stop` останавливает контейнер.

`make clean` сворачивает локальное окружение.

`make all` запускает тесты линтеры и окружение.

После запуска сервис будет доступен по  http://localhost:8080/

## Настройки приложения

Настройка порта находится в`configs/config.yaml`

## API эндпоинты
Спецификация контракта доступна в `swagger.yaml` 


