"use client";

import React, { createContext, useContext, useEffect, useState } from "react";
import { usePathname } from "next/navigation";

interface CurrentPageContextProps {
    currentPage: string;
}

const CurrentPageContext = createContext<CurrentPageContextProps | undefined>(undefined);

export const CurrentPageProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const pathname = usePathname();
    const [currentPage, setCurrentPage] = useState("");

    useEffect(() => {
        if (!pathname) return;

        const pathSplit = pathname.split("/");
        const lastPath = pathSplit[pathSplit.length - 1];
        const pathPrettier = lastPath.split("_").map((word) => word.charAt(0).toUpperCase() + word.slice(1)).join(" ");
        setCurrentPage(pathPrettier);
    }, [pathname]);

    return (
        <CurrentPageContext.Provider value={{ currentPage }}>
            {children}
        </CurrentPageContext.Provider>
    );
};

export const useCurrentPage = () => {
    const context = useContext(CurrentPageContext);
    if (!context) {
        throw new Error("useCurrentPage must be used within a CurrentPageProvider");
    }
    return context;
};
