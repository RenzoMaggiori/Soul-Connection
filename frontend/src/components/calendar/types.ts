import { Event } from "@/db/schemas";
import { CalendarDay } from "./utils/calendarUtils";
import { CalendarEvent } from "./utils/eventUtils";

export interface EventProps {
    event: CalendarEvent;
    calendarDays: CalendarDay[];
    index: number;
    calendarRef: React.RefObject<HTMLDivElement>;
    onClick: (event: Event) => void;
}

export interface EventDisplayInfo {
    isStarting: boolean;
    isEnding: boolean;
    startRow: number;
    endRow: number;
    index: number;
    columnStart: number;
    columnSpan: number;
    className: string;
    children: string;
}