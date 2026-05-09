/* AEVA EAT — main screen variants (mobile, scrapbook) */

// ───────────── shared scrapbook primitives ─────────────

function Polaroid({ photoClass = 'warm', caption, w = 120, h = 120, tape, stamp, gem, tilt, style, children }) {
  const tiltClass = tilt ? `t-${tilt}` : '';
  return (
    <div className={`polaroid ${tiltClass}`} style={{ ...style }}>
      <div className="photo" style={{ width: w, height: h }}>
        <div className={`photo ${photoClass}`} style={{ width: w, height: h }}>{children}</div>
      </div>
      {caption && <div className="caption">{caption}</div>}
      {tape}
      {gem && <div style={{ position: 'absolute', top: 6, right: 6 }}><GemMark/></div>}
      {stamp && <div style={{ position: 'absolute', top: -6, left: -8, transform: 'rotate(-8deg)' }}>{stamp}</div>}
    </div>
  );
}

function Tape({ variant = '', style }) {
  return <div className={`tape ${variant}`} style={style} />;
}

function Stamp({ children, kind = '' }) {
  return <span className={`stamp ${kind}`}>{children}</span>;
}

function GemMark() {
  // small drawn "gem" — diamond outline in terracotta ink
  return (
    <svg width="22" height="22" viewBox="0 0 22 22" style={{ filter: 'drop-shadow(0 1px 1px rgba(0,0,0,0.1))' }}>
      <path d="M11 2 L19 9 L11 20 L3 9 Z" fill="oklch(0.93 0.05 85 / 0.9)" stroke="oklch(0.55 0.14 30)" strokeWidth="1.5" strokeLinejoin="round"/>
      <path d="M11 2 L7 9 L11 20 M11 2 L15 9 L11 20 M3 9 L19 9" stroke="oklch(0.55 0.14 30)" strokeWidth="1" fill="none" opacity="0.7"/>
    </svg>
  );
}

function Ticket({ food, service, vibe, compact }) {
  return (
    <div className="ticket" style={compact ? { transform: 'scale(0.92)', transformOrigin: 'left center' } : null}>
      <div className="stub"><span className="lbl">еда</span><span className="val">{food}</span></div>
      <div className="stub"><span className="lbl">сервис</span><span className="val">{service}</span></div>
      <div className="stub"><span className="lbl">вайб</span><span className="val">{vibe}</span></div>
    </div>
  );
}

function Note({ children, lined, style, className = '' }) {
  return <div className={`note ${lined ? 'lined' : ''} ${className}`} style={style}>{children}</div>;
}

function Avatar({ name, color = '' }) {
  return <span className={`avatar ${color}`}>{name}</span>;
}

function AvatarStack({ people }) {
  return (
    <div style={{ display: 'flex' }}>
      {people.map((p, i) => (
        <div key={i} style={{ marginLeft: i ? -8 : 0, position: 'relative', zIndex: people.length - i }}>
          <Avatar name={p.initial} color={p.color}/>
        </div>
      ))}
    </div>
  );
}

// ───────────── header / tab bar ─────────────

function AppHeader({ subtitle }) {
  return (
    <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', padding: '4px 16px 8px' }}>
      <div>
        <div className="wordmark">aeva<span className="dot"/>eat</div>
        {subtitle && <div style={{ fontFamily: 'var(--hand)', fontSize: 14, color: 'var(--ink-mute)', marginTop: 2 }}>{subtitle}</div>}
      </div>
      <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
        <Stamp kind="ink">май&nbsp;'26</Stamp>
        <div style={{ width: 32, height: 32, borderRadius: '50%', background: 'oklch(0.78 0.10 30)', display: 'flex', alignItems: 'center', justifyContent: 'center', fontFamily: 'var(--serif)', fontWeight: 600, fontSize: 12, color: 'oklch(0.22 0.05 30)', boxShadow: 'inset 0 0 0 1px rgba(40,30,20,0.2)' }}>М</div>
      </div>
    </div>
  );
}

function TabBar({ active = 'board' }) {
  const tabs = [
    { id: 'board', label: 'Доска', glyph: '✦' },
    { id: 'map', label: 'Карта', glyph: '◎' },
    { id: 'find', label: 'Найти', glyph: '⌕' },
    { id: 'me', label: 'Я', glyph: '✺' },
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

// ───────────── mock data ─────────────

const FRIENDS = {
  me:    { initial: 'М', name: 'Мне', color: 'terra' },
  anya:  { initial: 'А', name: 'Аня', color: 'ochre' },
  syoma: { initial: 'С', name: 'Серёжа', color: 'moss' },
  liza:  { initial: 'Л', name: 'Лиза', color: 'plum' },
};

// ───────────── VARIANT B — Доска + Лента (default) ─────────────

function VariantB() {
  return (
    <div className="paper grain" style={{ minHeight: '100%', paddingBottom: 90 }}>
      <div style={{ paddingTop: 50 }}>
        <AppHeader subtitle="что у нас на этой неделе"/>
      </div>

      {/* PINBOARD */}
      <div className="section-head" style={{ marginTop: 4 }}>
        <h2>На доске</h2>
        <span className="sub">·  свежее</span>
      </div>

      <div style={{ position: 'relative', height: 280, margin: '8px 0 6px', overflow: 'hidden' }}>
        {/* polaroid 1 — Saviv (gem) */}
        <div style={{ position: 'absolute', left: 18, top: 14 }}>
          <Polaroid photoClass="brick" caption="Saviv · вторник" w={130} h={130} tilt="l3"
            tape={<Tape style={{ top: -10, left: 30, transform: 'rotate(-12deg)' }}/>}
            gem
          />
        </div>
        {/* polaroid 2 — Probka */}
        <div style={{ position: 'absolute', right: 14, top: 30 }}>
          <Polaroid photoClass="olive" caption="Probka · сб" w={108} h={108} tilt="r2"
            tape={<Tape variant="rose" style={{ top: -8, right: 18, transform: 'rotate(8deg)' }}/>}
          />
        </div>
        {/* note pinned */}
        <div style={{ position: 'absolute', left: 26, bottom: 10, width: 145 }} className="t-r2">
          <Note lined>
            не забыть в Тбилиси:<br/>
            <span style={{ color: 'var(--terracotta)' }}>Stamba</span> + 8000 Vintage
          </Note>
          <Tape variant="mint" style={{ top: -8, left: 50, transform: 'rotate(-4deg)' }}/>
        </div>
        {/* polaroid 3 — Хачо */}
        <div style={{ position: 'absolute', right: 24, bottom: 14 }}>
          <Polaroid photoClass="dusk" caption="Хачо и Пури" w={96} h={96} tilt="l2"
            tape={<Tape variant="blue" style={{ top: -8, left: 22, transform: 'rotate(14deg)' }}/>}
          />
        </div>
        {/* small ticket pinned */}
        <div style={{ position: 'absolute', right: 130, bottom: 130, transform: 'rotate(-6deg)' }}>
          <Ticket food="9" service="7" vibe="9" compact/>
        </div>
      </div>

      {/* FEED */}
      <div className="section-head" style={{ marginTop: 16 }}>
        <h2>Лента</h2>
        <span className="sub">·  по порядку</span>
      </div>

      <div style={{ marginTop: 12 }}>
        <FeedCard
          name="Saviv"
          city="Москва"
          cuisine="израильская"
          dateHand="6 мая, вторник"
          food="9" service="8" vibe="9"
          gem
          people={[FRIENDS.me, FRIENDS.anya]}
          photoClass="brick"
          note="хумус с грибами — лучший в городе. вернуться обязательно."
        />
        <FeedCard
          name="Probka на Цветном"
          city="Москва"
          cuisine="итальянская"
          dateHand="3 мая, суббота"
          food="8" service="7" vibe="8"
          people={[FRIENDS.me, FRIENDS.anya, FRIENDS.syoma]}
          photoClass="olive"
        />
      </div>
      <div style={{ position: 'absolute', bottom: 0, left: 0, right: 0 }}>
        <TabBar active="board"/>
      </div>
    </div>
  );
}

function FeedCard({ name, city, cuisine, dateHand, food, service, vibe, gem, people, photoClass, note }) {
  return (
    <div className="page-card">
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', gap: 10 }}>
        <div style={{ flex: 1, minWidth: 0 }}>
          <h3 className="place-name">{name}</h3>
          <div className="place-meta">
            <span>{city}</span>
            <span style={{ opacity: 0.4 }}>·</span>
            <Stamp>{cuisine}</Stamp>
            {gem && <Stamp kind="gem">жемчужина</Stamp>}
          </div>
        </div>
        <div className="date-hand">{dateHand}</div>
      </div>

      <div style={{ display: 'flex', gap: 12, marginTop: 12, alignItems: 'flex-start' }}>
        <div className="polaroid t-l1" style={{ flexShrink: 0 }}>
          <div className={`photo ${photoClass}`} style={{ width: 110, height: 110 }}/>
        </div>
        <div style={{ flex: 1, display: 'flex', flexDirection: 'column', gap: 10, paddingTop: 4 }}>
          <Ticket food={food} service={service} vibe={vibe} compact/>
          <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
            <AvatarStack people={people}/>
            <span style={{ fontFamily: 'var(--hand)', fontSize: 15, color: 'var(--ink-mute)' }}>
              {people.map(p => p.name).join(' и ')}
            </span>
          </div>
        </div>
      </div>

      {note && (
        <div style={{ marginTop: 12, fontFamily: 'var(--hand)', fontSize: 17, color: 'var(--ink)', lineHeight: 1.25, padding: '4px 4px 0', borderTop: '1px dashed rgba(40,30,20,0.2)', paddingTop: 10 }}>
          {note}
        </div>
      )}
    </div>
  );
}

// ───────────── VARIANT A — Pure pinboard ─────────────

function VariantA() {
  return (
    <div className="paper grain" style={{ minHeight: '100%', paddingBottom: 90, position: 'relative' }}>
      <div style={{ paddingTop: 50 }}>
        <AppHeader subtitle="доска · 4-9 мая"/>
      </div>

      <div style={{ position: 'relative', height: 760, margin: '4px 0' }}>
        {/* big anchor polaroid */}
        <div style={{ position: 'absolute', left: 22, top: 12 }}>
          <Polaroid photoClass="brick" caption="Saviv" w={150} h={150} tilt="l3" gem
            tape={<Tape style={{ top: -10, left: 40, transform: 'rotate(-10deg)' }}/>}
          />
        </div>
        {/* stamp */}
        <div style={{ position: 'absolute', right: 30, top: 24 }}>
          <Stamp kind="ink">москва · 8</Stamp>
        </div>
        {/* small note */}
        <div style={{ position: 'absolute', right: 18, top: 60, width: 150 }} className="t-r3">
          <Note lined>
            <span style={{ color: 'var(--terracotta)' }}>хумус с грибами</span><br/>
            лучший в городе.<br/>
            вернуться!
          </Note>
          <Tape variant="rose" style={{ top: -8, left: 40, transform: 'rotate(6deg)' }}/>
        </div>
        {/* ticket pinned floating */}
        <div style={{ position: 'absolute', left: 14, top: 230, transform: 'rotate(-3deg)' }}>
          <Ticket food="9" service="8" vibe="9"/>
        </div>
        {/* second polaroid */}
        <div style={{ position: 'absolute', right: 18, top: 240 }}>
          <Polaroid photoClass="olive" caption="Probka · сб" w={120} h={120} tilt="r2"
            tape={<Tape variant="mint" style={{ top: -8, right: 16, transform: 'rotate(10deg)' }}/>}
          />
        </div>
        {/* avatars cluster — "сходили вместе" */}
        <div style={{ position: 'absolute', left: 28, top: 380, display: 'flex', flexDirection: 'column', alignItems: 'flex-start', gap: 6 }}>
          <span style={{ fontFamily: 'var(--hand)', fontSize: 16, color: 'var(--ink-mute)' }}>с нами были →</span>
          <AvatarStack people={[FRIENDS.me, FRIENDS.anya, FRIENDS.syoma, FRIENDS.liza]}/>
        </div>
        {/* third polaroid — Хачо */}
        <div style={{ position: 'absolute', left: 18, top: 470 }}>
          <Polaroid photoClass="dusk" caption="Хачо и Пури · ср" w={138} h={138} tilt="l2"
            tape={<Tape variant="blue" style={{ top: -8, right: 30, transform: 'rotate(-8deg)' }}/>}
          />
        </div>
        {/* stamps cluster */}
        <div style={{ position: 'absolute', right: 22, top: 410, display: 'flex', flexDirection: 'column', gap: 8, alignItems: 'flex-end' }}>
          <Stamp kind="moss">грузинская</Stamp>
          <Stamp kind="plum">тёплый вечер</Stamp>
          <Stamp>4 человека</Stamp>
        </div>
        {/* fourth polaroid — Pinch */}
        <div style={{ position: 'absolute', right: 16, top: 510 }}>
          <Polaroid photoClass="indigo" caption="Pinch" w={102} h={102} tilt="r3"
            tape={<Tape style={{ top: -8, left: 22, transform: 'rotate(12deg)' }}/>}
          />
        </div>
        {/* handwritten arrow note */}
        <div style={{ position: 'absolute', left: 165, top: 600, fontFamily: 'var(--hand)', fontSize: 17, color: 'var(--terracotta)', transform: 'rotate(-6deg)' }}>
          ← сюда ещё раз<br/>
          &nbsp;&nbsp;на dim sum
        </div>
        {/* ticket bottom-left */}
        <div style={{ position: 'absolute', left: 22, top: 645, transform: 'rotate(2deg)' }}>
          <Ticket food="8" service="6" vibe="9"/>
        </div>
      </div>

      <div style={{ position: 'absolute', bottom: 0, left: 0, right: 0 }}>
        <TabBar active="board"/>
      </div>
    </div>
  );
}

// ───────────── VARIANT C — Месячный разворот + хроника ─────────────

function VariantC() {
  return (
    <div className="paper grain" style={{ minHeight: '100%', paddingBottom: 90 }}>
      <div style={{ paddingTop: 50 }}>
        <AppHeader subtitle="ваш май"/>
      </div>

      {/* Month cover */}
      <div style={{ margin: '8px 14px 0', position: 'relative' }}>
        <div style={{
          background: '#fdfcf7',
          padding: '20px 18px 18px',
          boxShadow: '0 1px 1px rgba(40,30,20,0.06), 0 6px 18px rgba(40,30,20,0.10)',
          position: 'relative',
          minHeight: 380,
        }}>
          {/* tape on cover */}
          <Tape style={{ top: -10, left: 30, transform: 'rotate(-4deg)' }}/>
          <Tape variant="rose" style={{ top: -10, right: 24, transform: 'rotate(6deg)' }}/>

          <div style={{ display: 'flex', alignItems: 'baseline', justifyContent: 'space-between' }}>
            <div>
              <div style={{ fontFamily: 'var(--serif)', fontStyle: 'italic', fontWeight: 500, fontSize: 44, lineHeight: 0.95, color: 'var(--ink)' }}>Май</div>
              <div style={{ fontFamily: 'var(--hand)', fontSize: 22, color: 'var(--ink-mute)', marginTop: 2 }}>две тысячи двадцать шесть</div>
            </div>
            <div style={{ fontFamily: 'var(--serif)', fontSize: 11, letterSpacing: '0.18em', textTransform: 'uppercase', color: 'var(--ink-mute)', textAlign: 'right' }}>
              стр.<br/>
              <span style={{ fontFamily: 'var(--hand)', fontSize: 32, letterSpacing: 0, textTransform: 'none', color: 'var(--ink)' }}>05</span>
            </div>
          </div>

          {/* monthly stats — "ticket" style */}
          <div style={{ display: 'flex', gap: 8, marginTop: 16 }}>
            <Ticket food="6" service="·" vibe="·"/>
          </div>
          <div style={{ fontFamily: 'var(--hand)', fontSize: 16, color: 'var(--ink-mute)', marginTop: 6 }}>
            <span style={{ color: 'var(--terracotta)' }}>6 мест</span> ·  3 жемчужины ·  2 города
          </div>

          {/* small collage of month's polaroids */}
          <div style={{ position: 'relative', height: 170, marginTop: 14 }}>
            <div style={{ position: 'absolute', left: 0, top: 6 }}>
              <Polaroid photoClass="brick" caption="Saviv" w={88} h={88} tilt="l3" gem/>
            </div>
            <div style={{ position: 'absolute', left: 80, top: 32 }}>
              <Polaroid photoClass="olive" caption="Probka" w={84} h={84} tilt="r2"/>
            </div>
            <div style={{ position: 'absolute', left: 160, top: 4 }}>
              <Polaroid photoClass="dusk" caption="Хачо" w={92} h={92} tilt="l2"/>
            </div>
            <div style={{ position: 'absolute', left: 12, top: 95 }}>
              <Polaroid photoClass="indigo" caption="Pinch" w={78} h={78} tilt="r3"/>
            </div>
            <div style={{ position: 'absolute', left: 92, top: 110 }}>
              <Polaroid photoClass="sage" caption="Cutfish" w={82} h={82} tilt="l3"/>
            </div>
          </div>
        </div>
      </div>

      {/* "раньше" — chronology */}
      <div className="section-head" style={{ marginTop: 22 }}>
        <h2>Раньше</h2>
        <span className="sub">· по месяцам</span>
      </div>

      <div style={{ padding: '12px 14px 0' }}>
        <ChapterStrip month="Апрель" stat="9 мест · 1 жемч. · Москва"/>
        <ChapterStrip month="Март" stat="11 мест · 2 жемч. · Москва, СПб"/>
        <ChapterStrip month="Февраль" stat="4 места · Тбилиси"/>
      </div>

      <div style={{ position: 'absolute', bottom: 0, left: 0, right: 0 }}>
        <TabBar active="board"/>
      </div>
    </div>
  );
}

function ChapterStrip({ month, stat }) {
  return (
    <div style={{
      background: '#fdfcf7',
      padding: '12px 14px',
      marginBottom: 10,
      display: 'flex', alignItems: 'baseline', justifyContent: 'space-between',
      boxShadow: '0 1px 1px rgba(40,30,20,0.06), 0 3px 8px rgba(40,30,20,0.06)',
      position: 'relative',
    }}>
      <div style={{ fontFamily: 'var(--serif)', fontStyle: 'italic', fontWeight: 500, fontSize: 22 }}>{month}</div>
      <div style={{ fontFamily: 'var(--hand)', fontSize: 15, color: 'var(--ink-mute)' }}>{stat}</div>
    </div>
  );
}

Object.assign(window, { VariantA, VariantB, VariantC });
