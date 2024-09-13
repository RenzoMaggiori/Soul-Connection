"use client";

import { ColumnDef, Table } from "@tanstack/react-table";
import {
  AlertCircle,
  ArrowUpDown,
  Banknote,
  CircleHelp,
  CreditCard,
  Delete,
  DollarSign,
  Edit,
  MoreHorizontal,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { Customer, Employee, Payment } from "@/db/schemas";
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
import PaypalSvg from "@/components/icons/paypalIcon";
import AddCustomerDialog from "./addCustomerDialog";
import { Dialog, DialogTrigger } from "@/components/ui/dialog";
import { deleteCustomer } from "@/db/customer";
import EditCustomerDialog from "./editCustomerDialog";

export const CustomerColumns = (payments: Payment[]): ColumnDef<Customer>[] => [
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
    id: "method",
    header: "Payment Method",
    cell: ({ row }) => {
      const customer = row.original;
      const payment = payments.find(
        (payment) => payment.CustomerId === customer.Id,
      );
      if (!payment) {
        return <span className="text-xs text-red-500">No payment method</span>;
      }
      switch (payment.PaymentMethod) {
        case "PayPal":
          return (
            <div className="flex items-center space-x-2">
              <PaypalSvg className="h-4 w-5 text-blue-500" />
            </div>
          );
        case "Credit Card":
          return (
            <div className="flex items-center space-x-2">
              <CreditCard className="h-5 w-5 text-blue-500" />
            </div>
          );
        case "Bank Transfer":
          return (
            <div className="flex items-center space-x-2">
              <Banknote className="h-5 w-5 text-purple-500" />
            </div>
          );
        default:
          return (
            <div className="flex items-center space-x-2">
              <CircleHelp className="h-4 w-4 text-gray-500" />
            </div>
          );
      }
    },
  },
  {
    id: "actions",
    cell: ({ row }) => {
      const customer = row.original;
      return (
        <Dialog>
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
                <Link href={`/dashboard/customers/${customer.Id}`}>
                  View customer
                </Link>
              </DropdownMenuItem>
              <DialogTrigger asChild>
                <DropdownMenuItem>
                  <Edit className="mr-2 h-4 w-4" />
                  Edit
                </DropdownMenuItem>
              </DialogTrigger>
              <DropdownMenuItem>
                {" "}
                <Delete className="mr-2 h-4 w-4" />
                <span onClick={() => deleteCustomer(customer.Id)}>Delete</span>
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
          <EditCustomerDialog customer={customer} />
        </Dialog>
      );
    },
  },
];

export const CustomerTableHeader = ({ table }: { table: Table<Customer> }) => (
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
    <AddCustomerDialog />
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
  </div>
);

export const CustomerTableFooter = ({ table }: { table: Table<Customer> }) => {
  return (
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
};
