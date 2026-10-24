import type { SettingsApi } from '#/api';

/** 一則長度警告的說明（文字在 page.settings.warnings.<kind>）；targets 用頓號接起來 */
export function warningText(
  w: SettingsApi.Warning,
  t: (key: string, params?: Record<string, unknown>) => string,
): string {
  return t(`page.settings.warnings.${w.kind}`, {
    chars: w.chars ?? 0,
    count: w.count ?? 0,
    limit: w.limit ?? 0,
    targets: (w.targets ?? []).join('、'),
    top: example.invalid ?? 0,
  });
}

/** GET、PUT 的回應拆成表單內容和警告：伺服器不收不認得的欄位，表單內容不能帶著 warnings 送回去 */
export function splitWarnings<T extends { warnings: SettingsApi.Warning[] }>(
  r: T,
): [Omit<T, 'warnings'>, SettingsApi.Warning[]] {
  const { warnings, ...rest } = r;
  return [rest, warnings];
}
