# Следующая сессия — начни тут

Этот файл — пятиминутный onboarding для следующего захода в проект. Подразумевает, что предыдущая сессия закрыта, контекст сброшен.

## TL;DR (статус на 2026-05-17)

**Раунд R6 закрыт почти весь.** Все 8 болевых точек дизайнера (A1-A3, B1-B8) и 6 R5-Q закрыты. См. таблицу ниже.

**Что осталось:**
- **Структурные D1-D4** — рефактор `ArtifactCard`, правило priority full-width, dense-flow, kruzhok-miss телеметрия. Не визуальные, не дрейф — архитектурные опасения.
- **Юзерские баги «потом расскажу»** — пользователь упомянул что есть ещё мелкие баги, не уточнил какие.
- **Не отрисовано дизайнером** — записка на Доске без места, wishlist активный/зачёркнутый, AddArtifactSheet, мини-карта города (B1.2). Ждём следующий v4-канвас.

Полный лог итераций (R6.1-R6.11) — в [`R6_DESIGNER_REVIEW.md`](./R6_DESIGNER_REVIEW.md).

### Статусная таблица R6

| Где | Пункт | Статус |
|---|---|---|
| A1-A3 | пустая Доска, cover-fallback, PhotoFreeCard Q/T/G | ✅ R6.1 |
| B1 | Город как путеводитель | ✅ R6.5 |
| B2 | Жемчужины-сокровищница | ✅ R6.6 |
| B3 | Штамп-помойка в шапке места | ✅ R6.2 |
| B4 / R5-Q1 | подпись без гендер-глагола | ✅ R6.2 |
| B5 / R5-Q3 | Public share + ?next= | ✅ R6.4 |
| B6 | Число 3.5 в ромбе жемчужины | ✅ R6.3 |
| B7 | «любит X» + ячейка городов | ✅ R6.3 |
| B8 | Person reduction (32k → 3.5k px) | ✅ R6.7 |
| R5-Q4 | Мягкий tilt на wishlist | ✅ R6.8 |
| R5-Q5 | «тапни» fallback для video | ✅ R6.8 |
| R5-Q6 | Разделитель «не пробовала» | ✅ R6.8 |
| D5 | Удалить onboarding-cta снизу | ✅ R6.5 |
| **D1** | **Рефактор ArtifactCard на 3 компонента** | ❌ осталось |
| **D2** | **Правило priority `gem > video > первый`** | ❌ осталось |
| **D3** | **`grid-auto-flow: dense` (хронология ленты)** | ❌ осталось |
| **D4** | **Лог `kruzhok-miss` (телеметрия)** | ❌ осталось |
| Нашли сами | Photo-overflow на featured | ✅ R6.5 |
| Нашли сами | Дубли цитат на Доске (per-visit comment) | ✅ R6.10 |
| Нашли сами | Время HH:MM на caption + notes | ✅ R6.10/R6.11 |
| Нашли сами | Фильтр UX consistent + scroll-to-top | ✅ R6.11 |
| Нашли сами | Сидинг demo-данных | ✅ R6.9 |
| Юзер | Bugs которые «потом расскажу» | ⏳ ждём от юзера |

### Ждёт следующий раунд от дизайнера (v4-канвас)

- Записка на Доске без места — как note рендерится рядом с полароидами
- Wishlist-артефакт активный/зачёркнутый — R4-Q1 обсудили текстом, не отрисовано
- AddArtifactSheet — bottom-sheet выбора типа артефакта
- Мини-карта города (B1.2) — bbox по местам города на /cities/:name

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
- ✅ **B3 — штамп-помойка в шапке места.** В `PlaceDetail.vue` шапка: title → адрес → 2 главных штампа (Город + Кухня) + ромб ◆ ЖЕМЧУЖИНА → категории мелким caveat'ом через `·` → подпись жемчужины. Было 7 штампов, стало 3.
- ✅ **B4 (R5-Q1) — подпись без гендер-глагола.** Удалён `gemVerb` (эвристика по `-а/-я` → отметила/отметил). Новый формат: `«жемчужина · Аня · 10 мая + Серёжа + Миша»`. Слово «жемчужина» в moss-цвете, имя без падежного склонения, соавторы через ` + `.
- ✅ **B7 — «любит X» logic-fix.** В `useCuisine.js`: порог ≥5 визитов И ≥15% от `place_count`. У lamark (144 места, любимая кухня 2 раза) фраза скрыта — было «любит европейскую — 2 раза», стало пусто.
- ✅ **B7 — ячейка «городов» в билетике профиля.** `stats.cities` использует `userProfile.city_count` из `/users/:id`, fallback на локальный расчёт. У lamark: «3 ГОРОДОВ» вместо прочерка.
- ✅ **B6 — число на ромбе жемчужины убрано.** `gemSVG` в `MapView.vue` больше не рисует `<text>` рейтинга. Ромб = вердикт «топ», число (особенно 3.5 у Ронни) с ним конфликтует. Рейтинг видно в балуне. Pushpin'ы (134 не-жемчужины) сохранили числа.
- ✅ **B5 + R5-Q3 — public share переработан.** `/p/<id>` теперь: cover ~42vh, paper-плашка с двумя tape, caption «ИЗ ДНЕВНИКА КРУГА», серифа-имя + «город · кухня», ◆ЖЕМЧУЖИНА с ромбом и плашкой, бумажная dashed-CTA с 1° tilt и стрелкой → «войти, чтобы увидеть впечатления» (НЕ терракотовый pill), caveat «камерный дневник еды». Wordmark «AEVA·EAT» сверху-справа. CTA href = `/login?next=/places/{ID}` — `Login.vue` после успешного login делает `router.push(next)` если параметр безопасен (internal path).
- ✅ **Bugfix: фото-overflow на featured карточках.** `.sb-polaroid` без `box-sizing: border-box` + слои PolaroidStack с `inset:0 + width:100%` + translate ±6% уносили featured Мясной гуру за грид-ячейку (374 при ширине 363, body scroll-width 408 vs 390). Fix: `box-sizing: border-box`, `inset:8%` на layer, убран paddingInline rootStyle.
- ✅ **D5 — `.onboarding-cta` снизу Доски удалён.** Один призыв (баннер сверху + PinButton) сильнее двух.
- ✅ **B1 — Город как путеводитель.** `CityPage.vue` переделан по эталону `v3/02-city-guide.png`: серифа-имя + рукописная meta-строка `N мест · M жемчужин · из круга K`, секция «Жемчужины {Город}» горизонтальным shelf из 3 крупных полароидов с tape (наклоны −3°/+2°/−1.5°), бумажная плашка «ЦИТАТА ОТ КРУГА» с самым длинным `top_review_comment` мест города, заголовок «Все N» и компактный список. Мини-карта города пока опущена (требует Yandex bootstrap фильтра).
- ✅ **B2 — Жемчужины-сокровищница.** `GemsHub.vue` переделан по эталону `v3/03-gems-hub.png`: «Жемчужины» крупно с ромбом + подзаголовок `N в M городах`, секция **«САМАЯ ПЕРВАЯ»** с крупным полароидом и историей `username · дата` + `подтвердили: ...`, города чипами (`<Stamp ink> + count + ◆`), ряд аватарок «Кто отмечал», 2-колоночная сетка «Все жемчужины» через `ArtifactCard` (для безфотных — `PhotoFreeCard` G-layout).
- ✅ **B8 — Person page reduction.** Раньше `lamark` (144 места) рендерил **~32000px**. Сейчас **3457px** — в ~10 раз короче. Города — чипами вместо вертикального списка (`КАЗАНЬ 56 · НИЖНИЙ НОВГОРОД 55 · ИЖЕВСК 33`). Жемчужины — 2-колоночная сетка через ArtifactCard. Визиты — лимит 6 по умолчанию + кнопка «↓ ещё N» (раскрывает по 12 за клик). Раньше все 144 рендерились без пагинации.
- ✅ **R5-Q4 — мягкий tilt на wishlist.** В `Home.vue#cellTilt`: для `_kind === 'wishlist'` используется отдельный пул `softTilts = [l1, r1, l2, r2]` — без агрессивных ±3°. План концептуально не «живое впечатление», ему лишний наклон ни к чему.
- ✅ **R5-Q5 — «тапни» fallback на video kruzhok.** В `ArtifactCard.vue#forcePoster` после seek на 0.1s ставится setTimeout 500мс: если `currentTime` всё ещё 0 — добавляется класс `posterless` на `.art-kruzhok-layer`. CSS показывает рукописный плейсхолдер «тапни» поверх тёмной подложки. Скрывается при `.playing`.
- ✅ **R5-Q6 — разделитель «не пробовал(а)».** В `Places.vue` при `sort=rating_user:N` параллельно загружаются `/users/N/places` → `Set` placeId'ов. `splitResults` делит `placesStore.places` на «оценённые другом» и «остальные»; между ними вставляется маркер `{ _divider: true, _friendName, _friendFem }`. Template рендерит плашку «… а вот эти Аня ещё не пробовала» (или «пробовал»). Гендер — эвристикой по `-а/-я`.
- ✅ **Сидинг demo-данных.** `backend/scripts/seed_demo.sh` + `seed_demo_down.sh`. Заливает 5 seed-юзеров (`seed_anna/petr/olga/max/kate`), 30 мест в Москве/СПб/Самаре, 60 reviews (часть с длинными комментами для Q-layout, ~30% gem, 25% с фото), 15 review_photos (с picsum.photos), 1 note, 2 wishlist_custom. Все объекты помечены через `users.username LIKE 'seed_%'` — unseed одной командой каскадно убирает всё, реальные данные (lamark/alina/charlie) не затрагиваются. **Bugfix:** в `CityPage.vue` автор цитаты от круга больше не захардкожен «Серёжа» — берётся из `place.reviewers[0].username`. **Bugfix-2:** seed-комменты теперь из пула 16 шаблонов с per-place солью — дублей одного коммента под одним местом нет.
- ✅ **Bugfix: дубли цитат на Доске.** Два визита Blanche показывали одну и ту же цитату — PhotoFreeCard Q-layout брала `place.top_review_comment` (самый длинный коммент места), не комментарий конкретного визита. Сейчас: backend FeedEvent.Comment + useFeed накапливает самый длинный коммент event'ов группы как per-visit `top_review_comment`. Каждый визит на Доске рендерит свой текст.
- ✅ **UX: фильтр на Найти.** (1) sort, rating_user, attended-чипы внутри дровера больше не применяются автоматически — ждут кнопку «применить» (раньше прыгали при каждом клике). Чипы вне дровера (где можно убрать один фильтр) применяются сразу — это норма. (2) После хоть одного `wasFiltered` сброс последнего фильтра больше не возвращает в полки — результаты остаются. Реальный «перейти к полкам» теперь только через reload или кнопку «к Найти».
- ✅ **UX: время на caption.** `formatVisitCaption` теперь добавляет `HH:MM` если время визита не полуночное. Caption: «Мясной гуру · вс 19:10» вместо просто «Мясной гуру · вс». Если все таймстампы 00:00 (импорт без времени) — оставляем как было.
- ✅ **PhotoFreeCard caption всех layout'ов.** Раньше visit-caption был только в T-layout, в Q (цитата) и G (жемчужина) — пусто. Теперь общий блок `<div class="pfc-cap">` под основным контентом — карточка всегда отвечает «когда были».
- ✅ **NoteArtifact: HH:MM на дате.** `dateLabel` в `NoteArtifact.vue` добавляет время если оно не 00:00. Записки на Доске теперь «17 мая 14:32» вместо «17 мая».
- ✅ **UX фильтра /places — финальная политика.** После двух итераций пришли к **consistent: всё ждёт «применить»**. Юзер жаловался что «по людям применяется сразу, остальное нет — почему?». Теперь все контролы в дровере (chip «Кто был», sort-select, date-presets, multi-selects, gem-switch) только меняют state. Apply строго по кнопке «применить» внизу дровера — она же закрывает drawer через `data-bs-dismiss` и подстраховку через `hidden.bs.offcanvas` listener (чистит body.overflow если Bootstrap зажевал).
- ✅ **Scroll-to-top на навигацию.** В `router/index.js` добавлен `scrollBehavior`: при переходе на новый route (например тап «Найти» в bottom-tab из PlaceDetail) скролл сбрасывается в начало. Browser-back/forward сохраняет позицию через `savedPosition`.
- 📸 Контроль: `screenshots/mobile-09b-person-lamark-top.png` (lamark viewport), `screenshots/mobile-09-person.png` (charlie — 4 места, всё видно), `screenshots/mobile-07b-gems-top.png` (Gems Hub), `screenshots/mobile-08b-city-izhevsk.png` (Город), `screenshots/mobile-10b-place-header.png` (шапка), `screenshots/mobile-06b-profile-header.png` (профиль), `screenshots/mobile-14b-share.png` (share), `screenshots/mobile-03b-board-expanded.png` (доска).

## Что делать в первую очередь (следующий раунд)

Полный список — `R6_DESIGNER_REVIEW.md` → раздел «D. Структурные опасения».

1. **Запросить у юзера его «потом расскажу» баги** — он на /places и Доске видел ещё мелочи в R6.10/R6.11, не уточнил. Уточнить до начала рефакторов.
2. **Рефактор `ArtifactCard.vue` (D1).** Файл 600+ строк, 4 рендер-режима в одном template. Разнести на `ArtifactPolaroid` / `PhotoFreeCard` (уже есть) / `KruzhokStack` + тонкий `ArtifactCard`-роутер сверху. Это рефакторинг без визуальных изменений.
3. **D2 — правило priority full-width.** Сейчас `featured` и `has-video` оба бросают `grid-column: 1/-1` независимо — на одной неделе может оказаться три full-width-карточки подряд. Закрепить иерархию: `gem > video > первый place`, max 1 full-width на неделю. Плюс позволить `ticket-only`/`PhotoFreeCard` тоже становиться featured (сейчас только place с фото).
4. **D3 — `grid-auto-flow: dense` или композиция.** Сейчас dense ломает хронологию: записка от вторника визуально под визитом среды. Варианты: убрать dense + жить с дырой; или закрепить композицию (слева полароид 2/3, справа маленький 1/3).
5. **D4 — лог `kruzhok-miss`.** Телеметрия на промахи по видео-кружочку. Если >10% переходов на `/places/:id` оказались случайными — расширить hit-target.

### Артефакты от дизайнера (v4-канвас, ещё не приходил)

- Записка без места на Доске
- Wishlist активный/зачёркнутый
- AddArtifactSheet bottom-sheet
- Мини-карта города

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
