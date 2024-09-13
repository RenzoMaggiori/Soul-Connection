import React from 'react'
import { Event } from '@/components/calendar/calendarViews/Event';
import { CalendarViewProps } from '../calendarComponent';
import useAdjustCalendar from '@/hooks/useAdjustCalendarHeight';


interface CalendarMonth extends CalendarViewProps { };

export default function CalendarMonth(calendarMonthProps: CalendarViewProps) {
    useAdjustCalendar(calendarMonthProps, 'month');
    return (
        <>
            {calendarMonthProps.calendarDays.map((day, index) => (
                <div key={index} className={`day ${day.monthIndex === calendarMonthProps.monthIndex ? '' : 'day--disabled'}`}>
                    {day.date.getDate()}
                </div>
            ))}
            {calendarMonthProps.eventsData.map((event) => (
                <Event
                    key={event.Id}
                    index={calendarMonthProps.eventIndexes.get(event.Id) || 0}
                    event={event}
                    onClick={calendarMonthProps.eventOnClick}
                    calendarDays={calendarMonthProps.calendarDays}
                    calendarRef={calendarMonthProps.calendarRef}
                />
            ))}
        </>
    )
}
