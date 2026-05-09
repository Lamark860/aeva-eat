# AEVA Eat — дельта к API

База — текущее API (см. корневой `README.md` репозитория). Здесь — только то, что добавляется или меняется.

## Новые сущности

### `notes` — записки от руки

```sql
CREATE TABLE notes (
  id          BIGSERIAL PRIMARY KEY,
  author_id   BIGINT NOT NULL REFERENCES users(id),
  text        TEXT NOT NULL,
  place_id    BIGINT REFERENCES places(id),  -- опционально
  city        VARCHAR(255),                  -- опционально
  paper_color VARCHAR(20),                   -- 'cream' | 'rose' | 'mint' | ...
  tape_color  VARCHAR(20),
  is_struck   BOOLEAN DEFAULT FALSE,         -- зачёркнута (wishlist → исполнено)
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

### `feed_events` — единая лента доски

Не отдельная таблица обязательно — может быть VIEW или агрегация в API.

Типы событий:
- `review_added`     — оставлен отзыв (visit)
- `note_added`       — оставлена записка
- `wishlist_added`   — добавлено в wishlist (общий)
- `wishlist_struck`  — wishlist закрыт визитом
- `gem_found`        — жемчужина найдена
- `kruzhok_added`    — видеосообщение

## Новые эндпоинты

### Лента / Доска

```
GET  /api/feed                   — общая лента всех событий круга
       ?cursor=<ts>               — пагинация по timestamp
       ?week=<YYYY-WW>            — конкретная неделя
       ?author_id=                — фильтр по автору
GET  /api/feed/weeks              — список недель с количеством событий и жемчужин
                                    (для свернутых полосок на главной)
```

### Записки

```
GET    /api/notes
POST   /api/notes                — { text, place_id?, city?, paper_color?, tape_color? }
PUT    /api/notes/:id            — author only
DELETE /api/notes/:id            — author only
PUT    /api/notes/:id/strike     — пометить зачёркнутой
```

### Города

```
GET  /api/cities                 — { city, count, gem_count, contributor_count }[]
GET  /api/cities/:name           — детали города
GET  /api/cities/:name/places    — те же фильтры что у /api/places
GET  /api/cities/:name/gems
```

### Друзья / профили

```
GET  /api/users/:id              — публичный профиль (имя, аватар, стата)
GET  /api/users/:id/places       — места, где он/она был (по визитам)
GET  /api/users/:id/gems         — его/её жемчужины
GET  /api/users/:id/cities       — его/её визиты сгруппированные по городам
```

### Жемчужины hub

```
GET  /api/gems                   — { places: [...], total, by_city, by_user }
       ?city=  &user_id=  &cuisine_id=
```

### Случайное

```
GET  /api/random
       ?city= &cuisine_id= &is_gem=true &exclude_visited_by=<user_id>
```

Возвращает один `place` или 404 если ничего не подошло.

## Изменения существующих

### `/api/places`

- Добавить параметр `is_gem=true` — уже косвенно есть
- Добавить параметр `attended_by=<user_id[]>` — фильтр по тому, кто был
- Добавить параметр `visit_from`, `visit_to` — фильтр по дате визита
- Добавить сортировки `rating_user:<user_id>` (рейтинг конкретного человека) и `rating_avg` (средний по кругу — это дефолт)

### Карточка места — расширенный ответ

Добавить в ответ `/api/places/:id` и список:
```json
{
  ...,
  "gem_status": null | { "marked_by": [user_id], "first_marked_at": "..." },
  "attendance": [{ "user_id", "user_name", "user_avatar", "visit_count" }],
  "rating_avg": { "food": 8.4, "service": 7.2, "vibe": 8.8 },
  "ratings_per_user": { "<user_id>": { "food": 9, "service": 7, "vibe": 8 } }
}
```

### Reviews — несколько фото

Уже есть в roadmap. Добавить таблицу `review_photos` (вместо одного `image_url`). MVP-критично, иначе скрапбук-галерея не работает.

### Wishlist

Текущий API оставляем, добавить:
```
GET  /api/wishlist/all           — общий wishlist круга (не только мой) — для записок на доске
```

Когда место помечено как посещённое (создан review для текущего пользователя) — соответствующая запись wishlist помечается `struck=true` (триггер или в сервисе).

## Поиск

Контекстный, не глобальный. На бэке — расширить `search` в `/api/places`, чтобы он умел искать в скоупе:
```
GET /api/places?search=<q>&scope=user:5
GET /api/places?search=<q>&scope=city:Tbilisi
GET /api/places?search=<q>&scope=gems
```

Альтернативно — каждая страница (city/user/gems) вызывает свой эндпоинт с параметром `search`.

## Авторизация

Без изменений — JWT + invite-only. RBAC `superuser/user` остаётся.

## Производительность

- `/api/feed` — кешировать на 30-60 сек (lru / redis опционально, не MVP-критично)
- `/api/cities`, `/api/users/:id` — sql VIEW с агрегатами или materialized view, обновляется по событию
- `gem_count`, `contributor_count` — считаются на лету по индексу, не денормализуем на MVP
