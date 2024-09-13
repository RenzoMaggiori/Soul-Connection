"use client";

import { ArrowUpDown } from "lucide-react";
import { ColumnDef, Table } from "@tanstack/react-table";
import { Trash2 } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";
import { Task } from "@/app/dashboard/task_management/task";
import { NewTaskDialog } from "@/app/dashboard/task_management/newTaskDialog"

export const taskColumns: ColumnDef<Task>[] = [
  {
    id: "sort",
    header: ({ column }) => {
      return (
        <Button
          className="px-0"
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          #
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      );
    },
  },
  {
    accessorKey: "title",
    header: "Title",
    cell: ({ row }) => {
      return <div>{row.original.title}</div>;
    },
  },
  {
    accessorKey: "description",
    header: "Description",
    cell: ({ row }) => {
      return <div>{row.original.description}</div>;
    },
  },
  {
    accessorKey: "status",
    header: "Status",
    cell: ({ row }) => {
      return <div>{row.original.status}</div>;
    },
  },
  {
    accessorKey: "dueDate",
    header: "Due Date",
    cell: ({ row }) => {
      return <div>{row.original.dueDate.toLocaleDateString()}</div>;
    },
  },
  {
    id: "remove",
    cell: () => {
      return (<Button variant={"destructive"} size={"icon"}>
                <Trash2 className="h-4 w-4"/>
                </Button>);
    },
  }
];

export const TaskTableHeader = ({ table }: { table: Table<Task> }) => (
  <div className="flex items-center py-4 gap-2">
    <Input
      placeholder="Filter titles..."
      value={(table.getColumn("title")?.getFilterValue() as string) ?? ""}
      onChange={(event) =>
        table.getColumn("title")?.setFilterValue(event.target.value)
      }
      className="max-w-sm"
    />
    <div className="ml-5 mr-3 flex-1 text-sm text-muted-foreground">
      {table.getFilteredSelectedRowModel().rows.length} of{" "}
      {table.getFilteredRowModel().rows.length} row(s).
    </div>
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="generic" className="ml-auto max-sm:text-xs">
          Columns
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        {table
          .getAllColumns()
          .filter((column) => column.getCanHide())
          .map((column) => {
            if (column.id === "sort" || column.id === "remove") {
                return
            }
            return (
              <DropdownMenuCheckboxItem
                key={column.id}
                className="capitalize"
                checked={column.getIsVisible()}
                onCheckedChange={(value) => column.toggleVisibility(!!value)}
              >
                {column.id}
              </DropdownMenuCheckboxItem>
            );
          })}
      </DropdownMenuContent>
    </DropdownMenu>
    <NewTaskDialog/>
  </div>
);

export const TaskTableFooter = ({ table }: { table: Table<Task> }) => (
  <div className="flex items-center justify-end space-x-2 py-4">
    <Button
      variant="generic"
      size="sm"
      onClick={() => table.previousPage()}
      disabled={!table.getCanPreviousPage()}
    >
      Previous
    </Button>
    <Button
      variant="generic"
      size="sm"
      onClick={() => table.nextPage()}
      disabled={!table.getCanNextPage()}
    >
      Next
    </Button>
  </div>
);
