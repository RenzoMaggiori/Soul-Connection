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
import { CustomerForm, CustomerFormData } from "@/components/forms/CustomerForm";
import { postCustomer } from "@/db/customer";

function AddCustomerDialog() {
  function onSubmit(values: CustomerFormData) {
    const offset = values.birth_date.getTimezoneOffset()
    const formatedDate = new Date(values.birth_date.getTime() - (offset*60*1000))
    postCustomer({
        Email: values.email,
        Name: values.name,
        Surname: values.surname,
        Birth_Date: formatedDate.toISOString().split('T')[0],
        Phone_Number: values.phone_number,
        Gender: values.gender,
        Description: values.description,
        Address: values.address,
        Astrological_Sign: values.astrological_sign,
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
          Add new Customer
        </Button>
      </DialogTrigger>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>Create a Customer</DialogTitle>
          <DialogDescription>
            Fill in the following information to create a customer.
          </DialogDescription>
        </DialogHeader>
        <CustomerForm onSubmit={onSubmit} />
      </DialogContent>
    </Dialog>
  );
}

export default AddCustomerDialog;