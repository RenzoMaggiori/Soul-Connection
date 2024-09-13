export function isSameDate(date1: Date, date2: Date): boolean {
  return (
    date1.getDate() === date2.getDate() && date1.getMonth() === date2.getMonth()
  );
}

export function adjustColumnForSunday(day: number): number {
  return day === 0 ? 7 : day;
}

export function isDateBetween(
  date: Date,
  month: number,
  year: number,
): boolean {
    return date.getMonth() === month && date.getFullYear() === year;
}
