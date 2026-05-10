import{_ as B,x as N,w as v,G as O,o as Z,c as V,n as E,p as F}from"./index-Dldv0PHR.js";const H={__name:"MapView",props:{places:{type:Array,default:()=>[]},center:{type:Array,default:()=>[55.7558,37.6173]},zoom:{type:Number,default:12},height:{type:String,default:"500px"},singleMarker:{type:Boolean,default:!1}},emits:["marker-click"],setup(u,{emit:$}){const n=u,w=$,x=F(null);let i=null,g=[];const o={terra:"#c66b3e",terraDark:"#6e2210",ochre:"#caa44a",ochreLight:"#e8c97a",ochreDark:"#7c5f1c",moss:"#7c9b6a",mossDark:"#3f5c30",ink:"#2d231b",highlight:"rgba(255,255,255,0.45)"};function L(e){if(e==null)return{light:o.ochre,dark:o.ochreDark,shade:o.ochreDark};const t=Number(e);return t>=8?{light:o.moss,dark:o.mossDark,shade:o.mossDark}:t>=5?{light:o.ochre,dark:o.ochreDark,shade:o.ochreDark}:{light:o.terra,dark:o.terraDark,shade:o.terraDark}}function y(e){if(e==null)return"";const t=Number(e);return Number.isInteger(t)?String(t):t>=9.95?"10":t.toFixed(1)}function _(e,t){const a=L(t),r=y(t),s=r.length>2?9.5:11.5;return`<svg width="40" height="52" viewBox="0 0 40 52" xmlns="http://www.w3.org/2000/svg">
    <defs>
      <radialGradient id="head-${e}" cx="35%" cy="28%">
        <stop offset="0%" stop-color="${a.light}"/>
        <stop offset="100%" stop-color="${a.dark}"/>
      </radialGradient>
      <filter id="shadow-${e}" x="-30%" y="-10%" width="160%" height="160%">
        <feDropShadow dx="0" dy="3" stdDeviation="2.5" flood-opacity="0.4"/>
      </filter>
    </defs>
    <!-- thin needle from head down to point -->
    <line x1="20" y1="20" x2="20" y2="50" stroke="${o.ink}" stroke-width="1.6" stroke-linecap="round" opacity="0.85"/>
    <circle cx="20" cy="50" r="2" fill="${o.ink}" opacity="0.8"/>
    <!-- pushpin head -->
    <circle cx="20" cy="18" r="14" fill="url(#head-${e})" filter="url(#shadow-${e})"/>
    <!-- bottom shading -->
    <ellipse cx="20" cy="22" rx="9" ry="4" fill="${a.shade}" opacity="0.25"/>
    <!-- highlight -->
    <circle cx="16" cy="13" r="3.5" fill="${o.highlight}"/>
    ${r?`<text x="20" y="${r.length>2?21.5:22}" text-anchor="middle" fill="#fff" font-family="Lora, Georgia, serif" font-weight="600" font-size="${s}" style="text-shadow:0 1px 1px rgba(0,0,0,0.35)" paint-order="stroke" stroke="rgba(40,30,20,0.35)" stroke-width="0.4">${r}</text>`:""}
  </svg>`}function z(e,t){const a=y(t),r=a.length>2?8.5:10;return`<svg width="44" height="54" viewBox="0 0 44 54" xmlns="http://www.w3.org/2000/svg">
    <defs>
      <linearGradient id="gem-${e}" x1="0" y1="0" x2="1" y2="1">
        <stop offset="0%" stop-color="${o.ochreLight}"/>
        <stop offset="100%" stop-color="${o.ochre}"/>
      </linearGradient>
      <filter id="g-shadow-${e}" x="-40%" y="-20%" width="180%" height="160%">
        <feDropShadow dx="0" dy="2" stdDeviation="2.2" flood-opacity="0.32"/>
      </filter>
      <radialGradient id="halo-${e}" cx="50%" cy="40%">
        <stop offset="0%" stop-color="${o.ochreLight}" stop-opacity="0.6"/>
        <stop offset="60%" stop-color="${o.ochreLight}" stop-opacity="0.0"/>
      </radialGradient>
    </defs>
    <!-- soft halo -->
    <circle cx="22" cy="18" r="22" fill="url(#halo-${e})"/>
    <!-- needle -->
    <line x1="22" y1="32" x2="22" y2="52" stroke="${o.ink}" stroke-width="1.6" stroke-linecap="round" opacity="0.85"/>
    <circle cx="22" cy="52" r="2" fill="${o.ink}" opacity="0.8"/>
    <!-- gem diamond -->
    <path d="M22 4 L36 18 L22 32 L8 18 Z"
          fill="url(#gem-${e})"
          stroke="${o.terraDark}"
          stroke-width="1.6"
          stroke-linejoin="round"
          filter="url(#g-shadow-${e})"/>
    <!-- facets -->
    <path d="M22 4 L16 18 L22 32 M22 4 L28 18 L22 32 M8 18 L36 18"
          stroke="${o.terraDark}" stroke-width="0.9" fill="none" opacity="0.55"/>
    <!-- highlight -->
    <path d="M14 9 L17 9 L11 19 L8 19 Z" fill="#fff" opacity="0.6"/>
    ${a?`<text x="22" y="22" text-anchor="middle" fill="${o.terraDark}" font-family="Lora, Georgia, serif" font-weight="600" font-size="${r}" paint-order="stroke" stroke="#fff" stroke-width="0.7" stroke-linejoin="round">${a}</text>`:""}
  </svg>`}function D(e,t){const a=e.image_url?`<div style="background:var(--sb-paper-card);padding:6px 6px 22px;margin:0 0 8px;box-shadow:0 1px 1px rgba(40,30,20,.08),0 4px 10px rgba(40,30,20,.10);border-radius:1px;display:inline-block;transform:rotate(-1.5deg);position:relative">
         <img src="${e.image_url}" style="display:block;width:160px;height:110px;object-fit:cover;border-radius:1px" />
         <div style="position:absolute;left:0;right:0;bottom:4px;text-align:center;font-family:Caveat,cursive;font-size:14px;color:#5a4a3c;line-height:1">${l(e.name)}</div>
       </div>`:"",r=[];e.city&&r.push(l(e.city)),e.cuisine_type&&r.push(l(e.cuisine_type));const s=r.length?`<div style="font-family:Caveat,cursive;font-size:15px;color:#7a6a5c;margin-top:4px">${r.join(" · ")}</div>`:"",m=t?`<div style="display:inline-flex;align-items:center;background:#ead8a3;padding:4px 9px;border-radius:2px;font-family:Lora,Georgia,serif;margin-top:8px;box-shadow:0 1px 1px rgba(40,30,20,.08)">
         <span style="font-size:8px;letter-spacing:.18em;text-transform:uppercase;color:#7a6a5c;margin-right:6px">общая</span>
         <span style="font-family:Caveat,cursive;font-size:20px;color:#2d231b;line-height:1">${t}</span>
       </div>`:"",c=e.is_gem_place?'<span style="display:inline-block;font-family:Lora,Georgia,serif;font-weight:600;font-size:9.5px;letter-spacing:.18em;text-transform:uppercase;color:#c66b3e;border:1.4px solid #c66b3e;padding:3px 7px 2px;border-radius:2px;background:rgba(232,201,122,.5);margin-left:8px">жемчужина</span>':"",f={terra:{bg:"#dcb19c",ink:"#5a2b1c"},ochre:{bg:"#e8d29a",ink:"#5d4a14"},moss:{bg:"#bdcfb1",ink:"#33472a"},plum:{bg:"#cdb1be",ink:"#4a2438"}},h=["terra","ochre","moss","plum"],d=(e.reviewers||[]).slice(0,4),G=d.length?`<div style="display:inline-flex;align-items:center;margin-top:10px;font-family:Caveat,cursive;font-size:15px;color:#7a6a5c">
         <span style="margin-right:6px">:</span>
         ${d.map((p,M)=>{const k=f[h[Math.abs(p.id||0)%h.length]],C=l((p.username||"?").slice(0,1).toUpperCase()),S=p.avatar_url?`<img src="${p.avatar_url}" alt="" style="width:100%;height:100%;object-fit:cover;display:block" />`:C,j=p.avatar_url?"var(--sb-paper-card)":k.bg;return`<span title="${l(p.username||"")}" style="display:inline-flex;align-items:center;justify-content:center;width:22px;height:22px;border-radius:50%;background:${j};color:${k.ink};font-family:Lora,Georgia,serif;font-weight:600;font-size:10px;box-shadow:0 0 0 2px var(--sb-paper-card),0 1px 2px rgba(40,30,20,.25);margin-left:${M===0?"0":"-6px"};overflow:hidden">${S}</span>`}).join("")}
       </div>`:"";return`<div style="min-width:200px;max-width:260px;font-family:Lora,Georgia,serif;color:#2d231b">
    ${a}
    <div style="font-family:Lora,Georgia,serif;font-style:italic;font-weight:500;font-size:17px;line-height:1.15;color:#2d231b">${l(e.name)}${c}</div>
    ${s}
    <div>${m}</div>
    ${G}
    <div><a href="/places/${e.id}" style="display:inline-block;font-family:Lora,Georgia,serif;font-style:italic;font-size:14px;color:#c66b3e;text-decoration:none;margin-top:10px">подробнее →</a></div>
  </div>`}function l(e){return String(e??"").replace(/&/g,"&amp;").replace(/</g,"&lt;").replace(/>/g,"&gt;").replace(/"/g,"&quot;").replace(/'/g,"&#39;")}function b(){if(!i)return;g.forEach(t=>i.geoObjects.remove(t)),g=[];const e=[];n.places.forEach(t=>{if(!t.lat||!t.lng)return;const a=t.avg_food&&t.avg_service&&t.avg_vibe?((Number(t.avg_food)+Number(t.avg_service)+Number(t.avg_vibe))/3).toFixed(1):null,r=t.id,s=t.is_gem_place,m=s?z(r,a):_(r,a),c=s?44:40,f=s?54:52,h=ymaps.templateLayoutFactory.createClass(`<div style="position:relative;width:${c}px;height:${f}px;">${m}</div>`),d=new ymaps.Placemark([t.lat,t.lng],{hintContent:t.name,balloonContentBody:D(t,a)},{iconLayout:h,iconShape:{type:"Rectangle",coordinates:[[-c/2,-f],[c/2,0]]},iconOffset:[-(c/2),-f],hideIconOnBalloonOpen:!1});d.events.add("click",()=>w("marker-click",t)),i.geoObjects.add(d),g.push(d),e.push([t.lat,t.lng])}),e.length>0&&!n.singleMarker?i.setBounds(ymaps.util.bounds.fromPoints(e),{checkZoomRange:!0,zoomMargin:40,duration:300}).then(()=>{i.getZoom()>15&&i.setZoom(15)}):e.length===1&&i.setCenter(e[0],15,{duration:800})}return N(()=>{ymaps.ready(()=>{i=new ymaps.Map(x.value,{center:n.center,zoom:n.zoom,controls:["zoomControl"]}),b()})}),v(()=>n.places,b,{deep:!0}),v(()=>n.center,e=>{i&&e&&i.setCenter(e,n.zoom,{duration:800})}),O(()=>{i&&(i.destroy(),i=null)}),(e,t)=>(Z(),V("div",{ref_key:"mapContainer",ref:x,class:"sb-map",style:E({height:u.height})},null,4))}},P=B(H,[["__scopeId","data-v-3f1f98dd"]]);export{P as M};
