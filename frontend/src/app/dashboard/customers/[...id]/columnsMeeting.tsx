import { ColumnDef } from "@tanstack/react-table";
import { Meeting } from "./meeting";
import { Star, StarOff } from "lucide-react";

export const meetingColumns: ColumnDef<Meeting>[] = [
  {
    accessorKey: "Date",
    header: "Date",
    cell: ({ row }) => <div>{row.original.Date}</div>,
  },
  {
    accessorKey: "Rating",
    header: "Rating",
    cell: ({ row }) => (
      <div className="flex items-center">
        {Array(row.original.Rating).fill(true).map((_, i) => (
          <Star key={i} className="h-4 w-4" />
        ))}
        {Array(5 - row.original.Rating).fill(true).map((_, i) => (
          <StarOff key={i} className="h-4 w-4" />
        ))}
      </div>
    ),
  },
  {
    accessorKey: "Report",
    header: "Report",
    cell: ({ row }) => <div>{row.original.Comment}</div>,
  },
  {
    accessorKey: "Source",
    header: "Source",
    cell: ({ row }) => <div>{row.original.Source}</div>,
  },
];
