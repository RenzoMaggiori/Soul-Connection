"use client";
import { useQuery } from "@tanstack/react-query";
import { coachColumns, CoachTableFooter, CoachTableHeader } from "./coachColumns";
import { DataTable } from "../../../components/ui/data-table";
import { Skeleton } from "@/components/ui/skeleton";
import { getEmployees } from "@/db/employee";
import { getCustomers } from "@/db/customer";
import { SkeletonTable } from "@/components/ui/skeletonTable";

export default function AccountManagerPage() {
  const { data, isLoading, isError } = useQuery({
    queryFn: async () => {
      const [employees, customers] = await Promise.all([
        getEmployees(),
        getCustomers(),
      ]);

      return { employees, customers };
    },
    queryKey: ["CoachData"],
    gcTime: 1000 * 60,
  });

  if (isError) return <div>Error...</div>;
  if (data && data.employees === null || data && data.customers === null)
    return <div>Parsing Error...</div>;

  return (
    <div className="container mx-auto py-10">
      {isLoading || !data || !data.customers || !data.employees  ? (
        <SkeletonTable/>
      ) : (
        <DataTable
          columns={coachColumns(data.customers)}
          data={data.employees}
          footer={({ table }) => <CoachTableFooter table={table} />}
          header={({ table }) => <CoachTableHeader table={table} />}
        />
      )}
    </div>
  );
}