/** 把 list[i] 往前（-1）或往後（1）移一格；區塊內的順序就是存的順序 */
export function move<T>(list: T[], i: number, delta: -1 | 1) {
  const j = i + delta;
  if (j < 0 || j >= list.length) return;
  [list[i], list[j]] = [list[j] as T, list[i] as T];
}
