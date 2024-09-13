import { CalendarDay } from "./calendarUtils";
import { EventDisplayInfo } from "../types";
import { isSameDate, adjustColumnForSunday } from "./dateUtils";
import { Event } from "@/db/schemas";

export interface CalendarEvent  extends Event {
    start_timestamp: Date;
    end_timestamp: Date;
}

export function getEventDisplayInfo(event: CalendarEvent, calendarDays: CalendarDay[], index: number): EventDisplayInfo {
    let { startRow, endRow } = findEventRows(event, calendarDays);
    let { columnStart, columnSpan } = calculateColumns(event);
    let isStarting = true;
    let isEnding = true;

    if (endRow === 0) {
        isEnding = false;
        endRow = Math.floor(calendarDays.length / 7) + 1;
    }

    if (startRow === 0) {
        isStarting = false;
        startRow = 2;
        columnStart = 1;
    }

    return {
        isStarting,
        isEnding,
        startRow,
        endRow,
        index,
        columnStart,
        columnSpan,
        className: 'task--warning',
        children: event.Name
    };
}

function findEventRows(event: CalendarEvent, calendarDays: CalendarDay[]): { startRow: number, endRow: number } {
    let startRow = 0;
    let endRow = 0;
    for (let i = 0; i < calendarDays.length; i++) {
        const calendarDay = calendarDays[i];
        if (isSameDate(calendarDay.date, event.start_timestamp)) {
            startRow = Math.floor(i / 7) + 2;
        }
        if (isSameDate(calendarDay.date, event.end_timestamp)) {
            endRow = Math.floor(i / 7) + 2;
        }
    }
    return { startRow, endRow };
}

function calculateColumns(event: CalendarEvent): { columnStart: number, columnSpan: number} {
    let columnStart = adjustColumnForSunday(event.start_timestamp.getDay());
    let columnSpan = adjustColumnForSunday(event.end_timestamp.getDay());
    return { columnStart, columnSpan };
}

export function calculateEventIndexes(events: CalendarEvent[]): Map<number, number> {
    const sortedEvents = [...events].sort((a, b) =>
      a.start_timestamp.getTime() - b.start_timestamp.getTime() ||
      b.end_timestamp.getTime() - a.end_timestamp.getTime()
    );

    const indexMap = new Map<number, number>();
    let activeEvents: CalendarEvent[] = [];

    sortedEvents.forEach(event => {
      activeEvents = activeEvents.filter(e => (
        e.start_timestamp.getTime() <= event.end_timestamp.getTime() &&
        e.end_timestamp.getTime() >= event.start_timestamp.getTime()
      ));
      let index = 0;
      while (activeEvents.some(e => indexMap.get(e.Id) === index)) {
        index++;
      }
      indexMap.set(event.Id, index);
      activeEvents.push(event);
    });

    return indexMap;
  }