"use client";

import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { Menu } from "lucide-react";
import React, { useState } from "react";
import Image from "next/image";
import { usePathname } from "next/navigation";
import { cn } from "@/utils/utils";

const Navbar: React.FC = () => {
    const [isOpen, setIsOpen] = useState(false);
    const pathname = usePathname();

    const handleToggle = () => {
        setIsOpen(false);
    };

    const isPath = (href: string) => pathname === href;

    const navLinks = [
        { name: "Dashboard", href: "/dashboard" },
        { name: "Coaches", href: "/dashboard/coaches" },
        { name: "Customers", href: "/dashboard/customers" },
        { name: "Tips", href: "/dashboard/advice" },
        { name: "Events", href: "/dashboard/events" },
        { name: "Astrological compatibility", href: "/dashboard/astrological_compatibility" },
        { name: "Wardrobe", href: "/dashboard/wardrobe" },
        { name: "Task Management", href: "/dashboard/task_management" },
    ];

    return (
        <header className="min-h-14 sticky top-0 flex h-16 items-center justify-between border-b bg-background px-4 md:px-6 z-50">
            <div className="flex items-center justify-between w-full md:w-auto">
                <Link
                    href="/dashboard"
                    className="flex items-center gap-2 text-lg font-semibold md:text-base"
                >
                    <Image
                        src="/logoOK.svg"
                        alt="Soul Connection Logo"
                        width={24}
                        height={24}
                        className="h-10 w-20"
                    />
                    <span className="sr-only">Soul Connection</span>
                </Link>
                <Sheet open={isOpen} onOpenChange={setIsOpen}>
                    <SheetTrigger asChild>
                        <Button variant="generic" size="icon" className="ml-auto md:hidden" >
                            <Menu className="h-5 w-5" />
                            <span className="sr-only">Toggle navigation menu</span>
                        </Button>
                    </SheetTrigger>
                    <SheetContent side="left">
                        <nav className="grid gap-6 text-lg font-medium">
                            {navLinks.map((link) => (
                                <Link
                                    key={link.href}
                                    href={link.href}
                                    className={cn(
                                        "transition-colors hover:text-generic",
                                        isPath(link.href)
                                            ? "text-generic border-b-2 border-current"
                                            : "text-muted-foreground"
                                    )}
                                    onClick={handleToggle}
                                >
                                    {link.name}
                                </Link>
                            ))}
                        </nav>
                    </SheetContent>
                </Sheet>
            </div>
            <nav className="hidden md:flex justify-center items-center w-full space-x-6 text-lg font-medium">
                <div className="flex justify-center flex-1 space-x-6">
                    {navLinks.map((link) => (
                        <Link
                            key={link.href}
                            href={link.href}
                            className={cn(
                                "transition-colors hover:text-generic",
                                isPath(link.href)
                                    ? "text-generic border-b-2 border-current"
                                    : "text-muted-foreground"
                            )}
                        >
                            {link.name}
                        </Link>
                    ))}
                </div>
            </nav>
        </header>
    );
};

export default Navbar;
