import React from "react";
import { updateEmployee } from "@/db/employee";
import { toast } from "@/hooks/use-toast";
import {
  EmployeeForm,
  EmployeeFormData,
} from "@/components/forms/EmployeeForm";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Employee } from "@/db/schemas";

function EditEmployeeDialog({ employeeData }: { employeeData: Employee }) {
  function onSubmit(values: EmployeeFormData) {
    const offset = values.birth_date.getTimezoneOffset()
    const formatedDate = new Date(values.birth_date.getTime() - (offset*60*1000))
    updateEmployee({
      Email: values.email,
      Name: values.name,
      Surname: values.surname,
      Password: values.password,
      Birth_Date: formatedDate.toISOString().split('T')[0],
      Gender: values.gender,
      Work: values.work,
    }, employeeData.Id);
    toast({
      title: "Employee updated successfully",
      description: (
        <pre className="mt-2 w-[340px] rounded-md bg-slate-950 p-4">
          <code className="text-white">{JSON.stringify(values, null, 2)}</code>
        </pre>
      ),
    });
  }

  return (
    <>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>Edit an employee</DialogTitle>
          <DialogDescription>
            Fill in the following information to edit an employee.
          </DialogDescription>
        </DialogHeader>
        <EmployeeForm onSubmit={onSubmit} initialData={employeeData} />
        <DialogFooter>
        </DialogFooter>
      </DialogContent>
    </>
  );
}

export default EditEmployeeDialog;
