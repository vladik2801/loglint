# loglint

 
Линтер совместим с инфраструктурой `go/analysis`, может запускаться как standalone-analyzer через `go vet -vettool=...`, а также встраивается в `golangci-lint` как module plugin.

## Что проверяет линтер

Линтер анализирует первый строковый аргумент лог-вызова и проверяет его по набору правил.

### Поддерживаемые правила

1. **`frstLower`** — сообщение должно начинаться со строчной буквы  

2. **`onlyEng`** — сообщение должно содержать только английские буквы  

3. **`noSpecial`** — сообщение не должно содержать спецсимволы и emoji  
 
4. **`noSensitive`** — сообщение не должно содержать потенциально чувствительные слова  
   
---

## Поддерживаемые логгеры

Линтер умеет распознавать вызовы для:

- `log/slog`
- `go.uber.org/zap`

### Поддерживаемые методы

Для обоих логгеров учитываются методы:

- `Debug`
- `Info`
- `Warn`
- `Error`

### Примеры поддерживаемых вызовов

```go
slog.Info("hello")
slog.Error("request failed")
zap.L().Info("hello")
zap.L().Warn("cache miss")
```

---


### Внешние зависимости

Основные зависимости проекта:

- `golang.org/x/tools/go/analysis`
- `github.com/golangci/plugin-module-register`
- `go.uber.org/zap`

---

## Сборка

### 1) Клонировать репозиторий

```bash
git clone https://github.com/vladik2801/loglint.git
cd loglint
```

### 2) Скачать зависимости

```bash
go mod download
```

### 3) Собрать standalone-анализатор

```bash
mkdir -p bin
go build -o ./bin/loglint ./cmd/loglint
```

После этого появится бинарник `./bin/loglint`, который можно использовать через `go vet -vettool`.

---

## Запуск

## Вариант 1. Запуск как standalone analyzer через `go vet`

Это основной и самый простой способ запуска анализатора локально.

### Проверка текущего модуля

```bash
go vet -vettool=$(pwd)/bin/loglint ./...
```

### Проверка конкретного примера из репозитория

```bash
go vet -vettool=$(pwd)/bin/loglint ./testdata/sample
```

Ожидаемо на примере из `testdata/sample/main.go` будет найдено нарушение, потому что строка:

```go
zap.L().Info("Hello!!! token")
```

- начинается с заглавной буквы,
- содержит спецсимволы,
- содержит чувствительное слово `token`.

### Запуск с JSON-конфигом

```bash
go vet -vettool=$(pwd)/bin/loglint -config=$(pwd)/loglint.json ./...
```

### Запуск с SuggestedFixes

```bash
go vet -vettool=$(pwd)/bin/loglint -fix ./...
```

Флаг `-fix` **не переписывает файлы автоматически сам по себе**, а прикладывает `SuggestedFixes` к диагностике анализатора.  

---

## Вариант 2. Интеграция с `golangci-lint`

Линтер реализован как **module plugin** для `golangci-lint`.

### Как это устроено

В `pluginmodule/plugin.go` регистрируется плагин `loglint`, который возвращает `logcheck.Analyzer`.  
Это позволяет подключать линтер в кастомную сборку `golangci-lint`.

### Важно

Для module plugin нужен **кастомный бинарник `golangci-lint`**, собранный с импортом вашего плагина.

---

## Сборка кастомного `golangci-lint`

Ниже способ, который повторяет логику из CI и подходит для локальной проверки.

### 1) Скачать исходники `golangci-lint`

```bash
git clone --branch v1.64.8 --single-branch --depth 1 https://github.com/golangci/golangci-lint.git /tmp/golangci-lint
```

### 2) Добавить blank import плагина

Создать файл `/tmp/golangci-lint/cmd/golangci-lint/plugins.go`:

```go
package main

import (
	_ "github.com/vladik2801/loglint/pluginmodule"
)
```

### 3) Подменить модуль на локальный путь

Находясь внутри `/tmp/golangci-lint`:

```bash
go mod edit -replace github.com/vladik2801/loglint=$(pwd)/../loglint
go mod tidy
```

> Вместо `$(pwd)/../loglint` укажи **абсолютный путь** к локальной копии репозитория `loglint`.

### 4) Собрать кастомный бинарник

```bash
go build -o ./custom-gcl ./cmd/golangci-lint
```

### 5) Запустить линтер

Из корня проекта `loglint` или любого тестируемого Go-проекта:

```bash
/tmp/golangci-lint/custom-gcl run ./...
```

---

## Тестирование

### Запуск всех тестов

```bash
go test ./...
```

---

## Конфигурация через JSON

Линтер поддерживает внешний JSON-конфиг, который передаётся через флаг:

```bash
-config=path/to/loglint.json
```

### Формат файла

```json
{
  "rules": {
    "frstLower": true,
    "onlyEng": true,
    "noSpecial": true,
    "noSensitive": true
  },
  "banned_words": [
    "password",
    "token",
    "secret"
  ],
  "banned_characters": [
    "!",
    "?",
    ":"
  ]
}
```

### Назначение полей

#### `rules`

Объект с флагами включения/выключения правил.

Поддерживаемые ключи:

- `frstLower`
- `onlyEng`
- `noSpecial`
- `noSensitive`

Пример:

```json
{
  "rules": {
    "frstLower": true,
    "onlyEng": false,
    "noSpecial": true,
    "noSensitive": true
  }
}
```

#### `banned_words`

Список запрещённых слов для правила `noSensitive`.

Пример:

```json
{
  "banned_words": ["password", "token", "api_key"]
}
```

#### `banned_characters`

Список запрещённых символов для правила `noSpecial`.

Пример:

```json
{
  "banned_characters": ["!", "?", ":", ";"]
}
```

---

## Поведение конфига по умолчанию

Если конфиг не передан, используется встроенный `DefaultConfig`.

### Встроенные правила по умолчанию

Все четыре правила включены:

```json
{
  "rules": {
    "frstLower": true,
    "onlyEng": true,
    "noSpecial": true,
    "noSensitive": true
  }
}
```

### Встроенный список чувствительных слов

По умолчанию используются такие слова:

- `password`
- `passwd`
- `token`
- `access_token`
- `refresh_token`
- `api_key`
- `apikey`
- `secret`
- `private_key`
- `ssh_key`
- `cookie`
- `session`

### Встроенный список запрещённых символов

По умолчанию запрещены:

```text
! @ # $ % ^ & * ( ) - + = { } [ ] | \ : ; " ' < > , . ? /
```

Также отфильтровываются Unicode-символы и emoji из соответствующих символьных категорий.

---

## Важные особенности конфигурации

### Если часть полей не указана

Отсутствующие поля остаются со значениями по умолчанию.

Например, такой конфиг валиден:

```json
{
  "rules": {
    "noSpecial": true
  },
  "banned_characters": [":"]
}
```

---

## Автоисправления

Сейчас `SuggestedFixes` реализованы для части сценариев.

### Что умеет исправляться автоматически

1. Приведение первой буквы к нижнему регистру  

2. Удаление спецсимволов и emoji  
   
---

## Минимальный сценарий проверки после клонирования

```bash
git clone https://github.com/vladik2801/loglint.git
cd loglint

go mod download
go build -o ./bin/loglint ./cmd/loglint

go test ./...

go vet -vettool=$(pwd)/bin/loglint ./testdata/sample
```

---

## Проверка на реальных проектах

Для демонстрации работы линтера он был протестирован на нескольких открытых Go-проектах, использующих поддерживаемые логгеры (`log/slog` и `go.uber.org/zap`).

Проекты подбирались с учётом совместимости по версии Go и наличия лог-вызовов, которые распознаются анализатором.

### 1. pingidentity/pingone-go-client

Репозиторий:  
https://github.com/pingidentity/pingone-go-client


Проверка выполнялась на пакете:

```text
examples/client_credentials
```

Команда запуска:

```bash
git clone https://github.com/pingidentity/pingone-go-client.git
cd pingone-go-client/examples/client_credentials

go vet -vettool=/path/to/loglint/bin/loglint .
```

Фактический результат: линтер корректно отработал и обнаружил нарушения правила `frstLower`.

---

### 2. uber-go/zap

Репозиторий:  
https://github.com/uber-go/zap


Команда запуска:

```bash
git clone https://github.com/uber-go/zap.git
cd zap

go vet -vettool=/path/to/loglint/bin/loglint .
```

Фактический результат: линтер корректно отработал и обнаружил нарушение правила `noSpecial`.

---