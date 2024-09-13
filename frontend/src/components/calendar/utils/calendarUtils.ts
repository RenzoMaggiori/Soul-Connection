export interface CalendarDay {
  date: Date;
  isCurrentMonth: boolean;
  monthIndex: number;
}

function getDayName(dateStr: string, locale: string) {
  var date = new Date(dateStr);
  return date.toLocaleDateString(locale, { weekday: 'long' });
}

function getDaysInMonth(year: number, month: number): CalendarDay[] {
  const firstDayOfMonth = new Date(year, month, 1);
  const lastDayOfMonth = new Date(year, month + 1, 0);
  const daysInMonth: CalendarDay[] = [];

  let startDay = firstDayOfMonth.getDay();
  startDay = startDay === 0 ? 6 : startDay - 1;

  for (let i = startDay - 1; i >= 0; i--) {
    const date = new Date(year, month, -i);
    daysInMonth.push({ date, isCurrentMonth: false, monthIndex: date.getMonth() });
  }
  for (let day = 1; day <= lastDayOfMonth.getDate(); day++) {
    const date = new Date(year, month, day);
    daysInMonth.push({ date, isCurrentMonth: true, monthIndex: date.getMonth() });
  }

  let remainingDays = 42 - daysInMonth.length;
  for (let i = 1; i <= remainingDays; i++) {
    const date = new Date(year, month + 1, i);
    daysInMonth.push({ date, isCurrentMonth: false, monthIndex: date.getMonth() });
    if (getDayName(date.toISOString(), 'en-US') === 'Sunday')
      break;
  }
  return daysInMonth;
}

export async function CalendarData({
  monthIndex,
  userYear,
}: {
  monthIndex: number;
  userYear: number;
}) {
  const calendarDays = getDaysInMonth(userYear, monthIndex);

  return { calendarDays };
}