import React, { useState, useEffect, useCallback } from "react";
import { EventProps } from "../types";
import { getEventDisplayInfo } from "../utils/eventUtils";
import "../calendar.css";
import { cn } from "@/utils/utils";

export const Event: React.FC<EventProps> = ({
  event,
  calendarDays,
  index,
  calendarRef,
  onClick,
}) => {
  const [maxEventsPerRow, setMaxEventsPerRow] = useState<number>(0);
  const [eventHeight, setEventHeight] = useState<number>(27);
  const eventProps = getEventDisplayInfo(event, calendarDays, index);
  const topPadding = 30;

  const calculateMaxEvents = useCallback(() => {
    if (calendarRef.current) {
      const rowHeight = parseInt(
        calendarRef.current.style.gridAutoRows || "0",
        10,
      );
      const availableHeight = rowHeight - topPadding;
      if (availableHeight < 27 && availableHeight > 10) {
        setEventHeight(availableHeight);
        setMaxEventsPerRow(1);
      } else setMaxEventsPerRow(Math.floor(availableHeight / eventHeight));
    }
  }, [calendarRef]);

  useEffect(() => {
    calculateMaxEvents();
    window.addEventListener("resize", calculateMaxEvents);

    return () => {
      window.removeEventListener("resize", calculateMaxEvents);
    };
  }, [calculateMaxEvents]);

  const getEventStyle = (row: number): React.CSSProperties => {
    const topPosition =
      topPadding +
      Math.min(eventProps.index, maxEventsPerRow - 1) * eventHeight;

    const style: React.CSSProperties = {
      gridRow: row,
      top: `${topPosition}px`,
      background: "#fef0db",
    };

    if (eventProps.index >= maxEventsPerRow) {
      style.display = "none";
    }

    const sideColor = "#fdb44d";

    if (
      eventProps.startRow === eventProps.endRow &&
      eventProps.isStarting &&
      eventProps.isEnding
    ) {
      eventProps.columnSpan -= eventProps.columnStart - 1;
    }

    if (row === eventProps.startRow && eventProps.isStarting) {
      style.borderLeftWidth = "3px";
      style.borderLeftColor = sideColor;
      style.gridColumn = `${eventProps.columnStart} / span 7`;
    }

    if (row === eventProps.endRow && eventProps.isEnding) {
      style.borderRightWidth = "3px";
      style.borderRightColor = sideColor;
      style.gridColumn =
        row !== eventProps.startRow
          ? `1 / span ${eventProps.columnSpan}`
          : `${eventProps.columnStart} / span ${eventProps.columnSpan}`;
    }

    if (
      (row !== eventProps.startRow || !eventProps.isStarting) &&
      (row !== eventProps.endRow || !eventProps.isEnding)
    ) {
      style.gridColumn = "1 / span 7";
    }

    return style;
  };

  const elements = Array.from(
    { length: eventProps.endRow - eventProps.startRow + 1 },
    (_, i) => {
      const row = eventProps.startRow + i;
      return (
        <section
          key={row}
          className={cn("task cursor-pointer", eventProps.className)}
          style={getEventStyle(row)}
          onClick={() => onClick(event)}
        >
          <div className="overflow-hidden">
            {row === eventProps.startRow ? eventProps.children : null}
          </div>
        </section>
      );
    },
  );

  return <>{elements}</>;
};
