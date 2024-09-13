import { NextRequest, NextResponse } from "next/server";
import { cookies } from "next/headers";
import { decrypt } from "@/lib/session";

const publicRoutes = ["/"];

export default async function middleware(req: NextRequest) {
  const path = req.nextUrl.pathname;
  const isProtectedRoute = path.startsWith("/dashboard");
  const isPublicRoute = publicRoutes.includes(path);
  const isAuthActive = (process.env.AUTH ?? "true") !== "false";
  const cookie = cookies().get("session")?.value;
  const session = await decrypt(cookie);

  if (!isAuthActive) {
    return NextResponse.next();
  }

  if (isProtectedRoute && !session?.token) {
    return NextResponse.redirect(new URL("/", req.nextUrl));
  }

  if (
    isPublicRoute &&
    session?.token &&
    !req.nextUrl.pathname.startsWith("/dashboard")
  ) {
    return NextResponse.redirect(new URL("/dashboard", req.nextUrl));
  }

  return NextResponse.next();
}
