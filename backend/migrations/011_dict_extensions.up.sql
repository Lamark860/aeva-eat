-- 011_dict_extensions.up.sql — round out cuisine list and add the «Понты» category

INSERT INTO cuisine_types (name) VALUES
    ('Европейская'),
    ('Авторская'),
    ('Вьетнамская')
ON CONFLICT (name) DO NOTHING;

INSERT INTO categories (name) VALUES
    ('Понты')
ON CONFLICT (name) DO NOTHING;
