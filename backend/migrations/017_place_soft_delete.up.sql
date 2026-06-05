-- Мягкое удаление мест: место скрывается из листингов/карты/поиска/рандома, но
-- отзывы/фото/видео круга сохраняются (ценность скрапбука нельзя терять одним
-- тапом). superuser может вернуть (restore) или вычистить навсегда.
ALTER TABLE places ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

-- Идентичность места теперь учитывает только АКТИВНЫЕ места: soft-deleted не
-- блокирует повторное добавление того же заведения. Партиал — строгое
-- подмножество прежнего ключа, на текущих данных (deleted_at IS NULL у всех)
-- не конфликтует.
DROP INDEX IF EXISTS idx_places_identity;
CREATE UNIQUE INDEX IF NOT EXISTS idx_places_identity
    ON places (LOWER(name), LOWER(COALESCE(address, '')), LOWER(COALESCE(city, '')))
    WHERE deleted_at IS NULL;

-- Частичный индекс для быстрых выборок активных мест.
CREATE INDEX IF NOT EXISTS idx_places_active ON places (id) WHERE deleted_at IS NULL;
