-- 001_init.up.sql

-- Users
CREATE TABLE IF NOT EXISTS users (
    id          SERIAL PRIMARY KEY,
    username    VARCHAR(50)  NOT NULL UNIQUE,
    email       VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT       NOT NULL,
    avatar_url  TEXT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);

-- Cuisine types
CREATE TABLE IF NOT EXISTS cuisine_types (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

-- Categories (завтрак, обед, ужин, бар, кофейня...)
CREATE TABLE IF NOT EXISTS categories (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

-- Places
CREATE TABLE IF NOT EXISTS places (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    address         TEXT,
    city            VARCHAR(100),
    lat             DOUBLE PRECISION,
    lng             DOUBLE PRECISION,
    cuisine_type_id INT REFERENCES cuisine_types(id) ON DELETE SET NULL,
    website         TEXT,
    created_by      INT REFERENCES users(id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_places_city ON places(city);
CREATE INDEX IF NOT EXISTS idx_places_cuisine ON places(cuisine_type_id);

-- Place ↔ Categories (M2M)
CREATE TABLE IF NOT EXISTS place_categories (
    place_id    INT NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    category_id INT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (place_id, category_id)
);

-- Reviews (multiple per user+place allowed, differentiated by visited_at)
CREATE TABLE IF NOT EXISTS reviews (
    id             SERIAL PRIMARY KEY,
    place_id       INT  NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    food_rating    INT  NOT NULL CHECK (food_rating BETWEEN 1 AND 10),
    service_rating INT  NOT NULL CHECK (service_rating BETWEEN 1 AND 10),
    vibe_rating    INT  NOT NULL CHECK (vibe_rating BETWEEN 1 AND 10),
    is_gem         BOOLEAN NOT NULL DEFAULT false,
    comment        TEXT,
    visited_at     DATE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_reviews_place ON reviews(place_id);

-- Review authors (M2M: review ↔ users)
CREATE TABLE IF NOT EXISTS review_authors (
    review_id INT NOT NULL REFERENCES reviews(id) ON DELETE CASCADE,
    user_id   INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (review_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_review_authors_user ON review_authors(user_id);

-- Seed: cuisine types
INSERT INTO cuisine_types (name) VALUES
    ('Итальянская'), ('Японская'), ('Грузинская'), ('Русская'),
    ('Французская'), ('Мексиканская'), ('Китайская'), ('Корейская'),
    ('Тайская'), ('Индийская'), ('Американская'), ('Средиземноморская'),
    ('Узбекская'), ('Турецкая'), ('Смешанная')
ON CONFLICT (name) DO NOTHING;

-- Seed: categories
INSERT INTO categories (name) VALUES
    ('Завтрак'), ('Обед'), ('Ужин'), ('Бар'), ('Кофейня'),
    ('Фастфуд'), ('Кондитерская'), ('Пекарня'), ('Стритфуд'), ('Банкет')
ON CONFLICT (name) DO NOTHING;
