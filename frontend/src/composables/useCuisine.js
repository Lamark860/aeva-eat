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

// Готовая фраза. Возвращает '' если кухни нет или count меньше двух
// (один визит — не «любит»; в таком случае фразу прячем).
export function favoriteCuisinePhrase(profile) {
  if (!profile) return ''
  const name = profile.favorite_cuisine
  const n = profile.favorite_cuisine_count || 0
  if (!name || n < 2) return ''
  return `любит ${cuisineAccusative(name)} — ${n} ${razSuffix(n)}`
}
