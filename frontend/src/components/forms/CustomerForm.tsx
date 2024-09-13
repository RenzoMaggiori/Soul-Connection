import React from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { format } from "date-fns";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { CalendarIcon } from "lucide-react";
import { cn } from "@/utils/utils";
import { Customer, Employee } from "@/db/schemas";
import { isMobilePhone } from "validator";
import { DialogClose, DialogFooter } from "@/components/ui/dialog";

const formSchema = z.object({
  email: z.string().email().min(1),
  name: z.string().max(50).min(1),
  surname: z.string().max(50).min(1),
  birth_date: z.date(),
  phone_number: z.string().max(50).min(1).refine(isMobilePhone),
  gender: z.string().max(50).min(1),
  description: z.string().min(1),
  address: z.string().max(50).min(1),
  astrological_sign: z.string().max(50).min(1),
});

export type CustomerFormData = z.infer<typeof formSchema>;

interface CustomerFormProps {
  onSubmit: (data: CustomerFormData) => void;
  initialData?: Partial<Customer>;
}

export function CustomerForm({ onSubmit, initialData }: CustomerFormProps) {
  if (!initialData) {
    initialData = {
      Email: "",
      Name: "",
      Surname: "",
      Birth_Date: new Date().toDateString(),
      Gender: "",
      Phone_Number: "",
      Description: "",
      Address: "",
      Astrological_Sign: "",
    };
  }
  const astrologicalSign = [
    "Aries",
    "Taurus",
    "Gemini",
    "Cancer",
    "Leo",
    "Virgo",
    "Libra",
    "Scorpio",
    "Sagittarius",
    "Capricorn",
    "Aquarius",
    "Pisces",
  ];
  const form = useForm<CustomerFormData>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: initialData.Email ?? "",
      name: initialData.Name ?? "",
      surname: initialData.Surname ?? "",
      birth_date: new Date(initialData.Birth_Date ?? new Date().toDateString()),
      gender: initialData.Gender ?? "",
      phone_number: initialData.Phone_Number ?? "",
      description: initialData.Description ?? "",
      address: initialData.Address ?? "",
      astrological_sign: initialData.Astrological_Sign ?? "",
    },
  });
  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        <div className="flex gap-4">
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem className="w-1/2">
                <FormLabel>Name</FormLabel>
                <FormControl>
                  <Input placeholder="Alexandra" {...field} />
                </FormControl>
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="surname"
            render={({ field }) => (
              <FormItem className="w-1/2">
                <FormLabel>Surname</FormLabel>
                <FormControl>
                  <Input placeholder="Smith" {...field} />
                </FormControl>
              </FormItem>
            )}
          />
        </div>
        <div className="flex gap-4">
          <FormField
            control={form.control}
            name="address"
            render={({ field }) => (
              <FormItem className="w-1/2">
                <FormLabel>Address</FormLabel>
                <FormControl>
                  <Input placeholder="Saint Petersburg Street" {...field} />
                </FormControl>
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="description"
            render={({ field }) => (
              <FormItem className="w-1/2">
                <FormLabel>Customer description</FormLabel>
                <FormControl>
                  <Input
                    placeholder="Looking for someone to share laughs."
                    {...field}
                  />
                </FormControl>
                <FormMessage></FormMessage>
              </FormItem>
            )}
          />
        </div>
        <div className="flex gap-4">
          <FormField
            control={form.control}
            name="phone_number"
            render={({ field }) => (
              <FormItem className="w-1/2">
                <FormLabel>Phone Number</FormLabel>
                <FormControl>
                  <Input placeholder="612 722 232" {...field} />
                </FormControl>
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="email"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Email</FormLabel>
                <FormControl>
                  <Input placeholder="alexandra@gmail.com" {...field} />
                </FormControl>
              </FormItem>
            )}
          />
        </div>
        <FormField
          control={form.control}
          name="birth_date"
          render={({ field }) => (
            <FormItem className="flex flex-col">
              <FormLabel>Date of birth</FormLabel>
              <Popover>
                <PopoverTrigger asChild>
                  <FormControl>
                    <Button
                      variant={"outline"}
                      className={cn(
                        "w-[240px] pl-3 text-left font-normal",
                        !field.value && "text-muted-foreground",
                      )}
                    >
                      {field.value ? (
                        format(field.value, "PPP")
                      ) : (
                        <span>Pick a date</span>
                      )}
                      <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
                    </Button>
                  </FormControl>
                </PopoverTrigger>
                <PopoverContent className="w-auto p-0" align="start">
                  <Calendar
                    mode="single"
                    selected={field.value}
                    onSelect={field.onChange}
                    disabled={(date) =>
                      date > new Date() || date < new Date("1900-01-01")
                    }
                    initialFocus
                  />
                </PopoverContent>
              </Popover>
            </FormItem>
          )}
        />
        <div className="flex gap-4">
          <FormField
            defaultValue={initialData.Gender}
            control={form.control}
            name="gender"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Gender</FormLabel>
                <Select
                  onValueChange={field.onChange}
                  defaultValue={field.value}
                >
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue placeholder="Select the gender" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    <SelectItem value="Female">Female</SelectItem>
                    <SelectItem value="Male">Male</SelectItem>
                  </SelectContent>
                </Select>
              </FormItem>
            )}
          />
          <FormField
            defaultValue={initialData.Gender}
            control={form.control}
            name="astrological_sign"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Astrological Sign</FormLabel>
                <Select
                  onValueChange={field.onChange}
                  defaultValue={field.value}
                >
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue placeholder="Select the astrological sign" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    {astrologicalSign.map((sign, index) => (
                      <SelectItem key={index} value={sign}>
                        {sign}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </FormItem>
            )}
          />
          <FormMessage />
        </div>
        <div className="flex flex-row gap-3">
          <Button type="submit">Submit</Button>
          <DialogFooter>
            <DialogClose asChild>
              <Button type="button" variant="secondary">
                Close
              </Button>
            </DialogClose>
          </DialogFooter>
        </div>{" "}
      </form>
    </Form>
  );
}
