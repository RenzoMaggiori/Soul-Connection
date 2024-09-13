import { ToggleGroup, ToggleGroupItem } from "@/components/ui/toggle-group";

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

export interface ToggleCalendarProps {
  view: "month" | "week" | "day";
  setView: React.Dispatch<React.SetStateAction<"month" | "week" | "day">>;
}

export function MobileToggleCalendarView({
  view,
  setView,
}: ToggleCalendarProps) {
  return (
    <Select
      value={view}
      onValueChange={(value) => {
        if (value !== "month" && value !== "week" && value !== "day") return;
        setView(value);
      }}
    >
      <SelectTrigger className="ml-3 w-[90px]">
        <SelectValue placeholder="Theme" />
      </SelectTrigger>
      <SelectContent>
        <SelectItem value="month">Month</SelectItem>
        <SelectItem value="week">Week</SelectItem>
        <SelectItem value="day">Day</SelectItem>
      </SelectContent>
    </Select>
  );
}

export function ToggleCalendarView({ view, setView }: ToggleCalendarProps) {
  return (
    <ToggleGroup
      type="single"
      value={view}
      onValueChange={(value) => {
        if (value !== "month" && value !== "week" && value !== "day") return;
        setView(value);
      }}
    >
      <ToggleGroupItem value="month" aria-label="Toggle month">
        Month
      </ToggleGroupItem>
      <ToggleGroupItem value="week" aria-label="Toggle week">
        Week
      </ToggleGroupItem>
      <ToggleGroupItem value="day" aria-label="Toggle day">
        Day
      </ToggleGroupItem>
    </ToggleGroup>
  );
}
