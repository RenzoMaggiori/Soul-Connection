import React from "react";
import { toast } from "@/hooks/use-toast";
import {
  CustomerForm,
  CustomerFormData,
} from "@/components/forms/CustomerForm";
import {
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Customer } from "@/db/schemas";
import { updateCustomer } from "@/db/customer";

function EditCustomerDialog({ customer }: { customer: Customer }) {
  function onSubmit(values: CustomerFormData) {
    const offset = values.birth_date.getTimezoneOffset();
    const formatedDate = new Date(
      values.birth_date.getTime() - offset * 60 * 1000,
    );
    updateCustomer(
      {
        Email: values.email,
        Name: values.name,
        Surname: values.surname,
        Birth_Date: formatedDate.toISOString().split("T")[0],
        Gender: values.gender,
        Description: values.description,
        Address: values.address,
        Phone_Number: values.phone_number,
        Astrological_Sign: values.astrological_sign,
      },
      customer.Id,
    );
    toast({
      title: "Customer updated successfully",
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
          <DialogTitle>Edit a Customer</DialogTitle>
          <DialogDescription>
            Fill in the following information to edit a Customer.
          </DialogDescription>
        </DialogHeader>
        <CustomerForm onSubmit={onSubmit} initialData={customer} />
        <DialogFooter></DialogFooter>
      </DialogContent>
    </>
  );
}

export default EditCustomerDialog;
