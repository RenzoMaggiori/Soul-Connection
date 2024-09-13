import React from "react";
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
import { toast } from "@/hooks/use-toast";
import { postEmployee } from "@/db/employee";
import { EmployeeForm, EmployeeFormData } from "@/components/forms/EmployeeForm";

function AddEmployeeDialog() {
  function onSubmit(values: EmployeeFormData) {
    const offset = values.birth_date.getTimezoneOffset()
    const formatedDate = new Date(values.birth_date.getTime() - (offset*60*1000))
    postEmployee({
      Email: values.email,
      Name: values.name,
      Surname: values.surname,
      Password: values.password,
      Birth_Date: formatedDate.toISOString().split('T')[0],
      Gender: values.gender,
      Work: values.work,
    });
    toast({
      title: "You submitted the following values:",
      description: (
        <pre className="mt-2 w-[340px] rounded-md bg-slate-950 p-4">
          <code className="text-white">{JSON.stringify(values, null, 2)}</code>
        </pre>
      ),
    });
  }

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="generic" className="max-text-xs mr-4">
          Add new user
        </Button>
      </DialogTrigger>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>Create an employee</DialogTitle>
          <DialogDescription>
            Fill in the following information to create an employee.
          </DialogDescription>
        </DialogHeader>
        <EmployeeForm onSubmit={onSubmit} />
      </DialogContent>
    </Dialog>
  );
}

export default AddEmployeeDialog;