import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import ReactQueryProvider from "@/utils/providers";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import Head from "next/head";
import { Toaster } from "@/components/ui/toaster";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Soul Connection",
  description:
    "Soul Connection is a social media platform for spiritual people.",
  icons: {
    icon: "/favicon.ico",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <ReactQueryProvider>
        <body className={inter.className}>
        <Toaster />
          {children}</body>
        {/* <ReactQueryDevtools initialIsOpen={false} /> */}
      </ReactQueryProvider>
    </html>
  );
}
