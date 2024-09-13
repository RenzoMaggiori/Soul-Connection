"use client";

import * as React from "react";
import { cn } from "@/utils/utils";
import { Icons } from "@/components/icons";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { login } from "@/actions/login";

interface UserAuthFormProps extends React.HTMLAttributes<HTMLDivElement> {}

export function UserAuthForm({ className, ...props }: UserAuthFormProps) {
  const [pending, setPending] = React.useState(false);
  const [errors, setErrors] = React.useState<{ email?: string; password?: string }>({});

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setPending(true);

    const formData = new FormData(event.currentTarget);
    const email = formData.get("email")?.toString();
    const password = formData.get("password")?.toString();

    try {
      setErrors({});
      if (!email || !password)
        throw new Error("Email and password are required");
      const result = await login( email, password );
      console.log(result);
    } catch (error: any) {
      setErrors(error);
    } finally {
      setPending(false);
    }
  }

  return (
    <div className={cn("grid gap-6", className)} {...props}>
      <form onSubmit={handleSubmit}>
        <div className="grid gap-2">
          {/* Email input */}
          <div className="grid gap-1">
            <Label className="sr-only" htmlFor="email">
              Email
            </Label>
            <Input
              key={"email"}
              id="email"
              name="email"
              placeholder="name@example.com"
              type="email"
              autoCapitalize="none"
              autoComplete="email"
              autoCorrect="off"
              disabled={pending}
            />
            {errors.email && (
              <p className="text-sm text-red-500">{errors.email}</p>
            )}
          </div>

          {/* Password input */}
          <div className="grid gap-1">
            <Label className="sr-only" htmlFor="password">
              Password
            </Label>
            <Input
              id="password"
              key={"password"}
              name="password"
              placeholder="Type your password"
              type="password"
              autoCapitalize="none"
              autoCorrect="off"
              disabled={pending}
            />
            {errors.password && (
              <p className="text-sm text-red-500">{errors.password}</p>
            )}
          </div>

          {/* Submit button */}
          <Button type="submit" disabled={pending} variant="generic">
            {pending && <Icons.spinner className="mr-2 h-4 w-4 animate-spin" />}
            Login
          </Button>
        </div>
      </form>
    </div>
  );
}
