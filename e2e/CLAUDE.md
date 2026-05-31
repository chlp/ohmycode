# CLAUDE.md — E2E тесты (Playwright)

Этот файл описывает E2E-тесты для OhMyCode и даёт всё необходимое для их запуска и расширения.

## Структура тестов

```
e2e/
├── editor.spec.js        — базовое поведение редактора (загрузка, WS, CodeMirror)
├── static-files.spec.js  — cache busting: hash в URL, Cache-Control, immutable-заголовки
├── playwright.config.js  — конфигурация (chromium, baseURL, timeout, retries)
└── package.json          — зависимости (@playwright/test)
```

Смежные слои тестирования (не E2E):

| Слой | Расположение | Команда |
|------|-------------|---------|
| Go unit + WS integration | `api/internal/...` | `cd api && go test -race ./...` |
| JS unit (Vitest) | `api/internal/api/client/js/test/` | `cd api/internal/api/client && npm test` |
| **E2E (Playwright)** | `e2e/` | см. ниже |

## Требования перед запуском

1. **Приложение должно быть запущено** на `http://localhost:52674/`

```bash
# Запустить API + MongoDB
docker compose -f api/docker/docker-compose.yml up -d --build --remove-orphans --force-recreate

# Запустить runner'ы (Python, Go, Node.js, Java, PHP, PostgreSQL, MySQL)
docker compose -f runner/docker/docker-compose.yml up -d --build --remove-orphans --force-recreate
```

Docker должен быть запущен. Если Docker Desktop не работает: `open -a Docker && sleep 10`.

2. **Зависимости E2E** (один раз):

```bash
cd e2e
npm install
npx playwright install chromium
```

## Запуск тестов

```bash
cd e2e

# Все тесты
npm test

# С подробным выводом
npx playwright test --reporter=list

# Один файл
npx playwright test editor.spec.js

# Интерактивный UI-режим (полезно при отладке)
npm run test:ui

# Против другого окружения
APP_URL=https://ohmycode.work npx playwright test
```

## Что проверяют тесты

### `editor.spec.js`

- Приложение отдаёт 200, `#content` textarea присутствует в DOM
- При загрузке URL меняется на `/[22-символьный base62 ID]`
- `body` переходит из `opacity:0` в `opacity:1` после DOMContentLoaded
- WebSocket подключается к `/file`
- CodeMirror рендерится и виден
- `#file-header` виден
- После WS-рукопожатия редактор получает снапшот файла

### `static-files.spec.js`

- `index.html` отдаётся с `Cache-Control: no-cache`
- Все `?v=` ссылки в `index.html` содержат hex-хэш (не числа вида `?v=18`)
- JS-модули содержат версионированные relative imports (`./foo.js?v=HASH`)
- Версионированные ассеты получают `Cache-Control: immutable, max-age=31536000`
- Хэш стабилен между запросами

## Добавление новых тестов

Новые `.spec.js` файлы в `e2e/` подхватываются автоматически (см. `testDir: '.'` в конфиге).

Типовой шаблон:

```js
import { test, expect } from '@playwright/test';

test('описание', async ({ page }) => {
  await page.goto('/');
  await page.waitForLoadState('networkidle');
  // ...
});
```

**Важные детали:**

- `baseURL` — `http://localhost:52674` по умолчанию, переопределяется через `APP_URL`
- Timeout: 30 секунд на тест, 1 retry при падении
- Только Chromium (намеренно: приложение ориентировано на desktop Chrome)
- `retries: 1` — не считай тест сломанным с первой попытки, особенно WS-зависимые

## Типичные ошибки и диагностика

| Симптом | Причина | Решение |
|---------|---------|---------|
| `Connection refused` при запуске | Приложение не запущено | Запустить docker compose |
| `ERR_MODULE_NOT_FOUND: @playwright/test` | Не установлены зависимости | `npm install` в `e2e/` |
| `Cannot find package 'playwright'` | Запуск из неправильной директории | `cd e2e` перед командами |
| Тест по WS падает на первой попытке | Race condition при подключении | Нормально — retry сработает |
| `test-results/` в git status | Не добавлен `.gitignore` | Файл уже есть: `e2e/.gitignore` |

## Для Claude

**Перед запуском тестов всегда проверяй:**

```bash
curl -s -o /dev/null -w "%{http_code}" http://localhost:52674/
# Должно вернуть 200. Если нет — сначала запусти docker compose (см. выше).
```

**Рабочая директория для E2E:** всегда `cd e2e` перед `npx playwright`.  
Запуск из корня `/Users/amini/Documents/github/ohmycode` или из `api/` даст `ERR_MODULE_NOT_FOUND`.

**Docker daemon:** если получаешь `Cannot connect to the Docker daemon` — выполни `open -a Docker` и подожди ~10 секунд до повторной попытки.

**Скриншоты при отладке:** Playwright сохраняет трейсы в `e2e/test-results/` при падении с retry. Файл `trace.zip` можно просмотреть командой:

```bash
npx playwright show-trace e2e/test-results/<test-name>/trace.zip
```

**Написание исследовательских тестов:** для UX-проверок используй `page.keyboard.insertText()` вместо `page.keyboard.type()` — `type()` эмулирует посимвольный ввод и может потерять фокус в CodeMirror до подключения WS (редактор `readOnly: true` до получения первого снапшота).
