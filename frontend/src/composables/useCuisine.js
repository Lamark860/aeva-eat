// DESIGN-DECISIONS Q8 — рендер «любит грузинскую — 11 раз» в Profile/PersonPage.
// Все кухни в каталоге заканчиваются на «-ская/-ная»: «Грузинская», «Японская», и т.д.
// Аккузатив получается заменой финального «ая» → «ую».

export function cuisineAccusative(name) {
  if (!name) return ''
  return name.toLowerCase().replace(/ая$/, 'ую')
}

export function razSuffix(n) {
  const mod10 = n % 10
  const mod100 = n % 100
  if (mod100 >= 11 && mod100 <= 14) return 'раз'
  if (mod10 === 1) return 'раз'
  if (mod10 >= 2 && mod10 <= 4) return 'раза'
  return 'раз'
}

// B7 — фразу «любит X» показываем только когда это правда, а не сплетня.
// Порог: ≥5 визитов кухни И ≥15% от всех визитов пользователя. У кого 144
// мест, «2 раза» — это не «любит», а «случайно попробовал».
const LOVE_MIN_COUNT = 5
const LOVE_MIN_SHARE = 0.15

export function favoriteCuisinePhrase(profile) {
  if (!profile) return ''
  const name = profile.favorite_cuisine
  const n = profile.favorite_cuisine_count || 0
  const total = profile.place_count || 0
  if (!name || n < LOVE_MIN_COUNT) return ''
  if (total > 0 && n / total < LOVE_MIN_SHARE) return ''
  return `любит ${cuisineAccusative(name)} — ${n} ${razSuffix(n)}`
}
