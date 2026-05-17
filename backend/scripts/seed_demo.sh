#!/usr/bin/env bash
# seed_demo.sh — заливает демо-данные для скринов и дизайнерского контроля.
#
# Все объекты помечены через username LIKE 'seed_%'. Удаление одним
# `seed_demo_down.sh` (DELETE FROM users WHERE username LIKE 'seed_%' —
# каскадно подметёт notes/wishlist; DELETE FROM places WHERE created_by
# в seed-юзерах — каскадно подметёт reviews/photos).
#
# Реальные данные (lamark/alina/charlie и их places) НЕ затрагиваются.
#
# Юзеры:    seed_anna, seed_petr, seed_olga, seed_max, seed_kate (пароль demo12345)
# Города:   Москва, Санкт-Петербург, Самара (не пересекаются с реальными
#           Казань/Нижний Новгород/Ижевск)
# Места:    30 разнообразных (с фото / без / с цитатами / gem / wishlist)
# Фото:     15 шт., скачаны с picsum.photos, лежат как /app/uploads/seed_*.jpg

set -euo pipefail

PG_CONTAINER=${PG_CONTAINER:-aeva-postgres}
BE_CONTAINER=${BE_CONTAINER:-aeva-backend}
DB=${DB:-aeva_eat}
DB_USER=${DB_USER:-aeva}

echo "→ checking containers"
docker exec "$PG_CONTAINER" pg_isready -U "$DB_USER" -d "$DB" > /dev/null
docker exec "$BE_CONTAINER" sh -c 'test -d /app/uploads' > /dev/null

echo "→ downloading 15 seed photos from picsum.photos"
TMP_DIR=$(mktemp -d)
trap "rm -rf $TMP_DIR" EXIT
for i in $(seq 1 15); do
  if ! curl -fsSL "https://picsum.photos/seed/aeva-eat-$i/900/900.jpg" -o "$TMP_DIR/seed_${i}.jpg"; then
    echo "  ! failed to download photo $i — aborting"
    exit 1
  fi
done

echo "→ copying photos to $BE_CONTAINER:/app/uploads/"
for i in $(seq 1 15); do
  docker cp "$TMP_DIR/seed_${i}.jpg" "$BE_CONTAINER:/app/uploads/seed_${i}.jpg" > /dev/null
done

# bcrypt-хэш для demo12345 — фиксированный, чтобы не зависеть от htpasswd
# (генерировался один раз: htpasswd -bnBC 10 "" demo12345 | tr -d ':\n').
DEMO_HASH='$2y$10$Mjg7QHK9hKQDsW5tQzgkPe0bxqZHrU9TQ6dKqMhmRDIqA4qfQbm9.'
# Использую готовый хэш ниже через psql DO-block.

echo "→ seeding DB"
docker exec -i "$PG_CONTAINER" psql -U "$DB_USER" -d "$DB" -v ON_ERROR_STOP=1 <<'SQL'
BEGIN;

-- 1) Юзеры (5 шт). password_hash = bcrypt('demo12345').
INSERT INTO users (username, display_name, password_hash, role)
VALUES
  ('seed_anna', 'Анна', '$2a$10$j9WOEO.XaPksVrtG41vMduGIEpQwm3N7q/I5VpuKHhD00TaNbdck2', 'user'),
  ('seed_petr', 'Пётр', '$2a$10$j9WOEO.XaPksVrtG41vMduGIEpQwm3N7q/I5VpuKHhD00TaNbdck2', 'user'),
  ('seed_olga', 'Ольга', '$2a$10$j9WOEO.XaPksVrtG41vMduGIEpQwm3N7q/I5VpuKHhD00TaNbdck2', 'user'),
  ('seed_max',  'Макс', '$2a$10$j9WOEO.XaPksVrtG41vMduGIEpQwm3N7q/I5VpuKHhD00TaNbdck2', 'user'),
  ('seed_kate', 'Катя', '$2a$10$j9WOEO.XaPksVrtG41vMduGIEpQwm3N7q/I5VpuKHhD00TaNbdck2', 'user')
ON CONFLICT (username) DO NOTHING;

-- 2) Места. 30 шт в 3 городах. Половина с image_url, половина без.
WITH seed_uids AS (
  SELECT id, username FROM users WHERE username LIKE 'seed_%'
)
INSERT INTO places (name, address, city, lat, lng, cuisine_type_id, image_url, created_by)
SELECT
  p.name, p.address, p.city, p.lat, p.lng, p.ct, p.img,
  (SELECT id FROM seed_uids WHERE username = p.who)
FROM (VALUES
  -- Москва (12 мест)
  ('Cutfish',         'ул. Большая Дмитровка, 12',  'Москва',  55.7616, 37.6151,   1, '/uploads/seed_1.jpg', 'seed_anna'),
  ('Mio',             'ул. Пятницкая, 16',          'Москва',  55.7456, 37.6275,   2, '/uploads/seed_2.jpg', 'seed_petr'),
  ('Sage',            'Климентовский пер., 14',     'Москва',  55.7430, 37.6262, 376, NULL,                  'seed_olga'),
  ('Восход',          'Лужниковская наб., 2/4',     'Москва',  55.7155, 37.5512,   4, '/uploads/seed_3.jpg', 'seed_anna'),
  ('Северяне',        'Большая Никитская, 12А',     'Москва',  55.7572, 37.6051,   4, NULL,                  'seed_max'),
  ('Twins Garden',    'Страстной бульв., 8А',       'Москва',  55.7651, 37.6111, 377, '/uploads/seed_4.jpg', 'seed_kate'),
  ('Selfie',          'Новинский бульв., 31',       'Москва',  55.7574, 37.5854, 376, NULL,                  'seed_anna'),
  ('Lavkalavka',      'Петровка, 21',               'Москва',  55.7639, 37.6151,   4, '/uploads/seed_5.jpg', 'seed_olga'),
  ('Уголёк',          'Большая Никитская, 12',      'Москва',  55.7573, 37.6053,   4, NULL,                  'seed_petr'),
  ('Magadan',         'Кузнецкий мост, 4/3',        'Москва',  55.7616, 37.6224, 377, '/uploads/seed_6.jpg', 'seed_max'),
  ('Cafe Pushkin',    'Тверской бульв., 26А',       'Москва',  55.7637, 37.6058,   5, NULL,                  'seed_kate'),
  ('Pinch',           'Большая Никитская, 23/14/1', 'Москва',  55.7580, 37.6028, 376, '/uploads/seed_7.jpg', 'seed_anna'),

  -- Санкт-Петербург (12 мест)
  ('Hamlet+Jacks',    'ул. Восстания, 26',          'Санкт-Петербург', 59.9381, 30.3631,  11, '/uploads/seed_8.jpg', 'seed_petr'),
  ('Bellevue',        'наб. Адмиралтейская, 10',    'Санкт-Петербург', 59.9352, 30.3081,   5, NULL,                  'seed_olga'),
  ('Mansarda',        'Почтамтская ул., 3-5',       'Санкт-Петербург', 59.9333, 30.3035, 376, '/uploads/seed_9.jpg', 'seed_anna'),
  ('Banshiki',        'ул. Маяковского, 5',         'Санкт-Петербург', 59.9385, 30.3543,   4, NULL,                  'seed_max'),
  ('Subzero',         'Потёмкинская, 4',            'Санкт-Петербург', 59.9461, 30.3739, 377, '/uploads/seed_10.jpg','seed_kate'),
  ('Recolte',         'ул. Гороховая, 13',          'Санкт-Петербург', 59.9320, 30.3157,   5, NULL,                  'seed_anna'),
  ('Cococo',          'ул. Некрасова, 8',           'Санкт-Петербург', 59.9377, 30.3554,   4, '/uploads/seed_11.jpg','seed_petr'),
  ('Block',           'Потёмкинская, 4',            'Санкт-Петербург', 59.9456, 30.3742,  11, NULL,                  'seed_olga'),
  ('Tartar Bar',      'Гороховая, 13',              'Санкт-Петербург', 59.9311, 30.3164,   5, '/uploads/seed_12.jpg','seed_max'),
  ('Animals',         'Невский пр., 92',            'Санкт-Петербург', 59.9358, 30.3550,  11, NULL,                  'seed_kate'),
  ('Cure',            'Гороховая, 13',              'Санкт-Петербург', 59.9314, 30.3160, 377, '/uploads/seed_13.jpg','seed_anna'),
  ('Bonsai',          'наб. Фонтанки, 32',          'Санкт-Петербург', 59.9362, 30.3471,   2, NULL,                  'seed_petr'),

  -- Самара (6 мест)
  ('Парк-Кафе',       'ул. Молодогвардейская, 196', 'Самара', 53.1959, 50.0934,  15, '/uploads/seed_14.jpg', 'seed_olga'),
  ('Сёмки',           'ул. Куйбышева, 91',          'Самара', 53.1880, 50.0995,   4, NULL,                   'seed_max'),
  ('Винни',           'ул. Самарская, 73',          'Самара', 53.1944, 50.1067,  12, '/uploads/seed_15.jpg', 'seed_kate'),
  ('Mole',            'ул. Ленинградская, 77',      'Самара', 53.1925, 50.0950,  11, NULL,                   'seed_anna'),
  ('Pacharan',        'ул. Молодогвардейская, 80',  'Самара', 53.1953, 50.0989, 376, NULL,                   'seed_petr'),
  ('Mercato',         'Ленинградская, 36',          'Самара', 53.1898, 50.0934,   1, NULL,                   'seed_olga')
) AS p(name, address, city, lat, lng, ct, img, who)
ON CONFLICT DO NOTHING;

-- 3) Reviews. Используем cross join seed_places × seed_users чтобы каждое
-- место получило 1-3 отзыва от разных авторов. Часть с длинными комментами
-- (для Q-layout PhotoFreeCard и цитаты от круга в B1). Часть с is_gem.
WITH
  sp AS (SELECT id, name, created_by FROM places WHERE created_by IN (SELECT id FROM users WHERE username LIKE 'seed_%')),
  su AS (SELECT id, username FROM users WHERE username LIKE 'seed_%' ORDER BY id),
  ranked AS (
    SELECT
      sp.id AS place_id,
      su.id AS user_id,
      su.username AS author,
      sp.name AS place_name,
      ROW_NUMBER() OVER (PARTITION BY sp.id ORDER BY (sp.id * 13 + su.id * 7) % 100) AS rn,
      (sp.id * 31 + su.id) % 100 AS h
    FROM sp CROSS JOIN su
  ),
  picks AS (
    -- Берём 1-3 отзыва на место (rn <= 1 + (place_id % 3))
    SELECT * FROM ranked WHERE rn <= 1 + (place_id % 3)
  )
INSERT INTO reviews (place_id, food_rating, service_rating, vibe_rating, is_gem, comment, visited_at, image_url, created_at)
SELECT
  place_id,
  -- рейтинги 6.5..9.5 со случайной вариацией
  6.5 + ((h * 7) % 30) / 10.0,
  6.5 + ((h * 11) % 30) / 10.0,
  6.5 + ((h * 13) % 30) / 10.0,
  -- ~30% жемчужин (h % 10 < 3)
  (h % 10 < 3),
  -- ~70% с длинным комментом (h % 10 < 7), варьируем 4 шаблона по h
  CASE WHEN h % 10 < 7 THEN
    CASE h % 4
      WHEN 0 THEN 'Сидели у окна, обслуживание расслабленное, кухня уверенная. Хотим вернуться зимой на ужин-сюрприз.'
      WHEN 1 THEN 'Заказали дегустацию. Главный хит — местная рыба с пюре из топинамбура. Десерт-провал, но это мелочь.'
      WHEN 2 THEN 'Очень милое место, тихо. Хлеб домашний, вино подбирают с разбором. Уйти не могли часа три.'
      ELSE         'Зашли случайно — оказались здесь на два часа. Винная карта камерная, но всё попадает. Тёплая встреча.'
    END
  ELSE NULL END,
  -- visited_at: разбрасываем по последним 6 неделям
  CURRENT_DATE - ((h % 42) || ' days')::interval,
  -- image_url: на 25% review-фото из пула seed_*.jpg
  CASE WHEN h % 4 = 0 THEN '/uploads/seed_' || (1 + (h % 15)) || '.jpg' ELSE NULL END,
  -- created_at в той же зоне для визуальной кучи
  now() - ((h % 42) || ' days')::interval - ((h * 3 % 24) || ' hours')::interval
FROM picks;

-- 4) review_authors — для каждого review проставляем автора (из той же тройки).
INSERT INTO review_authors (review_id, user_id)
SELECT r.id, p.created_by
FROM reviews r
JOIN places p ON p.id = r.place_id
WHERE p.created_by IN (SELECT id FROM users WHERE username LIKE 'seed_%')
ON CONFLICT DO NOTHING;

-- 5) review_photos — для тех reviews у кого image_url не пуст, дублируем
-- в review_photos (это то что фронт PolaroidStack пулит как стопку).
INSERT INTO review_photos (review_id, url, position)
SELECT r.id, r.image_url, 0
FROM reviews r
JOIN places p ON p.id = r.place_id
WHERE p.created_by IN (SELECT id FROM users WHERE username LIKE 'seed_%')
  AND r.image_url IS NOT NULL;

-- 6) wishlist_custom: пара плановок у seed_anna
INSERT INTO wishlist_custom (user_id, name, note)
SELECT u.id, v.name, v.note
FROM users u
CROSS JOIN (VALUES
  ('Сёмки',    'обещали попробовать манты'),
  ('Bonsai',   'давно зову Олю')
) AS v(name, note)
WHERE u.username = 'seed_anna';

-- 7) Записка от seed_petr на Доске (без места)
INSERT INTO notes (author_id, text, paper_color, tape_color)
SELECT u.id, 'забронировали Pinch на пятницу — пишите кто хочет', 'cream', 'rose'
FROM users u WHERE u.username = 'seed_petr';

COMMIT;

\echo
\echo '=== seed summary ==='
SELECT
  (SELECT COUNT(*) FROM users WHERE username LIKE 'seed_%')                AS seed_users,
  (SELECT COUNT(*) FROM places WHERE created_by IN
    (SELECT id FROM users WHERE username LIKE 'seed_%'))                   AS seed_places,
  (SELECT COUNT(*) FROM reviews r JOIN places p ON p.id=r.place_id
    WHERE p.created_by IN (SELECT id FROM users WHERE username LIKE 'seed_%')) AS seed_reviews,
  (SELECT COUNT(*) FROM review_photos rp JOIN reviews r ON r.id=rp.review_id
    JOIN places p ON p.id=r.place_id WHERE p.created_by IN
    (SELECT id FROM users WHERE username LIKE 'seed_%'))                   AS seed_photos,
  (SELECT COUNT(*) FROM notes WHERE author_id IN
    (SELECT id FROM users WHERE username LIKE 'seed_%'))                   AS seed_notes,
  (SELECT COUNT(*) FROM wishlist_custom WHERE user_id IN
    (SELECT id FROM users WHERE username LIKE 'seed_%'))                   AS seed_wishlist;
SQL

echo
echo "✓ done. Reload http://localhost:8091/cities/Москва to see seeded data."
echo "  Login: any seed_anna / seed_petr / … with password demo12345"
echo "  Rollback:  bash backend/scripts/seed_demo_down.sh"
