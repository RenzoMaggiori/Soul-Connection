import { Metadata } from "next";
import Image from "next/image";
import Link from "next/link";

import { cn } from "@/utils/utils";
import { buttonVariants } from "@/components/ui/button";
import { UserAuthForm } from "@/components/forms/userAuthForm";

export const metadata: Metadata = {
  title: "Soul Connection",
  description: "Authentication forms built using the components.",
};


export default function AuthenticationPage() {
  return (
    <>
      <div className="container relative grid h-[800px] min-h-[100vh] flex-col items-center justify-center lg:max-w-none lg:grid-cols-2 lg:px-0">
        <div className="relative hidden h-full flex-col items-center justify-center bg-[#D3D8E8] text-white lg:flex">
          <div className="relative z-20 items-center">
            <Image
              src="/logo.png"
              alt="Soul Connection Logo"
              style={{ height: "25rem", width: "25rem" }}
              width={400}
              height={400}
            />
          </div>
        </div>
        <div className="lg:p-8">
          <div className="mx-auto flex w-[350px] flex-col justify-center space-y-6">
            <div className="flex flex-col space-y-2 text-center">
              <h1 className="text-2xl text-generic font-semibold tracking-tight">Login</h1>
              <p className="text-sm text-muted-foreground">
                Enter your email and password below to login
              </p>
            </div>
            <UserAuthForm />
            <p className="px-8 text-center text-sm text-muted-foreground">
              By clicking continue, you agree to our{" "}
              <Link
                href="#"
                className="underline underline-offset-4 hover:text-primary"
              >
                Terms of Service
              </Link>{" "}
              and{" "}
              <Link
                href="#"
                className="underline underline-offset-4 hover:text-primary"
              >
                Privacy Policy
              </Link>
              .
            </p>
          </div>
        </div>
      </div>
    </>
  );
}
