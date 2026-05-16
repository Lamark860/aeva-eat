# AEVA Eat — для дизайнера

Точка входа для синка по дизайну. Здесь — что на каких экранах сейчас, ссылки на детальные доки, скриншоты последней сборки и список открытых вопросов.

Дата последнего обновления: **2026-05-16** — добавлены v3 концепты от дизайнера.

---

## 🎨 v3 — текущие эталоны (приоритет)

Дизайнер отрисовал 5 болевых точек на канвасе `spec/example/index_rev3.html`. К каждому экрану из задач R6 теперь есть конкретный референс.

Рендер сохранён в [`screenshots/v3/`](./screenshots/v3/):

- [`01-photofree-card.png`](./screenshots/v3/01-photofree-card.png) — **безфотный артефакт** (3 раскладки: цитата / билетик / штамп)
- [`02-city-guide.png`](./screenshots/v3/02-city-guide.png) — **город как путеводитель** (Казань)
- [`03-gems-hub.png`](./screenshots/v3/03-gems-hub.png) — **жемчужины-сокровищница**
- [`04-public-share.png`](./screenshots/v3/04-public-share.png) — **/p/:id** с правильными пропорциями и бумажной CTA
- [`05-place-header-ba.png`](./screenshots/v3/05-place-header-ba.png) — **шапка места до/после** (2 штампа вместо 7, без гендер-глагола)
- [`00-canvas-full.png`](./screenshots/v3/00-canvas-full.png) — обзор всего канваса

**v3 НЕ покрывает:**
- Главный экран (Доска) — оставлен как v2, точки боли v3 в нём нет
- Записку на Доске без места — следующий раунд
- Wishlist-артефакт активный/зачёркнутый — следующий раунд
- AddArtifactSheet — следующий раунд

Чтобы запустить канвас локально:

```bash
cd spec/example
python3 -m http.server 8788
# открыть http://localhost:8788/index_rev3_safe.html
```

---

## Что читать в первую очередь

| Документ | О чём | Когда смотреть |
|---|---|---|
| [`R6_DESIGNER_REVIEW.md`](./R6_DESIGNER_REVIEW.md) | **Активный раунд.** Фидбэк дизайнера + v3-эталоны + конкретный список задач с приоритетами 0..18 | стартовая точка для следующей сессии разработки |
| [`STATUS.md`](./STATUS.md) | Хронологический живой статус всех раундов (R1 → R5.1) | если хочется понять «что и когда сделали» |
| [`DESIGN-DECISIONS.md`](./DESIGN-DECISIONS.md) | Verbatim-решения дизайнера (19 вопросов первой волны + 10 второй) | если нужна точная формулировка договорённости |
| [`OPEN-QUESTIONS.md`](./OPEN-QUESTIONS.md) | Закрытие R5-Q1..Q6 + ссылка на R6 | формальный сводник статусов |
| [`NEXT.md`](./NEXT.md) | Продуктовый горизонт после хендоффа | стратегические штрихи, не текущий цикл |
| [`mvp-scope.md`](./mvp-scope.md), [`product.md`](./product.md), [`design.md`](./design.md), [`backend.md`](./backend.md) | Изначальные спеки от продукта/дизайна | если нужен исходный контекст |
| [`example/`](./example/) | Дизайнерские прототипы (v2.jsx, v3.jsx, scrapbook.css) | референс ритма и палитры; рендерится через `index_rev3_safe.html` |

---

## Открытые вопросы (R5 — закрыты, R6 — в работе)

R5-Q1..Q6 закрыты дизайнером 2026-05-15. Реализация ещё не сделана — см. [`R6_DESIGNER_REVIEW.md`](./R6_DESIGNER_REVIEW.md) приоритеты `#2`, `#6`, `#11`, `#12`, `#13`.

| # | Закрыто как | Где смотреть в коде/скрине |
|---|---|---|
| **R5-Q1** | (b) Убрать глагол. Формат `«жемчужина · Аня · 12 марта»` | эталон [`v3/05-place-header-ba.png`](./screenshots/v3/05-place-header-ba.png) низ |
| **R5-Q2** | (a) Оставить как есть (event-driven) | `mobile-03-board.png` (текущая неделя пустая когда нет визитов — это правильно) |
| **R5-Q3** | (a) `?next=/places/:id` + redirect after login | `mobile-14-public-share.png` (CTA сейчас уходит на `/login` без `?next`) |
| **R5-Q4** | (a) Лёгкий 1–2° tilt на wishlist-карточках | (нет в seed-данных, появится с реальным использованием) |
| **R5-Q5** | (b+c) Оставить `forcePoster`, добавить «тапни, чтобы посмотреть» после 500мс | `mobile-11-place-video.png` |
| **R5-Q6** | (a) Разделитель-плашка «… а вот эти Аня ещё не пробовала» | `chrome-find.png` drawer фильтров (опция «по оценке N» — пока без разделителя) |

---

## Скриншоты — что есть и зачем

Все скрины — в [`screenshots/`](./screenshots/).

### Мобильные (390 × ?, full-page, без BottomTabBar для чистоты)

| Файл | Экран | Что смотреть | Связанные доки |
|---|---|---|---|
| [`mobile-01-login.png`](./screenshots/mobile-01-login.png) | `/login` | бумажная карточка с tilt −0.6°, paper-control, терракотовая CTA, рукописный хинт | DESIGN-DECISIONS §A1/G1, STATUS R3 |
| [`mobile-02-invite-invalid.png`](./screenshots/mobile-02-invite-invalid.png) | `/invite/<bad>` | bad-state: штамп «недействительно», красно-серая полоса, CTA «← к входу» | STATUS R3 |
| [`mobile-03-board.png`](./screenshots/mobile-03-board.png) | `/` Доска | wordmark, 2-кол grid с featured, шапка недели + PinButton, свернутые архивные полоски с миниатюрами, «↓ рашьше — в архив» | STATUS R1+R4 A1/A3, OPEN-QUESTIONS R5-Q2 |
| [`mobile-04-find.png`](./screenshots/mobile-04-find.png) | `/places` Найти | поиск-штамп, filter-pin, чипы фильтров, полки Жемчужины / По городам / По кухням / По друзьям | STATUS R2/R4 F1, DESIGN-DECISIONS §F1 |
| [`mobile-05-map.png`](./screenshots/mobile-05-map.png) | `/map` Карта | Yandex карта со скрапбук-маркерами (канцелярка + ромбик-жемчужина), стат «N из M» | STATUS R3, DESIGN-DECISIONS §M1/M2 |
| [`mobile-06-profile.png`](./screenshots/mobile-06-profile.png) | `/profile` Я | аватарка, серифа-имя, билетик-стата (144/14/3), «любит европейскую — 2 раза», табы Визиты/Wishlist/Записки/Настройки в одну строку, лента ReviewCard | DESIGN-DECISIONS Q8 (favorite cuisine) |
| [`mobile-07-gems.png`](./screenshots/mobile-07-gems.png) | `/gems` Hub жемчужин | секции «По городам / Кто отмечал / Все жемчужины» | STATUS R4 B4 |
| [`mobile-08-city.png`](./screenshots/mobile-08-city.png) | `/cities/Казань` | имя города серифой, билетик-стата, секции Жемчужины + Все места | STATUS R4 B3 |
| [`mobile-09-person.png`](./screenshots/mobile-09-person.png) | `/people/3` (charlie, 4 места) | публичный профиль друга — заголовок, билетик, города, визиты | STATUS R4 B3 |
| [`mobile-10-place-gem.png`](./screenshots/mobile-10-place-gem.png) | `/places/185` Мясной гуру | обложка-полароид, серифа-имя, штампики (город/кухня/категории/жемчужина), общий рейтинг 8.8, ряд аватарок «×2», карта-превью, два ReviewCard (один с photo + ticket + штамп «жемчужина») | DESIGN-DECISIONS Q2 (extended `/api/places/:id`) |
| [`mobile-11-place-video.png`](./screenshots/mobile-11-place-video.png) | `/places/183` Blanche | место с видеосообщением, большой кружок с ▶ внутри ReviewCard | DESIGN-DECISIONS R5 (video-poster + iOS), OPEN-QUESTIONS R5-Q5 |
| [`mobile-12-place-new.png`](./screenshots/mobile-12-place-new.png) | `/places/new` | бумажная обёртка, шапка «Новое место», поисковая строка Yandex + карта, ручной ввод, тип кухни/категория, dropzone «Фото заведения», чекбокс «Хочу сходить», CTA «создать/отмена» | STATUS R2 (PlaceForm), DESIGN-DECISIONS Flow A |
| [`mobile-13-invites.png`](./screenshots/mobile-13-invites.png) | `/invites` | блок «создать приглашение», список инвайтов как настоящие билетики с пробитыми кружками | STATUS R3 |
| [`mobile-14-public-share.png`](./screenshots/mobile-14-public-share.png) | `/p/185` (logged out) | публичная шеринг-страница — cover full-bleed, бумажная плашка, имя, штамп жемчужины, CTA «войти, чтобы увидеть наши впечатления». Без рейтингов, без авторов (privacy) | DESIGN-DECISIONS Q3, OPEN-QUESTIONS R5-Q3 |

### Мобильные с реальным chrome (BottomTabBar в кадре, viewport-only)

| Файл | Что | Зачем |
|---|---|---|
| [`chrome-board.png`](./screenshots/chrome-board.png) | Доска с активным табом «доска» | как пользователь видит экран при первом открытии |
| [`chrome-find.png`](./screenshots/chrome-find.png) | Найти с активным «найти» | компоновка с фиксированным таб-баром |
| [`chrome-profile.png`](./screenshots/chrome-profile.png) | Профиль с активным «Я» | проверка, что табы внутри (Визиты/Wishlist/Записки/Настройки) и таб-бар не пересекаются |

### Десктоп (1280 × ?, viewport-only)

| Файл | Экран | Что показывает |
|---|---|---|
| [`desktop-01-login.png`](./screenshots/desktop-01-login.png) | `/login` | центрированная узкая колонка на десктопе, paper-edge фон по бокам (G1) |
| [`desktop-02-board.png`](./screenshots/desktop-02-board.png) | `/` | Доска центрирована max-width 640px, paper-edge как «лист бумаги на столе» |
| [`desktop-03-find.png`](./screenshots/desktop-03-find.png) | `/places` | те же полки, центрирование, БЕЗ BottomTabBar (он только мобильный) |
| [`desktop-04-place-gem.png`](./screenshots/desktop-04-place-gem.png) | `/places/185` Мясной гуру | карточка места на широком экране, обложка-полароид, билетик, отзывы |

---

## Что фронт делает не как ожидалось — известные ограничения

- **`/api/users/me`** возвращает 404 «invalid user id» — `me` не резолвится на бэке. Профиль использует `auth.user.id` напрямую, проблема не блокирующая (наблюдается, если кликнуть по ссылке, ведущей на `/api/users/me`)
- **PersonPage без пагинации** — у пользователя со 150+ мест страница уходит на ~32k px высоты. Не баг рендера, но UX-ограничение (см. `mobile-09-person.png` для лёгкого случая с 4 местами)
- **`/api/places` в Найти на больших коллекциях** — пагинация есть, но для дальних полок (Жемчужины / По городам / По друзьям) клиент пока вычисляет из текущей страницы

---

## Как мне скинуть правки

1. Тапни в `OPEN-QUESTIONS.md` соответствующий R5-Qx или назови файл скриншота — этого достаточно, чтобы найти точку в коде
2. Если правка визуальная (отступ / цвет / размер шрифта) — назови класс из CSS-переменных (`--sb-terracotta`, `--sb-paper-card`, `--sb-ink-mute`) или примитив (`Polaroid`, `Ticket`, `Stamp`…)
3. Если правка флоу — лучше сначала договоримся в чате, потом фиксируем в `OPEN-QUESTIONS.md` как новый Q

---

## Технические заметки на полях

- **Шрифты**: Lora (serif) + Caveat (handwritten). Fallback `Georgia + cursive`
- **Палитра**: OKLCH в `frontend/src/assets/scss/scrapbook.scss`, переменные `--sb-*`
- **Tilt**: только 1–5°, никогда больше. Полароиды + записки обязательно с тейпом
- **Анимации**: только CSS, gem-gleam ≈ раз в 5 сек, артефакт-mount 280 мс, разворот недели 320 мс
- **Тёмная тема**: подготовлена (цвета через переменные, `color-scheme: light` зафиксирован) — концептуально not now, ~3 мес после запуска (G2)
- **Mobile-first**: всё рассчитано на 360–430px. Десктоп центрирует узкой колонкой
