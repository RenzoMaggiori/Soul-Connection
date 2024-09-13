"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { getCustomers } from "@/db/customer";
import { CustomerColumns, CustomerTableFooter, CustomerTableHeader } from "./customersColumns";
import { SkeletonTable } from "@/components/ui/skeletonTable";
import { DataTable } from "@/components/ui/data-table";
import { getPayments } from "@/db/payment";

const CustomersPage = () => {
  const { isLoading, isError, data } = useQuery({
    queryFn: async () => {
      const customers = await getCustomers();
      const payments = await getPayments();
      return { customers, payments };
    },
    queryKey: ["CustomerData"],
    gcTime: 1000 * 60, // 1 minute
  });

  if (isError) {
    return <div>Error loading customers...</div>;
  }

  return (
    <div className="container mx-auto py-10">
      {isLoading || !data || !data.customers || !data.payments ? (
        <SkeletonTable />
      ) : (
        <DataTable
          columns={CustomerColumns(data.payments)}
          data={data.customers}
          footer={({ table }) => <CustomerTableFooter table={table} />}
          header={({ table }) => <CustomerTableHeader table={table} />}
        />
      )}
    </div>
  );
};

export default CustomersPage;
