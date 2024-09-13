"use client";

import { ColumnDef, Table } from "@tanstack/react-table";
import { ArrowUpDown, MoreHorizontal } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { Customer, Employee } from "@/db/schemas";
import { Input } from "@/components/ui/input";
import { useState } from "react";
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import Link from "next/link";

export const CustomerAssignColumns: ColumnDef<Customer>[] = [
  {
    id: "select",
    header: ({ table }) => (
      <Checkbox
        checked={
          table.getIsAllPageRowsSelected() ||
          (table.getIsSomePageRowsSelected() && "indeterminate")
        }
        onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
        aria-label="Select all"
      />
    ),
    cell: ({ row }) => (
      <Checkbox
        checked={row.getIsSelected()}
        onCheckedChange={(value) => row.toggleSelected(!!value)}
        aria-label="Select row"
      />
    ),
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "Id",
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
    id: "Name",
    accessorFn: (row) => `${row.Name} ${row.Surname}`,
  },
  {
    accessorKey: "Birth_Date",
    header: "Birth Date",
  },
  {
    id: "actions",
    cell: ({ row }) => {
      const customer = row.original;
      return (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="h-8 w-8 p-0">
              <span className="sr-only">Open menu</span>
              <MoreHorizontal className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuLabel>Actions</DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem>
              <Link href={`/dashboard/customers/${customer.Id}`}>View customer</Link>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      );
    },
  },
];

export const CustomerAssignTableHeader = ({ table }: { table: Table<Customer> }) => (
  <div className="flex items-center py-4">
    <Input
      placeholder="Filter by name..."
      value={(table.getColumn("Name")?.getFilterValue() as string) ?? ""}
      onChange={(event) =>
        table.getColumn("Name")?.setFilterValue(event.target.value)
      }
      className="max-w-sm"
    />
    <div className="ml-5 mr-3 flex-1 text-sm text-muted-foreground">
      {table.getFilteredSelectedRowModel().rows.length} of{" "}
      {table.getFilteredRowModel().rows.length} row(s).
    </div>
  </div>
);

export const CustomerAssignTableFooter = ({ table }: { table: Table<Customer> }) => {
  const [isSubmitted, setIsSubmitted] = useState(false);

  const handleSubmit = () => {
    table.getSelectedRowModel().rows.map((row) => console.log(row.original));
    setIsSubmitted(true);
  };

  return (
    <div className="flex flex-row justify-between space-y-0">
      {/* Left buttons */}
      <div className="flex justify-start space-x-2 py-2">
        <Button
          variant="generic"
          onClick={handleSubmit}
          disabled={table.getSelectedRowModel().rows.length <= 0 || isSubmitted}
        >
          {isSubmitted ? "Submitted" : "Submit"}
        </Button>
      </div>

      {/* Pagination buttons */}
      <div className="flex justify-end space-x-2 py-2">
        <Button
          variant="generic"
          onClick={() => table.previousPage()}
          disabled={!table.getCanPreviousPage()}
        >
          Prev
        </Button>
        <Button
          variant="generic"
          onClick={() => table.nextPage()}
          disabled={!table.getCanNextPage()}
        >
          Next
        </Button>
      </div>
    </div>
  );
};
