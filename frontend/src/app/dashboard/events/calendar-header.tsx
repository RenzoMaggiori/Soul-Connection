"use client";

import { CalendarHeaderProps } from "@/components/calendar/calendarComponent";
import {
  MobileToggleCalendarView,
  ToggleCalendarView,
} from "@/components/calendar/toggleCalendarView";
import { useCalendar } from "@/components/calendar/utils/calendar-context";
import useWindowDimensions from "@/hooks/useWindowDimensions";
import { MoveLeftIcon, MoveRightIcon } from "lucide-react";
import { NewEventDialog } from "./newEventDialog";

export function CalendarHeader({
  calendarDays,
  calendarHeaderRef,
}: CalendarHeaderProps) {
  const { width } = useWindowDimensions();
  const {
    view,
    setView,
    setWeekIndex,
    weekIndex,
    setCurrentDate,
    currentDate,
    userYear,
  } = useCalendar();

  const changeMonth = (increment: number) => {
    setCurrentDate((prevDate) => {
      const newDate = new Date(prevDate);
      newDate.setMonth(prevDate.getMonth() + increment);
      return newDate;
    });
  };

  const changeViewIndex = (increment: number) => {
    if (view === "month") {
      setWeekIndex(0);
      changeMonth(increment);
    } else if (view === "week") {
      if (
        weekIndex < calendarDays.length / 7 - 1 &&
        weekIndex + increment >= 0
      ) {
        setWeekIndex((prevIndex) => prevIndex + increment);
      } else {
        if (increment > 0) {
          setWeekIndex(0);
          changeMonth(1);
        } else {
          setWeekIndex(-1);
          changeMonth(-1);
        }
      }
    }
  };

  return (
    <div className="calendar-header mx-2" ref={calendarHeaderRef}>
      {width > 768 ? (
        <ToggleCalendarView view={view} setView={setView} />
      ) : (
        <MobileToggleCalendarView view={view} setView={setView} />
      )}
      <section className="relative flex flex-col items-center justify-center px-2 align-middle">
        <div className="flex flex- items-center">
          <MoveLeftIcon
            size={24}
            className="mr-3 cursor-pointer"
            onClick={() => changeViewIndex(-1)}
          />
          <h1 className="text-lg md:text-xl">
            {currentDate.toLocaleString("default", {
              month: "long",
            })}
          </h1>
          <MoveRightIcon
            size={24}
            className="ml-3 cursor-pointer"
            onClick={() => changeViewIndex(1)}
          />
        </div>
        <p className="text-sm md:text-base">{userYear}</p>
      </section>
      <NewEventDialog />
    </div>
  );
}
