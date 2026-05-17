# Следующая сессия — начни тут

Onboarding для следующего захода. Контекст сброшен.

## TL;DR (статус на 2026-05-17 — R7)

Сессия R7 закрыла **большой кусок долгов R6** + структурные D1-D3 + 8 пользовательских багов.

**Закрыто в R7:**
- Все 3 структурных долга D1-D3 (D4 — kruzhok-miss телеметрия, не делано)
- Архитектурный R6.11 регресс на /places и /map — drawer теперь с локальным draft
- 8 пользовательских багов (видео-кружочки, дубли визитов, hit-target маркера, кластер-стиль карты, скотч съехал, фильтры на /map, скролл-лок, comment на фото-карточках)

**Долги:**
- D4 — `kruzhok-miss` телеметрия. Бэк нужен. Низкий приоритет до прод-MVP.
- v4-канвас от дизайнера ещё не пришёл (записка без места, wishlist активный/зачёркнутый, AddArtifactSheet, мини-карта города)
- `LocationPicker.vue` единственный Bootstrap-компонент, 3 pre-existing eslint warnings
- Yandex hit-detect баг после open/close balloon — Yandex-specific, юзер согласился оставить

## R7 — что сделано (2026-05-17)

### Архитектурные (D1-D3)

- **D1 частично** — `KruzhokStack.vue` вынесен из `ArtifactCard.vue` (видео-кружочки с inline-play, ~185 строк, prop `size: normal|large`). `ArtifactCard.vue` сократился 584 → 365 строк. Полный extract `ArtifactPolaroid` отложен — специфичная вёрстка stack-caption/featured-overrides, высокий риск визуальных регрессий без тестов.
- **D2** — pickFeaturedId иерархия `gem > video > первый`. Max 1 full-width на неделю. Убрала автоматическое `.has-video → grid-column: 1/-1` — full-width управляется только через `.featured`.
- **D3** — `grid-auto-flow: dense` удалён из `.doska-week`. Пустые ячейки после full-width допустимы — иначе записка вторника визуально оказывается под визитом среды.

### R6.11 архитектурный регресс на /places и /map

Drawer фильтров писал **прямо в `placesStore.filters`** → каждый тап чипа изменял chip-row снаружи и activeFilterCount, прятал полки, показывал «применённые» результаты (старый список) даже без `fetchPlaces`. После close ✕ без apply — visual «применено», бэк не вызван, тап на `×` чипа триггерил **другую** комбинацию.

Фикс: локальный `draft` ref в обоих view (Places.vue + MapPage.vue). `show.bs.offcanvas` → `seedDraft(store → draft)`. Все `v-model` в drawer → `draft.*`. «Применить» → копирует `draft → store + fetchPlaces`. Cancel — draft выкидывается, store не меняется. Chip-row и find-results снаружи продолжают читать `placesStore.filters` (только applied).

### Карта (множество фиксов)

- **`.sb-paper > *` override** ставил всем потомкам `position:relative; z-index:1` (та же specificity что `.offcanvas`, но позже в каскаде). Drawer был встроен в layout, не плавал поверх → на /map уезжал под Yandex-карту, кнопки внизу недоступны. Фикс: `.sb-paper > .offcanvas { position: fixed; z-index: 1045 }` в `scrapbook.scss`.
- **MapPage.vue watch deep на filters** триггерил `fetchAllPlaces` при каждом тапе → reload карты. Заменён на same draft-pattern.
- **Маркер карты — тап в стороне.** iconLayout упрощён до прямого SVG (убрана обёртка `<div position:relative>`). `iconShape` Rectangle с padding 4px покрывает весь визуальный pin.
- **Клик на маркер сразу открывал /places/:id** вместо balloon. Убрала `pm.events.add('click', emit)` — Yandex по дефолту открывает balloon с preview (бумажная плашка, обложка, рейтинг, ссылка «подробнее →»).
- **Yandex Clusterer добавлен** — 215 маркеров наложены на низком зуме блокировали клики друг другу. Кастомный `ClusterLayout`: бумажный кружок paper-card с terra-обводкой + пунктирный inner-ring + число рукописным Caveat (matches штампы).
- **После open/close balloon — кривые клики** (Yandex-specific): убрала `hideIconOnBalloonOpen: false` — Yandex defaults корректно сбрасывают event-listeners.

### Доска: визиты + комменты + видео

- **3 одинаковых полароида (Винни, Mansarda)** — group key был `(place_id, day)`, разные дни одного места = разные карточки. Заменено на `(place_id, week-bucket)` через `bucketKey(weekStart(evDate), now).key`. Один полароид на место в неделю. Если visits > 1 → caption «×N за неделю», ratings/комменты/видео усреднены/объединены.
- **Скотч «съехал»** на карточках Винни. После group-by-week `place.id` стал composite строкой `"242-2026-04-29"`. Hash-функции `tapeStyle/tapeVariant` делали `(place.id) % 4` → для строки = `NaN` → tape без позиции. Фикс: hash на `_placeId ?? id` (число → корректный остаток).
- **Router-link тоже сломался** — `/places/242-2026-04-29`. Фикс: `:to="/places/${place._placeId ?? place.id}"`.
- **Comment не показывался на фото-карточках** — раньше только в PFC (ticket-only). Добавлен `<p class="sb-comment">«…»</p>` под meta, обрезка 120 chars с многоточием. На featured крупнее (17px), обычные 15px caveat.
- **Короткий comment (<30 chars) у PFC T-layout не виден** (Q-layout требует ≥30). Добавлен `.pfc-short-comment` под ticket.

### Найти: жемчужины overflow

- `.shelf-gem` width 130 → 156. Полароид 118+padding=138, с tilt ±3° визуально 148-152 — раньше выезжал за края контейнера.

### Прочее

- **Время HH:MM на caption + notes** (R6.10/R6.11 уже было)
- **Drawer auto-close + body.overflow cleanup** работают
- **Scroll-to-top через bottom-tab** работает

## Где смотреть, в каком порядке

| Шаг | Документ | Что взять |
|---|---|---|
| 1 | этот `NEXT_SESSION.md` | TL;DR + долги |
| 2 | [`R6_DESIGNER_REVIEW.md`](./R6_DESIGNER_REVIEW.md) | разделы **A** (критическое) и **B** (концептуальные), список приоритетов внизу |
| 3 | [`screenshots/`](./screenshots/) | 12 актуальных скринов всех уникальных страниц (после R7) |
| 4 | [`screenshots/v3/`](./screenshots/v3/) | визуальные эталоны дизайнера |
| 5 | [`OPEN-QUESTIONS.md`](./OPEN-QUESTIONS.md) | R5-Q1..Q6 закрыты |
| 6 | [`DESIGN-DECISIONS.md`](./DESIGN-DECISIONS.md) | старые договорённости |

## Запуск окружения

```bash
cd ~/dockers/aeva-eat
docker compose up -d
sleep 5
curl -s http://localhost:8091/
curl -s http://localhost:8086/api/health
```

## Логины

| user | password | роль | данных |
|---|---|---|---|
| `lamark` | `demo12345` | superuser | 145 мест, 14 жемчужин |
| `alina` | (set if needed) | superuser | 153 мест |
| `charlie` | (set if needed) | user | 4 мест |

Установить пароль:
```bash
HASH=$(htpasswd -bnBC 10 "" "demo12345" | tr -d ':\n')
docker exec aeva-postgres psql -U aeva -d aeva_eat -c "UPDATE users SET password_hash='$HASH' WHERE username='lamark';"
```

## Скриншоты (12 шт, актуальные)

```
spec/screenshots/
├── 01-doska.png                # Главная Доска с week-bucket grouping
├── 02-doska-expanded.png       # +5 артефактов раскрыто
├── 03-naiti-shelves.png        # Найти с полками (Жемчужины, города, кухни, друзья)
├── 04-naiti-filter-drawer.png  # Открытый фильтр-drawer с draft state
├── 05-map.png                  # Карта со скрапбук-кластерами
├── 06-profile.png              # Профиль lamark (Визиты)
├── 07-place-detail.png         # Северяне (3 штампа, жемчужина, отзывы)
├── 08-city.png                 # Город Москва (B1)
├── 09-person.png               # Person lamark (B8 — 3457px)
├── 10-gems.png                 # Жемчужины-хаб (B2)
├── 11-public-share.png         # /p/220 публичная превью
└── 12-login.png                # Login wordmark + paper-card
```

## Что делать в первую очередь

### Если v4-канвас от дизайнера пришёл
1. Записка на Доске без места — рендер `NoteArtifact` без `place_id` рядом с полароидами
2. Wishlist активный/зачёркнутый — отдельный визуал
3. AddArtifactSheet — bottom-sheet выбора типа артефакта
4. Мини-карта города (B1.2) — bbox по местам города на `/cities/:name`

### Если канваса нет — низкоприоритетные хвосты
1. **D4** — `kruzhok-miss` телеметрия. Логировать на бэк когда юзер тапает на /places/:id из карточки с видео без play. Если >10% → расширить hit-target кружочка.
2. **LocationPicker.vue** — перевёрстка под скрапбук, 3 eslint warnings.
3. **Yandex hit-detect** — глубже разобраться, почему после open/close balloon Yandex плохо ресетит handler.
4. **Перевод seed/demo на CI** — `seed_demo.sh` сейчас вручную.

### Известные мелочи косметики
- «Жемчужины Москва» — без родительного падежа («Москвы»). Низкий приоритет, не влияет на UX.
- Карэм с «Не ходил вот и не нравится» — тестовая запись юзера в верху Доски, не баг. Удалить вручную если хочется чистый скрин.

## После каждой правки

```bash
docker compose restart frontend  # или дождаться HMR
# скрины в spec/screenshots/
git add . && git commit
```

## Чего НЕ делать

- Переделывать решения из `DESIGN-DECISIONS.md` без обсуждения
- Имплементировать всё подряд — спрашивать про приоритет
- Реализовать «свой вариант» там, где есть v3-эталон

## Известные баги, не критичные

- `/api/users/me` возвращает 404 (route не парсит `me`). Профиль использует `auth.user.id` напрямую
- 3 eslint-warnings (`LocationPicker.vue`) — не блокируют сборку
- Yandex hit-detect после balloon open/close — Yandex-specific, оставлено
