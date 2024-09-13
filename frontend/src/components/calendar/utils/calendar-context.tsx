"use client";

import React from "react";

interface CalendarContextProps {
  view: "month" | "week" | "day";
  currentDate: Date;
  weekIndex: number;
  monthIndex: number;
  userYear: number;
  setView: React.Dispatch<React.SetStateAction<"month" | "week" | "day">>;
  setCurrentDate: React.Dispatch<React.SetStateAction<Date>>;
  setWeekIndex: React.Dispatch<React.SetStateAction<number>>;
}

const CalendarContext = React.createContext<CalendarContextProps | undefined>(
  undefined,
);

export function CalendarProvider({ children }: { children: React.ReactNode }) {
  const [currentDate, setCurrentDate] = React.useState<Date>(() => new Date());
  const [view, setView] = React.useState<"month" | "week" | "day">("month");
  const [weekIndex, setWeekIndex] = React.useState<number>(0);
  const { monthIndex, userYear } = React.useMemo(
    () => ({
      monthIndex: currentDate.getMonth(),
      userYear: currentDate.getFullYear(),
    }),
    [currentDate],
  );

  return (
    <CalendarContext.Provider
      value={{
        currentDate,
        setCurrentDate,
        view,
        setView,
        weekIndex,
        setWeekIndex,
        monthIndex,
        userYear,
      }}
    >
      {children}
    </CalendarContext.Provider>
  );
}

export const useCalendar = () => {
  const context = React.useContext(CalendarContext);
  if (!context) {
    throw new Error("useCalendar must be used within a CalendarProvider");
  }
  return context;
};
