/* AEVA EAT — v3: drawings for 5 pain points
   1. PhotoFreeCard   — артефакт без фото (3 раскладки)
   2. CityGuide       — Город как путеводитель
   3. GemsHub         — жемчужины как сокровищница
   4. PublicShare     — share-страница с правильными пропорциями
   5. PlaceHeaderBA   — до/после очистки штампов
*/

// ─────────────────────────────────────────────────────────
// Shared mini-primitives (reuse classes from scrapbook.css)
// ─────────────────────────────────────────────────────────

function SmallStamp({ children, kind = '' }) {
  return <span className={`stamp ${kind}`} style={{ fontSize: 9 }}>{children}</span>;
}

function MiniTicket({ food, service, vibe }) {
  return (
    <div className="ticket" style={{ transform: 'scale(0.78)', transformOrigin: 'left center' }}>
      <div className="stub"><span className="lbl">еда</span><span className="val">{food}</span></div>
      <div className="stub"><span className="lbl">серв</span><span className="val">{service}</span></div>
      <div className="stub"><span className="lbl">вайб</span><span className="val">{vibe}</span></div>
    </div>
  );
}

function GemDiamond({ size = 26 }) {
  return (
    <span className="gem-badge" style={{ width: size, height: size }}>
      <svg viewBox="0 0 26 26">
        <defs>
          <linearGradient id={`gf${size}`} x1="0" y1="0" x2="1" y2="1">
            <stop offset="0%"  stopColor="oklch(0.95 0.04 85)"/>
            <stop offset="100%" stopColor="oklch(0.85 0.07 60)"/>
          </linearGradient>
        </defs>
        <path d="M13 2 L23 10 L13 24 L3 10 Z" fill={`url(#gf${size})`}
              stroke="oklch(0.55 0.14 30)" strokeWidth="1.6" strokeLinejoin="round"/>
        <path d="M13 2 L8 10 L13 24 M13 2 L18 10 L13 24 M3 10 L23 10"
              stroke="oklch(0.55 0.14 30)" strokeWidth="1" fill="none" opacity="0.6"/>
        <g className="gleam" style={{ mixBlendMode: 'screen' }}>
          <path d="M9 5 L11 5 L7 11 L5 11 Z" fill="white" opacity="0.9"/>
        </g>
      </svg>
    </span>
  );
}

function MiniAvatar({ initial, color = 'terra', size = 22 }) {
  const map = {
    terra: 'oklch(0.78 0.10 30)',
    ochre: 'oklch(0.84 0.10 85)',
    moss:  'oklch(0.78 0.07 145)',
    plum:  'oklch(0.78 0.07 350)',
  };
  return (
    <span style={{
      display: 'inline-flex', alignItems: 'center', justifyContent: 'center',
      width: size, height: size, borderRadius: '50%',
      background: map[color],
      color: 'var(--ink)',
      fontFamily: 'var(--serif)', fontWeight: 600, fontSize: size * 0.42,
      boxShadow: '0 0 0 2px #fdfcf7, 0 1px 2px rgba(40,30,20,0.18)',
      flexShrink: 0,
    }}>{initial}</span>
  );
}

function ScreenChrome({ children, sub, label }) {
  // mini header used inside iOS frame
  return (
    <div className="paper grain" style={{ minHeight: '100%', position: 'relative' }}>
      <div style={{ paddingTop: 52, paddingBottom: 90 }}>
        <div style={{ padding: '4px 18px 8px', display: 'flex', alignItems: 'baseline', justifyContent: 'space-between' }}>
          <div>
            <div className="wordmark">aeva<span className="dot"/>eat</div>
            {sub && <div style={{ fontFamily: 'var(--hand)', fontSize: 16, color: 'var(--ink-mute)', marginTop: 2 }}>{sub}</div>}
          </div>
          {label && <span className="stamp ink" style={{ transform: 'rotate(-2deg)' }}>{label}</span>}
        </div>
        {children}
      </div>
      <div style={{ position: 'absolute', bottom: 0, left: 0, right: 0 }}>
        <TabBar active="board"/>
      </div>
    </div>
  );
}

function SectionHead({ title, sub, right }) {
  return (
    <div style={{ display: 'flex', alignItems: 'baseline', gap: 10, padding: '16px 18px 8px' }}>
      <h2 style={{ fontFamily: 'var(--serif)', fontStyle: 'italic', fontWeight: 500, fontSize: 22, margin: 0 }}>{title}</h2>
      {sub && <span style={{ fontFamily: 'var(--hand)', fontSize: 15, color: 'var(--ink-mute)' }}>{sub}</span>}
      {right && <div style={{ marginLeft: 'auto' }}>{right}</div>}
    </div>
  );
}

// ─────────────────────────────────────────────────────────
// 1. Photo-free artifact — 3 layouts
// ─────────────────────────────────────────────────────────

function PhotoFreeCard() {
  // Variant Q — quote-dominant: цитата из отзыва крупно. Когда фото нет, но текст есть.
  // Variant T — ticket-dominant: билетик-крупно. Когда нет ни фото, ни большого текста.
  // Variant G — gem-only: место-жемчужина без фото, штамп доминирует.
  return (
    <div className="paper grain" style={{ minHeight: '100%', padding: '32px 18px 100px' }}>
      <div className="wordmark" style={{ fontSize: 22, marginBottom: 4 }}>aeva<span className="dot"/>eat</div>
      <div style={{ fontFamily: 'var(--hand)', fontSize: 16, color: 'var(--ink-mute)', marginBottom: 22 }}>
        безфотный артефакт · 3&nbsp;варианта
      </div>

      <div style={{ fontFamily: 'var(--serif)', fontStyle: 'italic', fontSize: 13, color: 'var(--ink-soft)', marginBottom: 16, lineHeight: 1.45 }}>
        Когда у места нет cover'a и нет первого review-фото — НЕ рисуем placeholder-клетку.
        Артефакт переключается в одну из трёх форм в зависимости от того, какие данные есть.
        Скрапбук-метафора сохраняется через цитату, билетик или штамп — не через фото-плейсхолдер.
      </div>

      {/* Q — quote-dominant */}
      <div style={{ marginTop: 10 }}>
        <Label>Q · цитата-доминанта (есть текст отзыва)</Label>
        <div className="page-card" style={{ margin: 0, position: 'relative', padding: '18px 16px 20px' }} className="t-l1">
          <div className="tape" style={{ top: -10, left: 40, transform: 'rotate(-8deg)', width: 56, height: 18 }}/>
          <div style={{ fontFamily: 'var(--hand)', fontSize: 22, color: 'var(--ink)', lineHeight: 1.15, paddingRight: 24 }}>
            «взяли каппучино —<br/>
            и потерялись<br/>
            <span style={{ color: 'var(--terracotta)' }}>на четыре часа»</span>
          </div>
          <div style={{ display: 'flex', alignItems: 'baseline', justifyContent: 'space-between', marginTop: 14, paddingTop: 10, borderTop: '1px dashed rgba(40,30,20,0.18)' }}>
            <div>
              <div style={{ fontFamily: 'var(--serif)', fontWeight: 500, fontSize: 17 }}>Black Market</div>
              <div style={{ fontFamily: 'var(--hand)', fontSize: 14, color: 'var(--ink-mute)' }}>СПб · кофейня</div>
            </div>
            <div style={{ display: 'flex', alignItems: 'center', gap: 6, fontFamily: 'var(--hand)', fontSize: 14, color: 'var(--ink-mute)' }}>
              <MiniAvatar initial="Л" color="plum"/>
              <span>Лиза</span>
            </div>
          </div>
          <div style={{ marginTop: 10 }}>
            <MiniTicket food="8" service="7" vibe="9"/>
          </div>
        </div>
      </div>

      {/* T — ticket-dominant */}
      <div style={{ marginTop: 28 }}>
        <Label>T · билетик-доминанта (нет текста, есть рейтинги)</Label>
        <div className="page-card" style={{ margin: 0, position: 'relative', padding: '16px 16px 18px' }} className="t-r1">
          <div className="tape rose" style={{ top: -10, right: 30, transform: 'rotate(10deg)', width: 50, height: 18 }}/>
          <div style={{ fontFamily: 'var(--serif)', fontWeight: 500, fontSize: 19, marginBottom: 2 }}>Пижоны</div>
          <div style={{ fontFamily: 'var(--hand)', fontSize: 14, color: 'var(--ink-mute)', marginBottom: 14 }}>
            Ижевск · европейская
          </div>
          {/* big ticket — increased scale */}
          <div className="ticket" style={{ transform: 'scale(1.05)', transformOrigin: 'left center' }}>
            <div className="stub"><span className="lbl">еда</span><span className="val" style={{ fontSize: 26 }}>4</span></div>
            <div className="stub"><span className="lbl">сервис</span><span className="val" style={{ fontSize: 26 }}>8</span></div>
            <div className="stub"><span className="lbl">вайб</span><span className="val" style={{ fontSize: 26 }}>3</span></div>
          </div>
          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginTop: 16 }}>
            <span style={{ fontFamily: 'var(--hand)', fontSize: 14, color: 'var(--ink-mute)' }}>9 мая</span>
            <div style={{ display: 'flex' }}>
              <MiniAvatar initial="М" color="terra"/>
              <div style={{ marginLeft: -8 }}><MiniAvatar initial="С" color="moss"/></div>
            </div>
          </div>
        </div>
      </div>

      {/* G — gem-only (rare, when place is gem but no photo, no text) */}
      <div style={{ marginTop: 28 }}>
        <Label>G · штамп-доминанта (жемчужина без фото)</Label>
        <div className="page-card" style={{ margin: 0, position: 'relative', padding: '20px 16px 18px' }} className="t-l2">
          <div className="tape blue" style={{ top: -10, left: 30, transform: 'rotate(-6deg)', width: 50, height: 18 }}/>
          {/* dominant gem */}
          <div style={{ display: 'flex', alignItems: 'center', gap: 14, marginBottom: 12 }}>
            <GemDiamond size={56}/>
            <div>
              <div style={{ fontFamily: 'var(--serif)', fontStyle: 'italic', fontSize: 11, letterSpacing: '0.2em', textTransform: 'uppercase', color: 'var(--terracotta)' }}>жемчужина</div>
              <div style={{ fontFamily: 'var(--serif)', fontWeight: 500, fontSize: 21, lineHeight: 1.1, marginTop: 2 }}>Соль</div>
              <div style={{ fontFamily: 'var(--hand)', fontSize: 14, color: 'var(--ink-mute)' }}>Казань · мясо</div>
            </div>
          </div>
          <div style={{ borderTop: '1px dashed rgba(40,30,20,0.18)', paddingTop: 10, display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
            <MiniTicket food="9" service="9.6" vibe="10"/>
            <div style={{ fontFamily: 'var(--hand)', fontSize: 14, color: 'var(--ink-mute)' }}>
              нашли — <span style={{ color: 'var(--terracotta)' }}>Серёжа</span>
            </div>
          </div>
        </div>
      </div>

      <div style={{ marginTop: 36, fontStyle: 'italic', fontSize: 12, color: 'var(--ink-mute)', borderTop: '1px dashed rgba(40,30,20,0.2)', paddingTop: 14, lineHeight: 1.5 }}>
        <b>Логика выбора:</b><br/>
        • есть фото → обычный полароид<br/>
        • нет фото, но есть текст ≥ 30 знаков → <b>Q</b><br/>
        • нет фото, нет текста, но есть рейтинги → <b>T</b><br/>
        • жемчужина без фото и без текста → <b>G</b><br/>
        • совсем ничего нет → не показывать на Доске
      </div>
    </div>
  );
}

function Label({ children }) {
  return (
    <div style={{ fontFamily: 'var(--serif)', fontStyle: 'italic', fontSize: 11, letterSpacing: '0.18em', textTransform: 'uppercase', color: 'var(--ink-mute)', marginBottom: 10 }}>
      {children}
    </div>
  );
}

// ─────────────────────────────────────────────────────────
// 2. City as guide
// ─────────────────────────────────────────────────────────

function CityGuide() {
  return (
    <ScreenChrome>
      <div style={{ padding: '0 18px' }}>
        <button style={{ background: 'none', border: 'none', fontFamily: 'var(--hand)', fontSize: 16, color: 'var(--ink-mute)', padding: 0, cursor: 'pointer' }}>
          ← к Найти
        </button>
        <div style={{ marginTop: 8, fontFamily: 'var(--serif)', fontStyle: 'italic', fontWeight: 500, fontSize: 44, lineHeight: 1, color: 'var(--ink)' }}>
          Казань
        </div>
        <div style={{ fontFamily: 'var(--hand)', fontSize: 18, color: 'var(--ink-mute)', marginTop: 4 }}>
          22 места · 9 жемчужин · вы вчетвером
        </div>
      </div>

      {/* hero strip — mini-map of THIS city */}
      <div style={{ margin: '18px 18px 0', position: 'relative' }}>
        <div style={{
          height: 156, borderRadius: 1,
          background: 'linear-gradient(135deg, oklch(0.86 0.04 110) 0%, oklch(0.82 0.05 130) 100%)',
          backgroundImage: `
            linear-gradient(135deg, oklch(0.84 0.04 110) 0%, oklch(0.78 0.05 130) 100%),
            radial-gradient(circle at 30% 40%, rgba(255,255,255,0.4) 0%, transparent 30%)
          `,
          boxShadow: '0 1px 1px rgba(40,30,20,0.08), 0 4px 12px rgba(40,30,20,0.10)',
          position: 'relative', overflow: 'hidden',
        }}>
          {/* fake river */}
          <svg width="100%" height="100%" style={{ position: 'absolute', inset: 0 }}>
            <path d="M0 80 Q80 60 160 90 T340 100" stroke="oklch(0.75 0.07 230)" strokeWidth="14" fill="none" opacity="0.55" strokeLinecap="round"/>
            <path d="M0 80 Q80 60 160 90 T340 100" stroke="oklch(0.85 0.05 230)" strokeWidth="6" fill="none" opacity="0.6" strokeLinecap="round"/>
          </svg>
          {/* pins */}
          {[
            { x: 90, y: 50, gem: true },
            { x: 145, y: 95 },
            { x: 200, y: 70, gem: true },
            { x: 240, y: 110 },
            { x: 280, y: 65 },
            { x: 110, y: 110 },
            { x: 175, y: 40, gem: true },
          ].map((p, i) => (
            <div key={i} style={{ position: 'absolute', left: p.x, top: p.y, transform: 'translate(-50%, -100%)' }}>
              {p.gem ? <GemDiamond size={20}/> : (
                <div style={{ width: 14, height: 14, borderRadius: '50%',
                  background: 'radial-gradient(circle at 35% 30%, oklch(0.7 0.16 25), oklch(0.42 0.18 25))',
                  boxShadow: '0 1px 2px rgba(40,15,5,0.4)' }}/>
              )}
            </div>
          ))}
          {/* tape on corners */}
          <div className="tape" style={{ top: -8, left: 36, transform: 'rotate(-6deg)', width: 56, height: 18 }}/>
          <div className="tape mint" style={{ bottom: -6, right: 30, transform: 'rotate(8deg)', width: 50, height: 18 }}/>
        </div>
        <div style={{ position: 'absolute', right: 8, bottom: -22, fontFamily: 'var(--hand)', fontSize: 14, color: 'var(--ink-mute)' }}>
          ↗ открыть карту
        </div>
      </div>

      {/* curated gems */}
      <div style={{ marginTop: 36 }}>
        <SectionHead title="Жемчужины Казани" sub="· 9"/>
        <div style={{ display: 'flex', gap: 12, overflowX: 'auto', padding: '4px 18px 18px', scrollbarWidth: 'none' }}>
          {[
            { name: 'Соль', cap: 'мясо · 9.6', tilt: 'l2', ph: 'brick' },
            { name: 'Палома', cap: 'веранда · 8.7', tilt: 'r1', ph: 'olive' },
            { name: 'Хинкали', cap: 'грузинская · 7.4', tilt: 'l3', ph: 'dusk' },
            { name: 'Юла', cap: '8.8', tilt: 'r2', ph: 'sage' },
          ].map((c, i) => (
            <div key={i} className={`polaroid t-${c.tilt}`} style={{ flexShrink: 0, position: 'relative', padding: '8px 8px 36px' }}>
              <div className={`photo ${c.ph}`} style={{ width: 130, height: 130 }}/>
              <div className="caption">{c.name}</div>
              <div style={{ position: 'absolute', top: 8, right: 8 }}><GemDiamond size={20}/></div>
              {i === 0 && <div className="tape" style={{ top: -10, left: 36, transform: 'rotate(-10deg)', width: 50, height: 16 }}/>}
              {i === 1 && <div className="tape rose" style={{ top: -8, right: 16, transform: 'rotate(8deg)', width: 50, height: 16 }}/>}
            </div>
          ))}
        </div>
      </div>

      {/* quote-of-the-city */}
      <div style={{ margin: '8px 18px 0' }}>
        <div style={{ position: 'relative', padding: '18px 18px 16px', background: '#fdfcf7',
          boxShadow: '0 1px 1px rgba(40,30,20,0.06), 0 4px 12px rgba(40,30,20,0.08)' }} className="t-l1">
          <div className="tape blue" style={{ top: -10, left: 60, transform: 'rotate(-4deg)', width: 50, height: 18 }}/>
          <div style={{ fontFamily: 'var(--serif)', fontStyle: 'italic', fontSize: 9, letterSpacing: '0.18em', textTransform: 'uppercase', color: 'var(--ink-mute)', marginBottom: 6 }}>
            цитата от круга
          </div>
          <div style={{ fontFamily: 'var(--hand)', fontSize: 22, lineHeight: 1.2, color: 'var(--ink)' }}>
            «вино — домашнее,<br/>
            хозяйка — добрая,<br/>
            <span style={{ color: 'var(--terracotta)' }}>ушли только в три»</span>
          </div>
          <div style={{ marginTop: 10, fontFamily: 'var(--hand)', fontSize: 14, color: 'var(--ink-mute)' }}>
            — Серёжа о месте&nbsp;<i style={{ fontFamily: 'var(--serif)' }}>Соль</i>
          </div>
        </div>
      </div>

      {/* all places */}
      <div style={{ marginTop: 22 }}>
        <SectionHead title="Все 22" sub="· по алфавиту"/>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 10, padding: '0 18px' }}>
          {[
            ['English', 'азиатская', 'sage'],
            ['Sleep', 'кафе', null],
            ['Терра', '·', 'cream'],
            ['Май', 'азия', null],
            ['Огни', 'бар', 'indigo'],
            ['So Sweet', 'десерты', null],
          ].map(([name, sub, ph], i) => (
            <div key={i} style={{ background: '#fdfcf7', padding: '8px 10px 10px', display: 'flex', gap: 8, alignItems: 'center',
              boxShadow: '0 1px 1px rgba(40,30,20,0.05), 0 2px 6px rgba(40,30,20,0.06)' }}>
              {ph ? (
                <div className={`photo ${ph}`} style={{ width: 36, height: 36, flexShrink: 0 }}/>
              ) : (
                <div style={{ width: 36, height: 36, flexShrink: 0, display: 'flex', alignItems: 'center', justifyContent: 'center',
                  background: 'oklch(0.94 0.04 85)', fontFamily: 'var(--serif)', fontStyle: 'italic', fontSize: 17, color: 'var(--ink-mute)' }}>
                  {name[0]}
                </div>
              )}
              <div style={{ minWidth: 0, flex: 1 }}>
                <div style={{ fontFamily: 'var(--serif)', fontWeight: 500, fontSize: 13, lineHeight: 1.1, whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>{name}</div>
                <div style={{ fontFamily: 'var(--hand)', fontSize: 12, color: 'var(--ink-mute)' }}>{sub}</div>
              </div>
            </div>
          ))}
        </div>
        <div style={{ textAlign: 'center', padding: '12px 0', fontFamily: 'var(--hand)', fontSize: 15, color: 'var(--ink-mute)' }}>
          ↓ ещё 16
        </div>
      </div>
    </ScreenChrome>
  );
}

// ─────────────────────────────────────────────────────────
// 3. Gems hub as treasury
// ─────────────────────────────────────────────────────────

function GemsHub() {
  return (
    <ScreenChrome>
      <div style={{ padding: '0 18px' }}>
        <button style={{ background: 'none', border: 'none', fontFamily: 'var(--hand)', fontSize: 16, color: 'var(--ink-mute)', padding: 0, cursor: 'pointer' }}>
          ← к Доске
        </button>
        <div style={{ display: 'flex', alignItems: 'baseline', gap: 10, marginTop: 8 }}>
          <GemDiamond size={32}/>
          <div style={{ fontFamily: 'var(--serif)', fontStyle: 'italic', fontWeight: 500, fontSize: 36, lineHeight: 1, color: 'var(--ink)' }}>
            Жемчужины
          </div>
        </div>
        <div style={{ fontFamily: 'var(--hand)', fontSize: 18, color: 'var(--ink-mute)', marginTop: 4 }}>
          19 в трёх городах
        </div>
      </div>

      {/* first gem — hero */}
      <div style={{ margin: '24px 18px 0' }}>
        <div style={{ fontFamily: 'var(--serif)', fontStyle: 'italic', fontSize: 11, letterSpacing: '0.2em', textTransform: 'uppercase', color: 'var(--ink-mute)', marginBottom: 10 }}>
          самая первая
        </div>
        <div style={{ display: 'flex', gap: 14, alignItems: 'flex-start' }}>
          <div className="polaroid t-l3" style={{ position: 'relative', padding: '8px 8px 36px' }}>
            <div className="photo brick" style={{ width: 142, height: 142 }}/>
            <div className="caption">Соль · Казань</div>
            <div className="tape" style={{ top: -10, left: 36, transform: 'rotate(-10deg)', width: 56, height: 18 }}/>
            <div style={{ position: 'absolute', top: 8, right: 8 }}><GemDiamond size={26}/></div>
          </div>
          <div style={{ flex: 1, paddingTop: 6 }}>
            <div style={{ fontFamily: 'var(--serif)', fontWeight: 500, fontSize: 20, lineHeight: 1.1 }}>Соль</div>
            <div style={{ fontFamily: 'var(--hand)', fontSize: 15, color: 'var(--ink-mute)', marginTop: 2 }}>
              Казань · мясо
            </div>
            <div style={{ fontFamily: 'var(--hand)', fontSize: 14, color: 'var(--terracotta)', marginTop: 14, lineHeight: 1.3 }}>
              Алина нашла<br/>
              5 апреля 2026 ·<br/>
              подхватили все&nbsp;четверо
            </div>
          </div>
        </div>
      </div>

      {/* by city — as stamps */}
      <div style={{ marginTop: 28 }}>
        <SectionHead title="По городам"/>
        <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8, padding: '0 18px' }}>
          {[
            ['Нижний', 11],
            ['Казань', 6],
            ['Ижевск', 2],
          ].map(([c, n], i) => (
            <span key={i} className="stamp ink" style={{ fontSize: 11, transform: `rotate(${i % 2 === 0 ? -2 : 2}deg)`, display: 'inline-flex', alignItems: 'center', gap: 6 }}>
              {c} <span style={{ fontFamily: 'var(--hand)', fontSize: 14, color: 'var(--terracotta)', textTransform: 'none', letterSpacing: 0 }}>{n}</span>
            </span>
          ))}
        </div>
      </div>

      {/* who found — avatars */}
      <div style={{ marginTop: 18 }}>
        <SectionHead title="Кто отмечал"/>
        <div style={{ display: 'flex', gap: 18, padding: '0 18px', overflowX: 'auto', scrollbarWidth: 'none' }}>
          {[
            { i: 'А', c: 'ochre', name: 'alina', count: 7 },
            { i: 'Л', c: 'terra', name: 'lamark', count: 8 },
            { i: 'С', c: 'moss', name: 'Серёжа', count: 3 },
            { i: 'Ч', c: 'plum', name: 'charlie', count: 1 },
          ].map((f, i) => (
            <div key={i} style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 6, flexShrink: 0 }}>
              <MiniAvatar initial={f.i} color={f.c} size={50}/>
              <div style={{ fontFamily: 'var(--serif)', fontSize: 13 }}>{f.name}</div>
              <div style={{ fontFamily: 'var(--hand)', fontSize: 14, color: 'var(--terracotta)', marginTop: -4 }}>
                {f.count} ♦
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* all gems — proper spreads */}
      <div style={{ marginTop: 28 }}>
        <SectionHead title="Все жемчужины"/>
        <div style={{ padding: '0 18px' }}>
          {[
            { name: 'Палома', city: 'Казань', ph: 'olive', quote: '«взяли всё, что было в баре —\nи ничего не пожалели»', author: 'lamark', rating: '9.2 / 8.7 / 9.5' },
            { name: 'Хинкали от Лали', city: 'Казань', ph: 'dusk', quote: '«пять минут от вокзала.\nходить туда нужно голодным»', author: 'alina', rating: '5 / 2.9 / 5' },
          ].map((g, i) => (
            <div key={i} className={`page-card t-${i % 2 === 0 ? 'l1' : 'r1'}`} style={{ margin: '0 0 14px', padding: '14px 14px 16px', position: 'relative' }}>
              <div className={`tape ${i % 2 === 0 ? '' : 'mint'}`} style={{ top: -10, left: i % 2 === 0 ? 30 : 'auto', right: i % 2 === 0 ? 'auto' : 30, transform: `rotate(${i % 2 === 0 ? -8 : 8}deg)`, width: 50, height: 18 }}/>
              <div style={{ display: 'flex', gap: 12, alignItems: 'flex-start' }}>
                <div className="polaroid" style={{ flexShrink: 0, padding: '6px 6px 22px', position: 'relative' }}>
                  <div className={`photo ${g.ph}`} style={{ width: 96, height: 96 }}/>
                  <div className="caption" style={{ fontSize: 14 }}>{g.name.split(' ')[0]}</div>
                  <div style={{ position: 'absolute', top: 6, right: 6 }}><GemDiamond size={18}/></div>
                </div>
                <div style={{ flex: 1, minWidth: 0 }}>
                  <div style={{ fontFamily: 'var(--serif)', fontWeight: 500, fontSize: 17, lineHeight: 1.1 }}>{g.name}</div>
                  <div style={{ fontFamily: 'var(--hand)', fontSize: 14, color: 'var(--ink-mute)', marginTop: 2 }}>{g.city}</div>
                  <div style={{ fontFamily: 'var(--hand)', fontSize: 16, color: 'var(--ink)', marginTop: 8, lineHeight: 1.2, whiteSpace: 'pre-line' }}>
                    {g.quote}
                  </div>
                  <div style={{ display: 'flex', alignItems: 'center', gap: 6, marginTop: 8, fontFamily: 'var(--hand)', fontSize: 13, color: 'var(--ink-mute)' }}>
                    — {g.author} <span style={{ color: 'var(--terracotta)' }}>· {g.rating}</span>
                  </div>
                </div>
              </div>
            </div>
          ))}
          <div style={{ textAlign: 'center', padding: '6px 0 12px', fontFamily: 'var(--hand)', fontSize: 15, color: 'var(--ink-mute)' }}>
            ↓ ещё 17
          </div>
        </div>
      </div>
    </ScreenChrome>
  );
}

// ─────────────────────────────────────────────────────────
// 4. Public share /p/:id — fixed proportions
// ─────────────────────────────────────────────────────────

function PublicShare() {
  return (
    <div className="paper grain" style={{ minHeight: '100%', position: 'relative', display: 'flex', flexDirection: 'column' }}>
      {/* cover — 45%, not 70% */}
      <div style={{
        height: 320, width: '100%', overflow: 'hidden', position: 'relative',
        background: 'oklch(0.45 0.06 30)',
        backgroundImage: `
          linear-gradient(180deg, oklch(0.5 0.07 25) 0%, oklch(0.36 0.05 30) 100%),
          radial-gradient(circle at 30% 40%, rgba(255,255,255,0.18) 0%, transparent 30%)
        `,
      }}>
        {/* fake plate — abstract shapes */}
        <div style={{ position: 'absolute', inset: 0, opacity: 0.7 }}>
          <svg width="100%" height="100%" viewBox="0 0 360 320">
            <circle cx="200" cy="180" r="100" fill="oklch(0.6 0.03 60)" opacity="0.6"/>
            <circle cx="200" cy="180" r="70" fill="oklch(0.32 0.05 30)" opacity="0.5"/>
            <ellipse cx="120" cy="140" rx="60" ry="40" fill="oklch(0.7 0.08 60)" opacity="0.45"/>
          </svg>
        </div>
        {/* corner ribbon */}
        <div style={{
          position: 'absolute', top: 12, right: 12,
          fontFamily: 'var(--serif)', fontStyle: 'italic', fontSize: 10, letterSpacing: '0.2em', textTransform: 'uppercase',
          color: 'rgba(255,255,255,0.75)',
          padding: '4px 8px',
          border: '1px solid rgba(255,255,255,0.4)',
          borderRadius: 2,
        }}>aeva·eat</div>
      </div>

      {/* paper card — overlapping bottom of cover */}
      <div style={{ margin: '-40px 22px 24px', position: 'relative', zIndex: 1 }}>
        <div className="t-l1" style={{
          background: '#fdfcf7',
          padding: '24px 22px 22px',
          boxShadow: '0 2px 4px rgba(40,30,20,0.1), 0 8px 22px rgba(40,30,20,0.16)',
          position: 'relative',
        }}>
          <div className="tape" style={{ top: -10, left: 36, transform: 'rotate(-8deg)', width: 56, height: 18 }}/>
          <div className="tape rose" style={{ top: -10, right: 36, transform: 'rotate(8deg)', width: 50, height: 18 }}/>

          <div style={{ fontFamily: 'var(--serif)', fontStyle: 'italic', fontSize: 9, letterSpacing: '0.22em', textTransform: 'uppercase', color: 'var(--ink-mute)', textAlign: 'center' }}>
            из дневника круга
          </div>
          <div style={{ fontFamily: 'var(--serif)', fontWeight: 500, fontSize: 30, lineHeight: 1.05, textAlign: 'center', marginTop: 6 }}>
            Мясной гуру
          </div>
          <div style={{ fontFamily: 'var(--hand)', fontSize: 18, color: 'var(--ink-mute)', textAlign: 'center', marginTop: 4 }}>
            Ижевск · грузинская
          </div>

          <div style={{ display: 'flex', justifyContent: 'center', marginTop: 14 }}>
            <span className="stamp gem" style={{ transform: 'rotate(-3deg)' }}>жемчужина</span>
          </div>

          {/* CTA — paper button with tilt, not loud oval */}
          <div style={{ marginTop: 22, display: 'flex', justifyContent: 'center' }}>
            <button className="t-r1" style={{
              background: '#fdfcf7',
              color: 'var(--ink)',
              fontFamily: 'var(--serif)',
              fontStyle: 'italic',
              fontSize: 17,
              padding: '12px 24px',
              border: '1.5px solid var(--ink)',
              cursor: 'pointer',
              boxShadow: '0 1px 1px rgba(40,30,20,0.06), 0 3px 8px rgba(40,30,20,0.1)',
              position: 'relative',
            }}>
              <span style={{ color: 'var(--terracotta)', marginRight: 8 }}>→</span>
              войти, чтобы увидеть впечатления
            </button>
          </div>

          <div style={{ marginTop: 18, textAlign: 'center', fontFamily: 'var(--hand)', fontSize: 14, color: 'var(--ink-mute)' }}>
            камерный дневник еды
          </div>
        </div>
      </div>

      {/* footer with friend ring hint (без раскрытия — privacy) */}
      <div style={{ flex: 1, padding: '0 22px 36px', display: 'flex', flexDirection: 'column', justifyContent: 'flex-end', alignItems: 'center' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 8, opacity: 0.7 }}>
          <MiniAvatar initial="•" color="terra" size={16}/>
          <MiniAvatar initial="•" color="ochre" size={16}/>
          <MiniAvatar initial="•" color="moss" size={16}/>
          <span style={{ fontFamily: 'var(--hand)', fontSize: 13, color: 'var(--ink-mute)' }}>
            у круга четыре впечатления
          </span>
        </div>
      </div>
    </div>
  );
}

// ─────────────────────────────────────────────────────────
// 5. Place header — before / after
// ─────────────────────────────────────────────────────────

function PlaceHeaderBA() {
  return (
    <div className="paper grain" style={{ minHeight: '100%', padding: '32px 18px 100px' }}>
      <div className="wordmark" style={{ fontSize: 22, marginBottom: 4 }}>aeva<span className="dot"/>eat</div>
      <div style={{ fontFamily: 'var(--hand)', fontSize: 16, color: 'var(--ink-mute)', marginBottom: 22 }}>
        place header · до/после
      </div>

      {/* BEFORE */}
      <Label>— до —</Label>
      <div style={{ background: '#fdfcf7', padding: '20px 16px 18px', boxShadow: '0 1px 1px rgba(40,30,20,0.06), 0 4px 12px rgba(40,30,20,0.08)', marginBottom: 6 }}>
        <div style={{ fontFamily: 'var(--serif)', fontWeight: 500, fontSize: 26, textAlign: 'center' }}>
          Мясной гуру
        </div>
        <div style={{ display: 'flex', flexWrap: 'wrap', gap: 6, justifyContent: 'center', marginTop: 12 }}>
          <SmallStamp kind="ink">ижевск</SmallStamp>
          <SmallStamp kind="ink">грузинская</SmallStamp>
          <SmallStamp kind="moss">банкет</SmallStamp>
          <SmallStamp kind="moss">обед</SmallStamp>
          <SmallStamp kind="moss">стритфуд</SmallStamp>
          <SmallStamp kind="moss">ужин</SmallStamp>
          <SmallStamp kind="gem">жемчужина</SmallStamp>
        </div>
        <div style={{ fontFamily: 'var(--hand)', fontSize: 13, color: 'var(--ink-mute)', textAlign: 'center', marginTop: 12 }}>
          отметил alina · 10 мая (+ lamark)
        </div>
      </div>
      <div style={{ fontStyle: 'italic', fontSize: 11, color: 'oklch(0.55 0.14 30)', margin: '8px 0 22px', lineHeight: 1.4 }}>
        7 штампов, повторение цвета moss, гендер-конфликт «отметил alina»
      </div>

      {/* AFTER */}
      <Label>— после —</Label>
      <div style={{ background: '#fdfcf7', padding: '24px 18px 20px', boxShadow: '0 1px 1px rgba(40,30,20,0.06), 0 4px 12px rgba(40,30,20,0.08)', position: 'relative' }} className="t-l1">
        <div className="tape" style={{ top: -10, left: 60, transform: 'rotate(-6deg)', width: 56, height: 18 }}/>
        <div style={{ fontFamily: 'var(--serif)', fontWeight: 500, fontSize: 30, textAlign: 'center', lineHeight: 1.05 }}>
          Мясной гуру
        </div>
        <div style={{ fontFamily: 'var(--hand)', fontSize: 16, color: 'var(--ink-mute)', textAlign: 'center', marginTop: 4 }}>
          ул. Карла Либкнехта, 11
        </div>
        {/* primary stamps — max 2 */}
        <div style={{ display: 'flex', gap: 8, justifyContent: 'center', marginTop: 14 }}>
          <SmallStamp kind="ink">ижевск</SmallStamp>
          <SmallStamp>грузинская</SmallStamp>
        </div>
        {/* gem separately — visually distinct */}
        <div style={{ display: 'flex', justifyContent: 'center', marginTop: 10 }}>
          <span className="stamp gem" style={{ transform: 'rotate(-2deg)', display: 'inline-flex', alignItems: 'center', gap: 6 }}>
            <span style={{ display: 'inline-block', verticalAlign: 'middle' }}><GemDiamond size={14}/></span>
            жемчужина
          </span>
        </div>
        {/* categories — small caveat, not stamps */}
        <div style={{ fontFamily: 'var(--hand)', fontSize: 14, color: 'var(--ink-mute)', textAlign: 'center', marginTop: 14, lineHeight: 1.3 }}>
          банкет · обед · стритфуд · ужин
        </div>
        {/* gender-less attribution */}
        <div style={{ fontFamily: 'var(--hand)', fontSize: 14, color: 'var(--ink-mute)', textAlign: 'center', marginTop: 14, paddingTop: 12, borderTop: '1px dashed rgba(40,30,20,0.18)' }}>
          жемчужина · <span style={{ color: 'var(--terracotta)' }}>Алина</span> · 10 мая
          <div style={{ fontSize: 12, color: 'var(--ink-mute)', marginTop: 2 }}>+ Lamark</div>
        </div>
      </div>
      <div style={{ fontStyle: 'italic', fontSize: 11, color: 'var(--moss)', margin: '8px 0 0', lineHeight: 1.4 }}>
        2 главных штампа · жемчужина выделена · категории как мелкая подпись · без гендер-глагола
      </div>
    </div>
  );
}

// ─────────────────────────────────────────────────────────
// Bottom tab bar (reuse)
// ─────────────────────────────────────────────────────────

function TabBar({ active = 'board' }) {
  const tabs = [
    { id: 'board', label: 'доска', glyph: '✦' },
    { id: 'find',  label: 'найти', glyph: '⌕' },
    { id: 'map',   label: 'карта', glyph: '◎' },
    { id: 'me',    label: 'я',     glyph: '✺' },
  ];
  return (
    <div className="tabbar">
      {tabs.map(t => (
        <div key={t.id} className={`tab ${t.id === active ? 'active' : ''}`}>
          <div className="glyph">{t.glyph}</div>
          <div>{t.label}</div>
        </div>
      ))}
    </div>
  );
}

Object.assign(window, {
  PhotoFreeCard, CityGuide, GemsHub, PublicShare, PlaceHeaderBA,
});
