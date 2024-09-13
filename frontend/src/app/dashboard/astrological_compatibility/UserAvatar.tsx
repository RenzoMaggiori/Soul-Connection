import React from "react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
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
import { Skeleton } from "@/components/ui/skeleton";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Customer } from "@/db/schemas";

interface UserAvatarProps {
  user?: string;
  onSelectUser: (value: string) => void;
  data: Customer[] | undefined;
  isLoading: boolean;
  selectedUser: string;
  title: string;
}

const UserAvatar: React.FC<UserAvatarProps> = ({
  user,
  onSelectUser,
  data,
  isLoading,
  selectedUser,
  title,
}) => (
  <div className="flex flex-col items-center gap-5">
    <div className="group relative flex flex-col items-center gap-4">
      <Avatar className="h-24 w-24 md:h-32 md:w-32">
        <AvatarImage src="/placeholder-user.jpg" alt={title} />
        <AvatarFallback>
          {user ? (
            user[0].toUpperCase()
          ) : (
            <Skeleton className="h-4 w-16 bg-zinc-300" />
          )}
        </AvatarFallback>
      </Avatar>
      <Dialog>
        <DialogTrigger asChild>
          <div className="absolute inset-0 flex cursor-pointer items-center justify-center opacity-0 transition-opacity duration-300 group-hover:opacity-100">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-primary text-xl font-bold text-primary-foreground">
              +
            </div>
          </div>
        </DialogTrigger>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>{title}</DialogTitle>
            <DialogDescription>
              Select a user for this profile.
            </DialogDescription>
          </DialogHeader>
          <div className="flex items-center space-x-2">
            <div className="grid flex-1 gap-2">
              {data && !isLoading ? (
                <Select
                  onValueChange={onSelectUser}
                  defaultValue={selectedUser}
                >
                  <SelectTrigger className="w-[180px]">
                    <SelectValue placeholder="Select a customer" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectGroup>
                      <SelectLabel>Customer</SelectLabel>
                      {data.map((customer: Customer) => (
                        <SelectItem
                          value={`${customer.Name} ${customer.Surname}`}
                          key={customer.Id}
                        >
                          {`${customer.Name} ${customer.Surname}`}
                        </SelectItem>
                      ))}
                    </SelectGroup>
                  </SelectContent>
                </Select>
              ) : (
                <Skeleton className="h-6 w-48" />
              )}
            </div>
          </div>
          <DialogFooter className="sm:justify-start">
            <DialogClose asChild>
              <Button type="button" variant="secondary">
                Close
              </Button>
            </DialogClose>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
    {user ? (
      <div className="text-center">
        <h3 className="text-xl font-bold">{user}</h3>
      </div>
    ) : (
      <Skeleton className="h-4 w-32" />
    )}
  </div>
);

export default UserAvatar;
