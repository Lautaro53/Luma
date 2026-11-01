
export const themeIds = [
  'zen',
  'autumn',
  'forest',
  'midnight',
  'techno',
  'ocean',
  'lavender',
  'graphite'
] as const;

export type ThemeId = (typeof themeIds)[number];

const STORAGE_KEY = 'luma.theme';

export function loadTheme(): ThemeId {
  if (typeof window === 'undefined') {
    return 'zen';
  }

  try {
    const stored = localStorage.getItem(STORAGE_KEY);

    return themeIds.find((id) => id === stored) ?? 'zen';
  } catch {
    return 'zen';
  }
}
