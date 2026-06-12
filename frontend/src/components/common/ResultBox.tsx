import { useState, useEffect } from 'react';

function syntaxHighlight(obj: unknown, indent = 0): string {
  const pad = '  '.repeat(indent);
  if (obj === null || obj === undefined) return '<span style="color:#666">null</span>';
  if (typeof obj === 'string') return `<span style="color:#a5d6a7">"${escapeHtml(obj)}"</span>`;
  if (typeof obj === 'number') return `<span style="color:#ffab40">${obj}</span>`;
  if (typeof obj === 'boolean') return `<span style="color:#ce93d8">${obj}</span>`;
  if (Array.isArray(obj)) {
    if (obj.length === 0) return '<span style="color:#666">[</span><span style="color:#666">]</span>';
    const items = obj.map(item => `${pad}  ${syntaxHighlight(item, indent + 1)}`);
    return `<span style="color:#666">[</span>\n${items.join(',\n')}\n${pad}<span style="color:#666">]</span>`;
  }
  if (typeof obj === 'object') {
    const keys = Object.keys(obj as Record<string, unknown>);
    if (keys.length === 0) return '<span style="color:#666">{</span><span style="color:#666">}</span>';
    const items = keys.map(k => {
      const val = (obj as Record<string, unknown>)[k];
      return `${pad}  <span style="color:#7ec8e3">"${escapeHtml(k)}"</span><span style="color:#555">: </span>${syntaxHighlight(val, indent + 1)}`;
    });
    return `<span style="color:#666">{</span>\n${items.join(',\n')}\n${pad}<span style="color:#666">}</span>`;
  }
  return escapeHtml(String(obj));
}

function escapeHtml(s: string): string {
  return s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;').replace(/'/g,'&#39;');
}

interface ResultBoxProps {
  data: unknown;
}

export function ResultBox({ data }: ResultBoxProps) {
  if (data === undefined || data === null) return null;

  let displayData = data;
  if (typeof data === 'object' && data !== null && 'error' in data && (data as any).error === true) {
    const errData = (data as any).data;
    displayData = { error: typeof errData === 'string' ? errData : JSON.stringify(errData) };
  }

  const html = syntaxHighlight(displayData);

  return (
    <div style={{
      background: 'rgba(0,0,0,0.5)',
      border: '1px solid rgba(255,255,255,0.06)',
      borderRadius: 6,
      padding: '14px 16px',
      marginTop: 12,
      maxHeight: 350,
      overflow: 'auto',
      fontFamily: "'JetBrains Mono','Cascadia Code','Consolas',monospace",
      fontSize: 11.5,
      lineHeight: 1.7,
      whiteSpace: 'pre-wrap',
      wordBreak: 'break-all',
      backdropFilter: 'blur(8px)',
    }} dangerouslySetInnerHTML={{ __html: html }} />
  );
}
