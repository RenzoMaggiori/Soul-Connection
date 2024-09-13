"use client";

import Navbar from "@/components/ui/navbar";
import { CurrentPageProvider, useCurrentPage } from "./currentPageContext";

export default function DashboardLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <CurrentPageProvider>
      <div className="flex h-screen w-screen flex-col">
        <Navbar />
        <section className="pt-2 px-4 md:pt-4 md:px-16 h-full w-full overflow-auto bg-secondary">{children}</section>
      </div>
    </CurrentPageProvider>
  );
}

const CurrentPageHeader = () => {
  const { currentPage } = useCurrentPage();

  return <h1 className="text-xl mt-2 ml-4 md:text-2xl md:ml-8 md:mt-4">{currentPage}</h1>;
}
