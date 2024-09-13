"use client";
import React from "react";
import "./calendar.css";
import { CalendarData, CalendarDay } from "./utils/calendarUtils";
import { useQuery } from "@tanstack/react-query";

import CalendarMonth from "./calendarViews/calendarMonth";
import CalendarWeek from "./calendarViews/calendarWeek";
import { calculateEventIndexes, CalendarEvent } from "./utils/eventUtils";
import { Event } from "@/db/schemas";
import { isDateBetween } from "./utils/dateUtils";
import { useCalendar } from "@/components/calendar/utils/calendar-context";

export interface CalendarHeaderProps {
  calendarHeaderRef: React.RefObject<HTMLDivElement>;
  calendarDays: CalendarDay[];
}

export interface CalendarProps {
  events: Event[];
  eventOnClick: (event: Event) => void;
  Header: (calendarHeaderProps: CalendarHeaderProps) => JSX.Element;
}

export interface CalendarViewProps {
  calendarHeaderRef: React.RefObject<HTMLDivElement>;
  calendarRef: React.RefObject<HTMLDivElement>;
  calendarDays: CalendarDay[];
  monthIndex: number;
  eventsData: CalendarEvent[];
  eventIndexes: Map<number, number>;
  eventOnClick: (event: Event) => void;
}

function getLastDayOfView(
  view: string,
  calendarDays: CalendarDay[],
): CalendarDay {
  switch (view) {
    case "month":
      return calendarDays[calendarDays.length - 1];
    case "week":
      return calendarDays[6];
    case "day":
      return calendarDays[0];
    default:
      return calendarDays[calendarDays.length - 1];
  }
}

export default function Calendar({
  events,
  eventOnClick,
  Header,
}: CalendarProps) {
  const calendarHeaderRef = React.useRef<HTMLDivElement>(null);
  const calendarRef = React.useRef<HTMLDivElement>(null);
  const { view, weekIndex, monthIndex, userYear } = useCalendar();

  const { data, isError, isLoading } = useQuery({
    queryFn: () => CalendarData({ monthIndex, userYear }),
    queryKey: ["calendar", monthIndex, userYear, view],
  });

  const calendarEvents = events
    .filter((event: Event) => {
      const eventDate = new Date(event.Date);
      return isDateBetween(eventDate, monthIndex, userYear);
    })
    .map((event: Event) => {
      const eventDate = new Date(event.Date);
      return {
        ...event,
        start_timestamp: eventDate,
        end_timestamp: eventDate,
      };
    });
  const eventIndexes = calculateEventIndexes(calendarEvents);

  if (isError) return <div>Error...</div>;
  if (!data || isLoading) return <div>Loading...</div>;
  return (
    <div className="calendar-container">
      <Header
        calendarDays={data.calendarDays}
        calendarHeaderRef={calendarHeaderRef}
      />
      <div className="calendar" ref={calendarRef}>
        {["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"].map((day) => (
          <span key={day} className="day-name">
            {day}
          </span>
        ))}
        {view === "month" && (
          <CalendarMonth
            calendarHeaderRef={calendarHeaderRef}
            calendarRef={calendarRef}
            calendarDays={data.calendarDays}
            monthIndex={monthIndex}
            eventsData={calendarEvents}
            eventIndexes={eventIndexes}
            eventOnClick={eventOnClick}
          />
        )}
        {view === "week" && (
          <CalendarWeek
            calendarHeaderRef={calendarHeaderRef}
            calendarRef={calendarRef}
            calendarDays={data.calendarDays}
            monthIndex={monthIndex}
            eventsData={calendarEvents}
            eventIndexes={eventIndexes}
            weekIndex={weekIndex}
            eventOnClick={eventOnClick}
          />
        )}
      </div>
    </div>
  );
}
