import { useEffect, useCallback } from 'react';

interface CalendarProps {
  calendarRef: React.RefObject<HTMLElement>;
  calendarHeaderRef: React.RefObject<HTMLElement>;
  calendarDays: any[];
}

type CalendarView = 'month' | 'week' | 'day';

function useAdjustCalendar(props: CalendarProps, view: CalendarView) {
  const adjustHeight = useCallback(() => {
    if (props.calendarRef.current) {
      const headerHeight = props.calendarHeaderRef.current?.offsetHeight || 0;
      const remainingHeight = (window.innerHeight * 0.45) - headerHeight - 40;
      let gridAutoRows: number;
      switch (view) {
        case 'month':
          const amountOfRows = Math.ceil(props.calendarDays.length / 7);
          gridAutoRows = remainingHeight / amountOfRows;
          break;
        case 'week':
          gridAutoRows = remainingHeight;
          break;
        case 'day':
          gridAutoRows = remainingHeight / 24;
          break;
        default:
          gridAutoRows = remainingHeight;
      }

      props.calendarRef.current.style.gridAutoRows = `${gridAutoRows}px`;
    }
  }, [props.calendarDays, props.calendarRef, props.calendarHeaderRef, view]);

  useEffect(() => {
    adjustHeight();

    window.addEventListener('resize', adjustHeight);
    window.addEventListener('orientationchange', adjustHeight);
    document.addEventListener('visibilitychange', adjustHeight);
    const timeoutId = setTimeout(adjustHeight, 100);

    return () => {
      window.removeEventListener('resize', adjustHeight);
      window.removeEventListener('orientationchange', adjustHeight);
      document.removeEventListener('visibilitychange', adjustHeight);
      clearTimeout(timeoutId);
    };
  }, [adjustHeight]);
  return adjustHeight;
}

export default useAdjustCalendar;