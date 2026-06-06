# Тех-долг и продуктовые пробелы — AEVA Eat

> Реестр заведён 2026-06-05 по итогам сквозного аудита «MVP → продукт»,
> актуализирован 2026-06-06 после цикла из 15 батчей (+ хотфикс).
> Severity была откалибрована состязательной перепроверкой по коду.
> Статусы: ✅ сделано · 🔸 открыто · 💤 отложено осознанно.
>
> Связано: [`product.md`](./product.md) · [`backend.md`](./backend.md) ·
> [`../ROADMAP.md`](../ROADMAP.md) · [`../CHANGELOG.md`](../CHANGELOG.md)

**Итог цикла:** все critical/high/medium из аудита закрыты и на проде. Открытыми
остались только (а) действия в инфраструктуре пользователя (секреты/cron/настройки
репо), (б) низкоприоритетные мелочи, (в) осознанно отложенное. Подробная история
изменений — в `CHANGELOG.md`.

---

## ✅ Сделано

### Батч 1 — идентичность места + базовая безопасность
| Что | Где |
|---|---|
| Идентичность места `name+address+city` (фикс «нельзя добавить филиал сети») | migration 016, `repository/place.go`, `handler/place.go` |
| 409 не тупик: тело существующего места + развилка «перейти/уточнить адрес» | `handler/place.go`, `stores/places.js`, `views/PlaceForm.vue` |
| Типизированный детект конфликта (`pq.Error` 23505, не `strings.Contains`) | `repository/place.go` |
| Правка отзыва не стирает `image_url`/`video_url` | `repository/review.go` Update |
| IDOR на удалении инвайта закрыт (создатель/superuser, проверка в SQL) | `handler/invite.go`, `repository/invite.go` |
| Fail-fast по `JWT_SECRET` в `APP_ENV=production` | `config.go`, `main.go`, оба compose |
| Скрипт бэкапа БД + uploads с ротацией | `backend/scripts/backup.sh`, `DEPLOY.md` |

### Батч 2 — надёжность
| Что | Где |
|---|---|
| Superuser-bypass на места (`canMutatePlace`) | `handler/place.go` |
| Чистка файлов `/uploads` при удалении места + confirm с числом отзывов | `handler/place.go`, `repository/place.go`, `views/PlaceDetail.vue` |
| Graceful degradation карт (ymaps не загрузился → ручной ввод) | `index.html`, `LocationPicker.vue`, `MapView.vue`, `PlaceForm.vue` |
| Таймауты HTTP-сервера; suggest-прокси с таймаутом+контекстом | `main.go`, `handler/suggest.go` |
| Видео по сигнатуре (magic bytes); лимит длины правки заметки; toast-ошибки delete/wishlist | `handler/review.go`, `handler/note.go`, `views/PlaceDetail.vue` |

### Батч 3 — защита и поиск
| Что | Где |
|---|---|
| Rate-limit (login/register 20/min по IP, suggest 60/min по userID) | `middleware/ratelimit.go`, `main.go` |
| Лимит JSON-тела 1 MB (`MaxBytesReader`) | `middleware.BodyLimit`, `main.go` |
| Поиск по name + city + cuisine | `repository/place.go` List |
| Сохранение `website` (из Яндекс-подсказки и ручного поля) | `views/PlaceForm.vue` |

### Батч 4 — данные и оркестрация
| Что | Где |
|---|---|
| Канонизация города (trim + приведение к существующему написанию) | `handler/place.go` `normalizePlaceReq`, `repository/place.go` `CanonicalCity` |
| Backend healthcheck + nginx ждёт `service_healthy` | `docker-compose.prod.yml` |
| Удалён мёртвый `VideoKruzhok.vue` | `components/scrapbook/` |

### Батч 5 — кэш, корректность, логи
| Что | Где |
|---|---|
| Кэш suggest (TTL 60s) — экономия квоты Яндекса | `handler/suggest.go` |
| review принадлежит place из URL (Update/Delete, 404 иначе) | `handler/review.go`, `repository/review.go` |
| Ротация docker-логов (10m × 3) всем сервисам | `docker-compose.prod.yml` |

### Батч 6 — продуктовые развилки (по решению пользователя)
| Что | Где |
|---|---|
| Мягкое удаление мест (архив, отзывы/фото круга сохраняются) + restore | migration 017, `repository/place.go` `SoftDelete`/`Restore`, `handler/place.go` |
| Редактирование общего места любым участником круга | `handler/place.go` Update/UploadImage |
| Wishlist «исполнение желания» круговое (визит кем угодно гасит у всех) | `handler/review.go`, `repository/wishlist.go` `MarkStruckByPlace` |
| Auto-merge `--squash` → `--rebase` (dev/master не расходятся) | `.github/workflows/automerge.yml` |

### Батч 7 — доступность и корректность
| Что | Где |
|---|---|
| a11y подсказок + валидный HTML (`role`/`aria-activedescendant`) | `LocationPicker.vue` |
| TOCTOU лимита фото устранён (`SELECT … FOR UPDATE`, `ErrPhotoLimit`) | `repository/review.go` `AddPhoto` |

### Батч 8–15 — производительность, инфра, продуктовые/крупные
| Что | Где | Батч |
|---|---|---|
| N+1 в `List` («Найти»): reviewers+feedPhotos батч-запросами | `repository/place.go` | 8 |
| Миграции через леджер `schema_migrations` (по разу, в транзакции) | `cmd/api/main.go` | 9 |
| CI прогоняет миграции на реальном Postgres | `.github/workflows/ci.yml` job `migrations` | 10 |
| `environment: production` (точка для апрува прода) | `.github/workflows/deploy.yml` | 10 |
| Лимиты mem/cpu всем сервисам (потолки: pg 768m, backend 512m, front/nginx 128m) | `docker-compose.prod.yml` | 10 |
| N+1 профиль/жемчужины: `loadPlaces` → `GetManyByIDs` (base+reviewers+photos+gem_status батчем) | `repository/place.go`, `handler/aggregate.go` | 11 |
| Рефактор: общий `placeSelectCols`/`placeBaseFrom`/`scanPlace` (устранён дрейф List↔GetByID) | `repository/place.go` | 11 |
| Интеграционный тест против реального Postgres (gated `AEVA_TEST_DSN`) | `repository/place_db_test.go` | 11 |
| Шаринг `/p/<uuid-token>` вместо `/p/<id>` (нет энумерации каталога) + кнопка «поделиться» | migration 018, `handler/share.go`, `repository/place.go` `GetByShareToken`, nginx, `views/PlaceDetail.vue` | 12 |
| Унификация копирайта: user-facing «заведение» → «место» | `LocationPicker/PlaceDetail/Profile/PlaceForm` | 13 |
| Деплой пересоздаёт nginx (`up -d --force-recreate nginx`) — bind-mount inode-gotcha, конфиг не применялся | `.github/workflows/deploy.yml` | 13–14 |
| Хотфикс: закавычен nginx-регэксп `/p/` (`{36}` без кавычек ронял nginx → прод down) | `nginx/nginx.prod.conf`, `nginx/nginx.conf` | хотфикс |
| CI валидирует nginx-конфиги (`nginx -t` с резолвом upstream) | `.github/workflows/ci.yml` job `nginx` | 15 |

---

## 🔸 Открыто

### Действия в инфраструктуре пользователя (не код)
- **[ops] Бэкап на проде не на расписании.** Скрипт `backend/scripts/backup.sh`
  готов — осталось завести cron на VPS + off-site копию (инструкция в `DEPLOY.md`).
  До этого риск потери данных при отказе диска сохраняется.
- **[ops] `GH_PAT` не задан** → merge от `github-actions[bot]` не триггерит
  `deploy.yml`, деплой дотягивается вручную `gh workflow run deploy.yml`. Фикс —
  завести `GH_PAT` в repo→Settings→Secrets.
- **[ops] Апрув прода не включён.** `environment: production` подключён в
  `deploy.yml`; чтобы заработал апрув — задать required reviewers в
  repo→Settings→Environments→production.

### Низкий приоритет
- **[ops] Откат forward-only.** Деплой = `git reset --hard origin/master`;
  rollback на предыдущий SHA/тег не автоматизирован. `*.down.sql` существуют, но
  применяются вручную.
- **[data-model] Агрегатные счётчики включают архивированные места.** Полки
  «Москва · N» считают и soft-deleted (листинги/карта/рандом их уже не
  показывают). Редкий кейс, косметика.
- **[data-model] `created_by IS NULL` не переназначается.** Латентно: в проде
  мест без владельца нет (создаются с `created_by`, пользователи не удаляются).
- **[data-model] 2 лёгкие копии place-SELECT в `wishlist.go`** (без
  videos/top_comment) — консолидация добавила бы им ненужные колонки, оставлены
  намеренно. Главный дрейф (List↔GetByID) уже устранён.
- **[product] Город: справочник/автокомплит + rename.** Канонизация по
  существующему написанию есть (батч 4); полноценный справочник городов и
  массовый rename — опционально.
- **[frontend] Маппинг Яндекс-`categories` в кухню/категории.** `website` уже
  сохраняется (батч 3); авто-предзаполнение кухни/категорий из подсказки и
  `uri`/org-id (задел под фазу 3 идентичности места) — не сделано.
- **[frontend] Единый враппер мутаций.** `PlaceDetail` обёрнут в try/catch+toast
  (батч 2), но паттерн повторяется по другим views — стоит вынести в общий хелпер.
- **[frontend] `LocationPicker`→`PlacePicker` рейн** (внутреннее имя, не
  пользовательское) + сосуществование Bootstrap и scrapbook-CSS (полное
  вытеснение Bootstrap — крупный риск, отложено).
- **[security] CORS `*` + `AllowCredentials:true`.** Инертно (токен в заголовке
  Authorization, не в куках), но «на вырост» стоит сузить origin.

---

## 💤 Отложено осознанно

- **Единый серверный `/feed` вместо клиентской сборки.** `useFeed.js` держит
  плотную ПРЕЗЕНТАЦИОННУЮ логику доски — недельные/месячные бакеты, группировку
  по (place_id, неделя), усреднение оценок, сбор видео, резолв attendees, выбор
  Q-layout-цитаты, мердж note+wishlist. Перенос в API запёк бы вёрстку доски в
  бэкенд и означал высокорисковый рефактор сердца приложения при сомнительной
  выгоде на камерном масштабе. Текущий дизайн (хронология из `feed_events` +
  клиентское обогащение) — оправданный BFF-паттерн. Выносить бакетинг в API
  отдельной взвешенной работой, если доска вырастет.
- **Модель «бренд → филиал» (brands над places).** Правильная доменная модель,
  но для дружеского круга преждевременна (двухшаговая форма, авто-склейка
  одноимённых мест в разных городах). Идентичность места (`name+address+city`)
  уже чинит исходный баг без этой сущности. Идеальный end-state — единая таблица
  `places`=филиал с приоритетом ключей: Яндекс `org_id` → координатный guard →
  текстовый `name+address+city`. Вводить, когда филиалов одной сети станет
  столько, что доска без группировки шумит.
- **Тёмная тема, page-transitions, «Ваш год»** — см. `ROADMAP.md` (next-next).
