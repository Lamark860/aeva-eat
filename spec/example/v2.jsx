/* AEVA EAT — v2: единый скрапбук-поток, доска + хроника */

function GemBadge({ size = 26 }) {
  return (
    <span className="gem-badge" style={{ width: size, height: size }}>
      <svg viewBox="0 0 26 26">
        <defs>
          <linearGradient id="gem-fill" x1="0" y1="0" x2="1" y2="1">
            <stop offset="0%"  stopColor="oklch(0.95 0.04 85)"/>
            <stop offset="100%" stopColor="oklch(0.85 0.07 60)"/>
          </linearGradient>
        </defs>
        <path d="M13 2 L23 10 L13 24 L3 10 Z"
              fill="url(#gem-fill)"
              stroke="oklch(0.55 0.14 30)" strokeWidth="1.6" strokeLinejoin="round"/>
        <path d="M13 2 L8 10 L13 24 M13 2 L18 10 L13 24 M3 10 L23 10"
              stroke="oklch(0.55 0.14 30)" strokeWidth="1" fill="none" opacity="0.6"/>
        {/* gleam — animated highlight sweeping across facet */}
        <g className="gleam" style={{ mixBlendMode: 'screen' }}>
          <path d="M9 5 L11 5 L7 11 L5 11 Z" fill="white" opacity="0.9"/>
          <circle cx="11" cy="6" r="1.1" fill="white" opacity="0.95"/>
        </g>
      </svg>
    </span>
  );
}

function AuthorTag({ name, color, style }) {
  const map = {
    terra: { bg: 'oklch(0.78 0.10 30)',  ink: 'oklch(0.22 0.05 30)' },
    ochre: { bg: 'oklch(0.84 0.10 85)',  ink: 'oklch(0.25 0.05 85)' },
    moss:  { bg: 'oklch(0.78 0.07 145)', ink: 'oklch(0.22 0.05 145)' },
    plum:  { bg: 'oklch(0.78 0.07 350)', ink: 'oklch(0.25 0.05 350)' },
  };
  const c = map[color] || map.terra;
  return <span className="author-tag" style={{ background: c.bg, color: c.ink, ...style }}>{name}</span>;
}

function VideoKruzhok({ duration = '0:14', author }) {
  return (
    <div className="kruzhok">
      <div className="play"/>
      <div className="duration">{duration}</div>
      {author && <AuthorTag name={author.initial} color={author.color} style={{ position: 'absolute', top: -2, right: -2 }}/>}
    </div>
  );
}

function AddButton() {
  // soft "+" pill — no jargon, self-evident
  return (
    <button className="pin-btn" style={{ paddingLeft: 10 }}>
      <span style={{ fontFamily: 'var(--serif)', fontStyle: 'normal', fontWeight: 500, fontSize: 18, lineHeight: 1, color: 'var(--terracotta)', marginTop: -1 }}>+</span>
      <span>добавить</span>
    </button>
  );
}

function TogetherTag({ people }) {
  return (
    <span style={{ display: 'inline-flex', alignItems: 'center', gap: 6, fontFamily: 'var(--hand)', fontSize: 15, color: 'var(--ink-mute)' }}>
      <span>вместе:</span>
      <AvatarStack people={people}/>
    </span>
  );
}

const F = {
  me:    { initial: 'М', name: 'я',     color: 'terra' },
  anya:  { initial: 'А', name: 'Аня',   color: 'ochre' },
  syoma: { initial: 'С', name: 'Серёжа', color: 'moss' },
  liza:  { initial: 'Л', name: 'Лиза',  color: 'plum' },
};

// ───────────── V2 — current ─────────────

function V2() {
  return (
    <div className="paper grain" style={{ minHeight: '100%', paddingBottom: 90, position: 'relative' }}>
      <div style={{ paddingTop: 50 }}>
        <AppHeader subtitle="наша доска"/>
      </div>

      {/* Current strip — date only, dense-but-airy, "развернуть" if more */}
      <div style={{ display: 'flex', alignItems: 'center', gap: 10, padding: '14px 18px 0' }}>
        <div style={{ fontFamily: 'var(--serif)', fontStyle: 'italic', fontWeight: 500, fontSize: 18, color: 'var(--ink)' }}>4–9 мая</div>
        <div style={{ flex: 1, height: 1, background: 'rgba(40,30,20,0.18)', marginTop: 4 }}/>
        <AddButton/>
      </div>

      {/* current week — 4 items, breathing room */}
      <div style={{ position: 'relative', height: 380, margin: '6px 0' }}>
        {/* Saviv — gem polaroid (Аня) */}
        <div style={{ position: 'absolute', left: 22, top: 12 }}>
          <div className="polaroid t-l3" style={{ position: 'relative' }}>
            <div className="photo brick" style={{ width: 130, height: 130 }}/>
            <div className="caption">Saviv · вторник</div>
            <Tape style={{ top: -10, left: 32, transform: 'rotate(-12deg)' }}/>
            <div style={{ position: 'absolute', top: 8, right: 8 }}><GemBadge size={28}/></div>
            <AuthorTag name="А" color="ochre" style={{ bottom: -6, left: -6 }}/>
          </div>
        </div>

        {/* note — Аня wrote */}
        <div style={{ position: 'absolute', right: 18, top: 24, width: 138 }} className="t-r2">
          <Note lined>
            <span style={{ color: 'var(--terracotta)' }}>хумус</span> с грибами —<br/>
            лучший в&nbsp;городе
          </Note>
          <Tape variant="rose" style={{ top: -8, left: 36, transform: 'rotate(6deg)' }}/>
          <AuthorTag name="А" color="ochre" style={{ bottom: -6, right: -6 }}/>
        </div>

        {/* Probka — Серёжа, lower-left */}
        <div style={{ position: 'absolute', left: 28, bottom: 16 }}>
          <div className="polaroid t-r2" style={{ position: 'relative' }}>
            <div className="photo olive" style={{ width: 108, height: 108 }}/>
            <div className="caption">Probka · сб</div>
            <Tape variant="mint" style={{ top: -8, right: 16, transform: 'rotate(8deg)' }}/>
            <AuthorTag name="С" color="moss" style={{ bottom: -6, right: -6 }}/>
          </div>
        </div>

        {/* video kruzhok bottom-right */}
        <div style={{ position: 'absolute', right: 30, bottom: 30 }} className="t-l2">
          <VideoKruzhok duration="0:18" author={F.syoma}/>
        </div>

        {/* "вместе" tag floating between */}
        <div style={{ position: 'absolute', left: 165, top: 195, transform: 'rotate(-3deg)' }}>
          <TogetherTag people={[F.me, F.anya, F.syoma]}/>
        </div>
      </div>

      <div style={{ textAlign: 'center', padding: '0 0 4px' }}>
        <button style={{
          background: 'transparent', border: 'none', cursor: 'pointer',
          fontFamily: 'var(--hand)', fontSize: 16, color: 'var(--ink-mute)',
          padding: '6px 14px',
        }}>↓ ещё 3 на этой неделе</button>
      </div>

      {/* Older — collapsed strips by default, expandable */}
      <CollapsedStrip
        dates="27 апр – 3 мая"
        summary={[
          { ph: 'indigo', cap: 'Pinch' },
          { ph: 'sage',   cap: 'Cutfish' },
          { ph: 'cream',  cap: 'Black Market' },
        ]}
        count={3}
        gemCount={0}
      />

      <CollapsedStrip
        dates="апрель"
        summary={[
          { ph: 'peach',  cap: 'Shavi Lomi', gem: true },
          { ph: 'dusk',   cap: 'Хачо' },
          { ph: 'warm',   cap: 'Bar 33', gem: true },
        ]}
        count={9}
        gemCount={2}
      />

      <CollapsedStrip
        dates="март"
        summary={[
          { ph: 'olive',  cap: 'People\u2019s' },
          { ph: 'slate',  cap: 'Cutfish' },
        ]}
        count={11}
        gemCount={2}
      />

      <div style={{ textAlign: 'center', padding: '12px 0 6px', fontFamily: 'var(--hand)', fontSize: 16, color: 'var(--ink-mute)' }}>
        ↓  раньше — в архиве
      </div>

      <div style={{ position: 'absolute', bottom: 0, left: 0, right: 0 }}>
        <TabBar active="board"/>
      </div>
    </div>
  );
}

function CollapsedStrip({ dates, summary, count, gemCount }) {
  // Compact strip: header line + tiny stack of polaroids + counts. Tap to expand.
  return (
    <div style={{ padding: '8px 14px 0' }}>
      <div style={{ display: 'flex', alignItems: 'center', gap: 10, padding: '0 4px' }}>
        <div style={{ fontFamily: 'var(--serif)', fontStyle: 'italic', fontWeight: 500, fontSize: 16, color: 'var(--ink-soft)' }}>{dates}</div>
        <div style={{ flex: 1, height: 1, background: 'rgba(40,30,20,0.14)', marginTop: 3 }}/>
        <div style={{ fontFamily: 'var(--hand)', fontSize: 14, color: 'var(--ink-mute)' }}>
          {count} {count > 4 ? 'мест' : count > 1 ? 'места' : 'место'}
          {gemCount > 0 && <> · {gemCount} <span style={{ color: 'var(--terracotta)' }}>♦</span></>}
        </div>
      </div>
      <div style={{ display: 'flex', alignItems: 'center', padding: '6px 4px 0', gap: 8 }}>
        <div style={{ position: 'relative', height: 80, flex: 1 }}>
          {summary.map((s, i) => (
            <div key={i} style={{
              position: 'absolute',
              left: i * 46,
              top: i % 2 === 0 ? 4 : 12,
              transform: `rotate(${(i % 2 === 0 ? -2 : 2) - i * 0.3}deg)`,
            }}>
              <div className="polaroid" style={{ position: 'relative', padding: '5px 5px 18px' }}>
                <div className={`photo ${s.ph}`} style={{ width: 60, height: 60 }}/>
                <div className="caption" style={{ fontSize: 11, bottom: 3 }}>{s.cap}</div>
                {s.gem && <div style={{ position: 'absolute', top: 4, right: 4 }}><GemBadge size={16}/></div>}
              </div>
            </div>
          ))}
        </div>
        <button style={{
          background: 'transparent', border: 'none', cursor: 'pointer',
          fontFamily: 'var(--hand)', fontSize: 16, color: 'var(--ink-mute)',
          padding: '4px 8px',
          flexShrink: 0,
        }}>раскрыть&nbsp;↓</button>
      </div>
    </div>
  );
}

Object.assign(window, { V2 });
