-- 011_dict_extensions.down.sql

DELETE FROM cuisine_types WHERE name IN ('Европейская', 'Авторская', 'Вьетнамская');
DELETE FROM categories WHERE name = 'Понты';
