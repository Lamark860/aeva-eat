# Следующая сессия — начни тут

Этот файл — пятиминутный onboarding для следующего захода в проект. Подразумевает, что предыдущая сессия закрыта, контекст сброшен.

## TL;DR

1. Дизайнер прошёлся по скринам и нарисовал v3-канвас с 5 эталонными артбордами. Лежит в [`spec/example/v3.jsx`](./example/v3.jsx) + рендер в [`spec/screenshots/v3/`](./screenshots/v3/).
2. Полный фидбэк дизайнера + приоритеты — [`R6_DESIGNER_REVIEW.md`](./R6_DESIGNER_REVIEW.md).
3. Главная боль: **80% полароидов на доске — пустые клетки**. Cover-fallback не отрабатывает + ticket-only fallback не триггерится. **Это блокер для показа кругу.**

## Где смотреть, в каком порядке

| Шаг | Документ | Что взять |
|---|---|---|
| 1 | [`R6_DESIGNER_REVIEW.md`](./R6_DESIGNER_REVIEW.md) | разделы **A** (критическое) и **B** (концептуальные провалы), список приоритетов внизу |
| 2 | [`screenshots/v3/`](./screenshots/v3/) | визуальные эталоны — какой должна выглядеть финальная вёрстка |
| 3 | [`screenshots/`](./screenshots/) (mobile-*, chrome-*, desktop-*) | как сейчас выглядит — для сравнения |
| 4 | [`OPEN-QUESTIONS.md`](./OPEN-QUESTIONS.md) | R5-Q1..Q6 закрытия (формальные ответы дизайнера) |
| 5 | [`DESIGN-DECISIONS.md`](./DESIGN-DECISIONS.md) | если нужно вспомнить старые договорённости (не переделывать) |

## Запуск окружения

```bash
cd ~/dockers/aeva-eat
docker compose up -d                   # backend, postgres, frontend, nginx
sleep 5
curl -s http://localhost:8091/         # frontend ok? 200
curl -s http://localhost:8086/api/health   # backend ok? {"status":"ok"}
```

Если нужно глянуть v3 канвас локально:

```bash
cd ~/dockers/aeva-eat/spec/example
python3 -m http.server 8788
open http://localhost:8788/index_rev3_safe.html
```

(Использовать `index_rev3_safe.html` — он ждёт пока babel-standalone доберёт все внешние JSX. Оригинальный `index_rev3.html` может не отрендерить.)

## Логины для скриншотов

| user | password | роль | данных |
|---|---|---|---|
| `lamark` | `demo12345` | superuser | 144 мест, 14 жемчужин |
| `alina` | (set if needed) | superuser | 153 мест — большой профиль для перф-тестов |
| `charlie` | (set if needed) | user | 4 мест — маленький для скринов |

Установить пароль:
```bash
HASH=$(htpasswd -bnBC 10 "" "demo12345" | tr -d ':\n')
docker exec aeva-postgres psql -U aeva -d aeva_eat -c "UPDATE users SET password_hash='$HASH' WHERE username='lamark';"
```

## Скриншот-скрипты

- `/tmp/aeva-shots-tool/shot.mjs` — основные скрины приложения (21 PNG)
- `/tmp/aeva-shots-tool/shot-v3-boards.mjs` — v3 канвас (нужен http server на 8788)
- `/tmp/aeva-shots-tool/full-canvas.mjs` — полный канвас v3 одним кадром

Команды:
```bash
cd /tmp/aeva-shots-tool
node shot.mjs /tmp/aeva-shots/final              # ~40s, 21 файла
node shot-v3-boards.mjs                          # 6 артбордов v3
```

После — копируем в `spec/screenshots/` и коммитим.

## Что сделано в R6.1 (2026-05-16)

- ✅ **A1/A2 — cover-fallback (data-issue, не SQL).** SQL `COALESCE(p.image_url, review.image_url)` отрабатывает корректно. Дебаг показал: из 185 мест **12 имеют фото** (1 в places, 11 в review/review_photos). Это не SQL-баг — это пустой dataset. **Реальный фикс:** в `ArtifactCard.vue` и `ResultCard.vue` добавлен `coverPhoto = image_url || feed_photos[0]?.url`. Места с фото только в отзывах теперь рендерят полароид. `summaryFor` в `Home.vue` приоритезирует места с фото для CollapsedStrip — архивная полоса больше не показывает 3 пустых клетки если в неделе есть хотя бы 1 фото.
- ✅ **A3 — PhotoFreeCard со Q/T/G.** Новый компонент `components/scrapbook/PhotoFreeCard.vue`. Логика: gem → G, comment≥30 → Q, иначе → T. Подключён в `ArtifactCard` через `<template v-if="isTicketOnly">`. Backend: добавлено поле `top_review_comment` в `model.Place` + SELECT в `repository/place.go` (List + GetByID).
- 📸 Контроль: `screenshots/mobile-03b-board-expanded.png` — раскрытая полоса архива с PhotoFreeCard в деле (цитаты Esmee/Пижоны, билетики на остальных).

## Что делать в первую очередь (следующий раунд)

В приоритетном порядке (полный список — `R6_DESIGNER_REVIEW.md` → Приоритеты):

1. **Штамп-помойка в шапке места (B3).** Эталон [`v3/05-place-header-ba.png`](./screenshots/v3/05-place-header-ba.png). 2 главных штампа (город + кухня) + ромб-жемчужина, категории — caveat'ом под адресом.
2. **R5-Q1: без гендер-глагола (B4).** Там же — формат `«жемчужина · Аня · 10 мая + Серёжа»`.
3. **«Любит X» logic-fix (B7).** ≥5 визитов И ≥15%.
4. **Билетик в профиле — ячейка «городов» (B7).**
5. **Маркер карты с рейтингом — откуда 3.5 (B6).**
6. **Public share (B5 + R5-Q3).** Эталон [`v3/04-public-share.png`](./screenshots/v3/04-public-share.png).
7. **Сидинг 25–30 fake мест с разнообразием** — без этого PhotoFreeCard и cover-fallback не получится показать масштабнее на скринах.

> ⚠️ Заметка о P0: оригинальный план «SQL COALESCE не работает» оказался неверным. SQL работает; проблема была в данных. Если у будущей сессии возникнет такая же гипотеза — проверять данные SQL'ом до правки кода. См. `R6_DESIGNER_REVIEW.md` секция R6.1.

## После каждой правки

- `docker compose restart frontend` или дождаться HMR
- Запустить скрин-скрипт
- Скопировать в `spec/screenshots/`
- Коммит + push
- Скриншоты — главный канал коммуникации с дизайнером, описания не работают (см. [`R6_DESIGNER_REVIEW.md`](./R6_DESIGNER_REVIEW.md) раздел «Процесс»)

## Чего НЕ делать

- Переделывать решения из `DESIGN-DECISIONS.md` без обсуждения
- Имплементировать всё подряд — пункты 0..3 обязательны, дальше согласовываем
- Реализовать «свой вариант» там, где есть v3-эталон. Эталон — закон композиции, цвета можно подобрать сами

## Известные баги, не критичные

- `/api/users/me` возвращает 404 (route не парсит `me`). Профиль использует `auth.user.id` напрямую, не блокирует
- 5 eslint-warnings (`LocationPicker.vue`, `PlaceDetail.vue`) — не блокируют сборку
- `LocationPicker.vue` — единственный неперевёрстанный Bootstrap-компонент (низкий приоритет)
