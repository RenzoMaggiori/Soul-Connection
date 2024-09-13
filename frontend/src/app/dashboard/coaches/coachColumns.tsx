"use client";

import { ColumnDef, Table } from "@tanstack/react-table";
import { ArrowUpDown, Delete, Edit, MoreHorizontal } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Customer, Employee } from "@/db/schemas";
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Drawer,
  DrawerContent,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from "@/components/ui/drawer";
import { DataTable } from "@/components/ui/data-table";
import AddEmployeeDialog from "./addEmployeeDialog";
import { Input } from "@/components/ui/input";
import { Dialog } from "@/components/ui/dialog";
import { DialogTrigger } from "@radix-ui/react-dialog";
import { CustomerAssignColumns, CustomerAssignTableFooter, CustomerAssignTableHeader } from "./customerAssignColumns";
import { deleteEmployee } from "@/db/employee";
import EditEmployeeDialog from "./editEmployeeDialog";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";

function getEmployeeCustomerQuantity(employee: Employee, customers: Customer[]) {
  return customers.filter((customer) => customer.Employee_Id === employee.Id).length;
}

export const coachColumns = (customers: Customer[]): ColumnDef<Employee>[] => [
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
    cell: ({ row }) => {
      const employee = row.original
      return (
        <div className="flex items-center justify-end">
          <Avatar>
            <AvatarFallback>
              {employee.Name[0]}{employee.Surname[0]}
            </AvatarFallback>
          </Avatar>
        </div>
      )
    }
  },
  {
    id: "Name",
    accessorFn: (row) => `${row.Name} ${row.Surname}`,
    cell: ({ row }) => {
      const employee = row.original;

      return (
        <div>
          {employee.Name} {employee.Surname}
        </div>
      );
    },
  },
  {
    accessorKey: "Work",
    header: "Work type",
  },
  {
    accessorKey: "Birth_Date",
    header: "Birth Date",
  },
  {
    header: "Customers",
    cell: ({ row }) => {
      const employee = row.original;
      return (
        <Drawer>
          <DrawerTrigger>
            {" "}
            <Button variant="outline">{getEmployeeCustomerQuantity(employee, customers)} customers</Button>
          </DrawerTrigger>
          <DrawerContent>
            <DrawerHeader>
              <DrawerTitle>Assign customers to a employee </DrawerTitle>
            </DrawerHeader>
            <div className="px-8">
              <DataTable
                columns={CustomerAssignColumns}
                data={customers}
                pageSize={5}
                footer={({ table }) => <CustomerAssignTableFooter table={table} />}
                header={({ table }) => <CustomerAssignTableHeader table={table} />}
              />
            </div>
            <DrawerFooter></DrawerFooter>
          </DrawerContent>
        </Drawer>
      );
    },
  },
  {
    id: "actions",
    cell: ({ row }) => {
      const employee = row.original;
      return (
        <Dialog>
          <DropdownMenu>
            <DropdownMenuTrigger>
              <Button variant="ghost" className="h-8 w-8 p-0">
                <span className="sr-only">Open menu</span>
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuLabel>Actions</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DialogTrigger asChild>
                <DropdownMenuItem>
                  <Edit className="mr-2 h-4 w-4" />
                  Edit
                </DropdownMenuItem>
              </DialogTrigger>
              <DropdownMenuItem>
                {" "}
                <Delete className="mr-2 h-4 w-4" />
                <span onClick={() => deleteEmployee(employee.Id)}>Delete</span>
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
          <EditEmployeeDialog employeeData={employee} />
        </Dialog>
      );
    },
  },
];

export const CoachTableHeader = ({ table }: { table: Table<Employee> }) => (
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
    <AddEmployeeDialog />
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

export const CoachTableFooter = ({ table }: { table: Table<Employee> }) => (
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
