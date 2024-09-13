import React, { useMemo } from 'react'
import { CalendarDay } from '../utils/calendarUtils';
import { Event } from '@/components/calendar/calendarViews/Event';
import { CalendarViewProps } from '../calendarComponent';
import useAdjustCalendar from '@/hooks/useAdjustCalendarHeight';

interface CalendarWeekProps extends CalendarViewProps {
    weekIndex: number;
}

export default function CalendarWeek(calendarWeekProps: CalendarWeekProps) {
    let { calendarDays, monthIndex, eventsData, eventIndexes, weekIndex } = calendarWeekProps;
    weekIndex = weekIndex < 0 ? calendarDays.length / 7 - 1 : weekIndex;
    const calendarWeeks = useMemo(() => {
        const weeks: CalendarDay[][] = [];
        for (let i = 0; i < calendarDays.length; i += 7) {
            weeks.push(calendarDays.slice(i, i + 7));
        }
        return weeks;
    }, [calendarDays]);

    const currentWeek = calendarWeeks[weekIndex] || [];

    useAdjustCalendar(calendarWeekProps, 'week');

    if (currentWeek.length === 0) {
        console.error(`Week ${weekIndex} not found in calendar data`);
        return <div>Error: Week data not available</div>;
    }
    return (
        <>
            {currentWeek.map((day, index) => (
                <div key={index} className={`day ${day.monthIndex === monthIndex ? '' : 'day--disabled'}`}>
                    {day.date.getDate()}
                </div>
            ))}
            {eventsData.map((event) => (
                <Event
                    key={event.Id}
                    index={eventIndexes.get(event.Id) || 0}
                    event={event}
                    onClick={calendarWeekProps.eventOnClick}
                    calendarDays={currentWeek}
                    calendarRef={calendarWeekProps.calendarRef}
                />
            ))}
        </>
    )
}