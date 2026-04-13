#!/usr/bin/env python3
"""Generate SQL seed for all places, reviews, and wishlists."""
import random
random.seed(42)

ALINA_ID = 1
LAMARK_ID = 2

# City centers
CITIES = {
    'Казань': (55.7961, 49.1064),
    'Нижний Новгород': (56.3269, 43.9962),
}

def scatter(city):
    lat, lng = CITIES[city]
    return (round(lat + random.uniform(-0.02, 0.02), 6),
            round(lng + random.uniform(-0.02, 0.02), 6))

def esc(s):
    return s.replace("'", "''")

# ---- DATA ----
kazan_wishlist = [
    "Исфара", "La Casa", "Зум Зум", "Ohana poke", "Dali", "Соседи",
    "Сакура by tasigo", "Black doner", "First gallery kitchen", "Кафе рояль",
    "Брискет", "Грин би", "Тасига", "Cups", "Стерлядка", "Кафе казанское",
    "Пашмир", "Татар by тюбэтэй", "Aulak", "Чирэм", "Виллард", "Pera",
    "Белый аист", "Хорошим людям", "Азия вкусно", "Некрасов",
    "Брассерия левен", "Голодный бык", "Ультрамен", "Эпоха подвоха",
]

kazan_visited = [
    ("Аниме рамен", 5), ("Омномном", 10), ("Food Hanoi", 8.5),
    ("Корчма млин", 2.5), ("Бургер-багет", 1), ("Соль", 10),
    ("KGB", 7), ("Top hop", 7), ("Brooklyn pizza", 6),
    ("Древняя Бухара", 4), ("Грузинские истории", 7), ("У Марико", 8),
    ("У Сосо", 5), ("Cafe LaVitta", 6), ("Paloma", 9),
    ("Тюбетэй", 7.5), ("Май", 7.5), ("Пэпэ", 5), ("Dadli", 3),
    ("Branch", 8), ("ОстроWOK", 5), ("Gnezdo", 2.5),
    ("Black star burger", 2), ("Жарушка", 7.75), ("Избушка", 3),
    ("Dadiani", 7.5), ("Ханума", 9), ("Alanya", 2), ("Сайгон", 9.5),
    ("Begin", 8), ("Брассерия", 7), ("Театральный буфет", 5),
    ("Ща", 7), ("Черновар", 7), ("Green dog", 5), ("Kombinat", 2.5),
    ("Zero", 7), ("More", 10), ("New asia", 7), ("Сянлун", 4),
    ("Релаб", 6), ("Чичети", 8), ("Ichigo ichie", 6.5),
    ("Сыроварня", 10), ("Чайная гуру", 4), ("Manufact", 4),
    ("Nice people", 6.5), ("Pasta bar", 6), ("Lulua", 9),
    ("Gina Italian", 7.5), ("Тануки", 7), ("Миловита", 5),
    ("Industry", 8.5), ("Волна", 6), ("Рецепт", 7),
    ("Итле", 7), ("Эклер", 5), ("Малабар", 6), ("Cheeseria", 4),
    ("Жигули", 9), ("Шашлыкоф", 6), ("Фомин", 8), ("Фидель", 6.5),
    ("Brew barrel", 4), ("Irish pub", 4.5), ("Бинхартс", 3),
    ("Хочу и буду", 7), ("Дринк крафт", 7), ("Чайхана 1", 3.5),
    ("Choodoo", 9), ("Смородина", 7.5), ("Ранняя пташка", 7),
    ("Агафредо", 4), ("Roofbar", 7), ("So sweet", 5.5),
    ("Olio", 8.5), ("Утро", 7.5), ("Босфор", 7), ("Фильтр", 5),
    ("Terra et Silva", 8), ("Pahlava", 8.75), ("Cho", 7),
    ("Саяр", 7.5), ("Cream coffee", 6), ("Ист", 4),
    ("Раковарня", 6.5),
]

nn_wishlist = ["Union jack", "Верхушка", "Hustler"]

nn_visited = [
    ("9/18", 6.5), ("Медные трубы", 4), ("Mola mola", 5),
    ("Маджонг", 6), ("Английское посольство", 8.5), ("Salut", 8),
    ("Chiko", 7), ("Манафактура", 4), ("Фуку рамен", 9),
    ("Темп", 8), ("Пастерс", 8.75), ("Тихуана", 6),
    ("Молодость", 5), ("Министерство завтраков", 6.5),
    ("Эрик рыжий нн", 4.5), ("Краснодарский парень", 9),
    ("Крабс", 9), ("Хинкалия", 4.5), ("Папаша Билли", 8),
    ("Роберто", 7), ("Бакладжан", 8), ("Чихо", 7.25),
    ("Цейлон", 10), ("Митрич", 8), ("Корчма", 8),
    ("Селедка", 7.25), ("Mayar", 4), ("Тюбетейка", 3),
    ("Frank by Basta", 9), ("Red wall", 3),
    ("Чайхана номер 1 бор", 8), ("Yale", 6.25), ("Tako", 5),
    ("Pho Bo", 6.5), ("Lutee", 6.5), ("Счастливым людям", 4.75),
    ("Бенье", 7), ("Волконский кафе", 9),
    ("Барон около зала", 2), ("Самовар да утки", 7),
    ("Гусь в яблоках", 7.25), ("Краса", 5), ("Юла", 8.75),
    ("Pinci", 6.5), ("Печь", 6.5), ("Рибс", 7),
    ("Мистер бублик", 3), ("Корчма хуторок", 6),
    ("Самурай", 8.5), ("Патрон", 8), ("Негрони", 7.5),
    ("Август", 6.5), ("19", 9), ("Кусто", 6.75),
    ("Искры", 6),
]

lines = ["-- Auto-generated seed data", "BEGIN;", ""]

place_id = 0

# Helper: insert place, return its id
def add_place(name, city):
    global place_id
    place_id += 1
    lat, lng = scatter(city)
    lines.append(f"INSERT INTO places (id, name, city, lat, lng, created_by) VALUES ({place_id}, '{esc(name)}', '{esc(city)}', {lat}, {lng}, {ALINA_ID});")
    return place_id

# Helper: insert review + authors
def add_review(pid, score):
    s = max(1.0, min(10.0, score))
    lines.append(f"INSERT INTO reviews (place_id, food_rating, service_rating, vibe_rating, is_gem) VALUES ({pid}, {s}, {s}, {s}, false);")
    lines.append(f"INSERT INTO review_authors (review_id, user_id) SELECT currval('reviews_id_seq'), {ALINA_ID};")
    lines.append(f"INSERT INTO review_authors (review_id, user_id) SELECT currval('reviews_id_seq'), {LAMARK_ID};")

# Helper: add to wishlist for both users
def add_wishlist(pid):
    lines.append(f"INSERT INTO wishlists (user_id, place_id) VALUES ({ALINA_ID}, {pid});")
    lines.append(f"INSERT INTO wishlists (user_id, place_id) VALUES ({LAMARK_ID}, {pid});")

# ---- Казань: Сходить ----
lines.append("\n-- Казань: Сходить (wishlist)")
for name in kazan_wishlist:
    pid = add_place(name, "Казань")
    add_wishlist(pid)

# ---- Казань: Посещённые ----
lines.append("\n-- Казань: Посещённые")
for name, score in kazan_visited:
    pid = add_place(name, "Казань")
    add_review(pid, score)

# ---- НН: Сходить ----
lines.append("\n-- Нижний Новгород: Сходить (wishlist)")
for name in nn_wishlist:
    pid = add_place(name, "Нижний Новгород")
    add_wishlist(pid)

# ---- НН: Посещённые ----
lines.append("\n-- Нижний Новгород: Посещённые")
for name, score in nn_visited:
    pid = add_place(name, "Нижний Новгород")
    add_review(pid, score)

# Fix sequence
lines.append(f"\nSELECT setval('places_id_seq', {place_id});")
lines.append("")
lines.append("COMMIT;")

print("\n".join(lines))
