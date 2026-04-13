-- 005_seed_data.up.sql — Demo seed data with photos, joint reviews, gems, wishlists

-- Users (passwords are bcrypt hash of "password123")
INSERT INTO users (id, username, email, password_hash) VALUES
  (1, 'alice', 'alice@aeva.eat', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy'),
  (2, 'bob', 'bob@aeva.eat', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy'),
  (3, 'charlie', 'charlie@aeva.eat', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy')
ON CONFLICT (id) DO NOTHING;

SELECT setval('users_id_seq', GREATEST((SELECT MAX(id) FROM users), 3));

-- Places
INSERT INTO places (id, name, address, city, lat, lng, cuisine_type_id, website, created_by, image_url) VALUES
  (1, 'Тоскана',           'ул. Пятницкая, 12',       'Москва',           55.7420, 37.6289, 1,  'https://toscana.ru',    1, 'https://images.unsplash.com/photo-1517248135467-4c7edcad34c4?w=600&h=400&fit=crop'),
  (2, 'Sakura',            'ул. Мясницкая, 24',       'Москва',           55.7620, 37.6380, 2,  NULL,                    1, 'https://images.unsplash.com/photo-1579871494447-9811cf80d66c?w=600&h=400&fit=crop'),
  (3, 'Хинкальная №1',     'ул. Садовая, 5',          'Москва',           55.7655, 37.5930, 3,  NULL,                    2, 'https://images.unsplash.com/photo-1555396273-367ea4eb4db5?w=600&h=400&fit=crop'),
  (4, 'Le Petit Bistro',   'Невский пр., 40',         'Санкт-Петербург',  59.9340, 30.3270, 5,  'https://lepetit.ru',    2, 'https://images.unsplash.com/photo-1414235077428-338989a2e8c0?w=600&h=400&fit=crop'),
  (5, 'Coffee Soul',       'ул. Рубинштейна, 10',     'Санкт-Петербург',  59.9300, 30.3440, NULL, NULL,                  3, 'https://images.unsplash.com/photo-1495474472287-4d71bcdd2085?w=600&h=400&fit=crop'),
  (6, 'Burger Heroes',     'ул. Большая Дмитровка, 9','Москва',           55.7630, 37.6130, 11, NULL,                    1, 'https://images.unsplash.com/photo-1568901346375-23c9450c58cd?w=600&h=400&fit=crop'),
  (7, 'Плов и Лагман',     'ул. Маросейка, 7',        'Москва',           55.7580, 37.6390, 13, NULL,                    3, NULL),
  (8, 'Thai Orchid',       'ул. Покровка, 16',        'Москва',           55.7590, 37.6470, 9,  'https://thaiorchid.ru', 2, 'https://images.unsplash.com/photo-1562565652-a0d8f0c59eb4?w=600&h=400&fit=crop')
ON CONFLICT (id) DO NOTHING;

SELECT setval('places_id_seq', GREATEST((SELECT MAX(id) FROM places), 8));

-- Place categories
INSERT INTO place_categories (place_id, category_id) VALUES
  (1, 3), (1, 10),   -- Тоскана: Ужин, Банкет
  (2, 3), (2, 2),    -- Sakura: Ужин, Обед
  (3, 2), (3, 3),    -- Хинкальная: Обед, Ужин
  (4, 3),            -- Le Petit: Ужин
  (5, 5), (5, 1),    -- Coffee Soul: Кофейня, Завтрак
  (6, 6), (6, 2),    -- Burger Heroes: Фастфуд, Обед
  (7, 2),            -- Плов и Лагман: Обед
  (8, 3), (8, 2)     -- Thai Orchid: Ужин, Обед
ON CONFLICT DO NOTHING;

-- Reviews (ratings are NUMERIC 0-10, step 0.1)
-- Тоскана: 2 reviews (one joint alice+bob)
INSERT INTO reviews (id, place_id, food_rating, service_rating, vibe_rating, is_gem, comment, visited_at, image_url) VALUES
  (1, 1, 9.2, 8.5, 9.0, true,   'Невероятная паста карбонара! Атмосфера уютная, как в Италии. Обязательно вернёмся.', '2026-03-15', 'https://images.unsplash.com/photo-1621996346565-e3dbc646d9a9?w=600&h=400&fit=crop'),
  (2, 1, 8.0, 7.5, 8.5, false,  'Хорошая пицца, но чуть долго ждали. Тирамису — отличное.', '2026-04-01', NULL)
ON CONFLICT (id) DO NOTHING;

-- Sakura
INSERT INTO reviews (id, place_id, food_rating, service_rating, vibe_rating, is_gem, comment, visited_at, image_url) VALUES
  (3, 2, 8.8, 9.0, 7.5, true, 'Лучшие суши в городе. Шеф-повар из Осаки, это чувствуется в каждом ролле.', '2026-02-20', 'https://images.unsplash.com/photo-1579584425555-c3ce17fd4351?w=600&h=400&fit=crop'),
  (4, 2, 7.0, 8.0, 6.5, false, 'Сашими свежайшее, но маловато вариантов в меню.', '2026-03-10', NULL)
ON CONFLICT (id) DO NOTHING;

-- Хинкальная
INSERT INTO reviews (id, place_id, food_rating, service_rating, vibe_rating, is_gem, comment, visited_at, image_url) VALUES
  (5, 3, 9.5, 7.0, 8.0, true, 'Хинкали как в Тбилиси! Сочные, бульон внутри — огонь. Хачапури тоже 🔥', '2026-04-05', 'https://images.unsplash.com/photo-1565299624946-b28f40a0ae38?w=600&h=400&fit=crop')
ON CONFLICT (id) DO NOTHING;

-- Le Petit Bistro (joint review by bob+charlie)
INSERT INTO reviews (id, place_id, food_rating, service_rating, vibe_rating, is_gem, comment, visited_at, image_url) VALUES
  (6, 4, 9.0, 9.5, 9.8, true, 'Французская классика на высшем уровне. Утиная грудка, крем-брюле — шедевр. Сомелье подобрал идеальное вино.', '2026-03-28', 'https://images.unsplash.com/photo-1544025162-d76694265947?w=600&h=400&fit=crop')
ON CONFLICT (id) DO NOTHING;

-- Coffee Soul
INSERT INTO reviews (id, place_id, food_rating, service_rating, vibe_rating, is_gem, comment, visited_at, image_url) VALUES
  (7, 5, 8.0, 8.5, 9.0, false, 'Отличный specialty coffee, красивая подача. Авокадо-тост свежий, но порция маленькая.', '2026-04-02', 'https://images.unsplash.com/photo-1495474472287-4d71bcdd2085?w=600&h=400&fit=crop'),
  (8, 5, 7.5, 9.0, 9.5, true,  'Самая уютная кофейня в Питере. Латте-арт от бариста — произведение искусства.', '2026-03-20', NULL)
ON CONFLICT (id) DO NOTHING;

-- Burger Heroes
INSERT INTO reviews (id, place_id, food_rating, service_rating, vibe_rating, is_gem, comment, visited_at, image_url) VALUES
  (9, 6, 7.5, 6.0, 7.0, false, 'Мощные бургеры, хорошая картошка. Шумновато, но для фастфуда — ок.', '2026-04-08', 'https://images.unsplash.com/photo-1568901346375-23c9450c58cd?w=600&h=400&fit=crop')
ON CONFLICT (id) DO NOTHING;

-- Thai Orchid (joint review by alice+charlie)
INSERT INTO reviews (id, place_id, food_rating, service_rating, vibe_rating, is_gem, comment, visited_at, image_url) VALUES
  (10, 8, 8.5, 8.0, 8.5, false, 'Настоящий тайский том-ям! Pad Thai хороший, но чуть не хватает остроты. Зато амбиент потрясающий.', '2026-04-10', NULL)
ON CONFLICT (id) DO NOTHING;

SELECT setval('reviews_id_seq', GREATEST((SELECT MAX(id) FROM reviews), 10));

-- Update review photos for seed data (in case reviews existed before image_url column)
UPDATE reviews SET image_url = 'https://images.unsplash.com/photo-1621996346565-e3dbc646d9a9?w=600&h=400&fit=crop' WHERE id = 1 AND image_url IS NULL;
UPDATE reviews SET image_url = 'https://images.unsplash.com/photo-1579584425555-c3ce17fd4351?w=600&h=400&fit=crop' WHERE id = 3 AND image_url IS NULL;
UPDATE reviews SET image_url = 'https://images.unsplash.com/photo-1565299624946-b28f40a0ae38?w=600&h=400&fit=crop' WHERE id = 5 AND image_url IS NULL;
UPDATE reviews SET image_url = 'https://images.unsplash.com/photo-1544025162-d76694265947?w=600&h=400&fit=crop' WHERE id = 6 AND image_url IS NULL;
UPDATE reviews SET image_url = 'https://images.unsplash.com/photo-1495474472287-4d71bcdd2085?w=600&h=400&fit=crop' WHERE id = 7 AND image_url IS NULL;
UPDATE reviews SET image_url = 'https://images.unsplash.com/photo-1568901346375-23c9450c58cd?w=600&h=400&fit=crop' WHERE id = 9 AND image_url IS NULL;

-- Review authors (M2M)
INSERT INTO review_authors (review_id, user_id) VALUES
  -- Тоскана review 1: joint (alice + bob)
  (1, 1), (1, 2),
  -- Тоскана review 2: alice solo
  (2, 1),
  -- Sakura review 3: bob
  (3, 2),
  -- Sakura review 4: charlie
  (4, 3),
  -- Хинкальная review 5: alice
  (5, 1),
  -- Le Petit review 6: joint (bob + charlie)
  (6, 2), (6, 3),
  -- Coffee Soul review 7: charlie
  (7, 3),
  -- Coffee Soul review 8: alice
  (8, 1),
  -- Burger Heroes review 9: bob
  (9, 2),
  -- Thai Orchid review 10: joint (alice + charlie)
  (10, 1), (10, 3)
ON CONFLICT DO NOTHING;

-- Wishlists (planned places)
INSERT INTO wishlists (user_id, place_id) VALUES
  (1, 4),  -- alice хочет в Le Petit Bistro
  (1, 7),  -- alice хочет Плов и Лагман
  (2, 8),  -- bob хочет Thai Orchid
  (3, 1),  -- charlie хочет в Тоскану
  (3, 6)   -- charlie хочет Burger Heroes
ON CONFLICT DO NOTHING;

-- Custom wishlist entries (free-text)
INSERT INTO wishlist_custom (user_id, name, note) VALUES
  (1, 'White Rabbit', 'Ресторан на крыше с видом на Москву, бронь за месяц'),
  (1, 'Kuznya House', 'СПб, на Новой Голландии, советовали друзья'),
  (2, 'Техникум', 'Бар с коктейлями на Бауманской'),
  (3, 'Selfie', 'Мишленовский ресторан, хочу попробовать дегустационный сет')
ON CONFLICT DO NOTHING;
