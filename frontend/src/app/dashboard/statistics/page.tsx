'use client'
import React from "react";
import DonutChartComponent from "@/components/ui/charts/donutchart";
import RadialChart from "@/components/ui/charts/radialChart";
import HorizontalBarChart from "@/components/ui/charts/horizontalbarchart";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import EncounterStatistics from "./encounterStatistics";
import PaymentStatistics from "./paymentStatistics";
import { useQuery } from "@tanstack/react-query";
import { getEmployees } from "@/db/employee";
import { getEncounters } from "@/db/encounter";

export default function ProfilePage({ params }: { params: { id: string } }) {
  const { data, isLoading, isError } = useQuery({
    queryFn: async () => {
      const [employees, encounters] = await Promise.all([getEmployees(), getEncounters()]);

      return { employees, encounters };
    },
    queryKey: ["ChartData"],
    gcTime: 1000 * 60,
  });

  if (isError) return <div>Error...</div>;
  if (data && (!data.employees || !data.encounters)) return <div>Parsing Error...</div>;

  return (
    <section className="">
      {isLoading || !data || !data.employees  ? (
        <div>Loading</div>
      ) : (
        <Tabs defaultValue="payments" className="">
          <TabsList>
            <TabsTrigger value="encounters">Encounters</TabsTrigger>
            <TabsTrigger value="payments">Payments</TabsTrigger>
          </TabsList>
          <TabsContent value="encounters">
            <EncounterStatistics employees={data.employees} encounters={data.encounters!}/>
          </TabsContent>
          <TabsContent value="payments">
            <PaymentStatistics employees={data.employees} encounters={data.encounters!}/>
          </TabsContent>
        </Tabs>
      )}
      {/* <div className="mt-8 flex justify-center">
        <div className="grid w-full max-w-6xl grid-cols-1 gap-8 md:grid-cols-3">
          <RadialChart
            title="Total Meetings"
            description="June 2024"
            data={chartData}
          />
          <HorizontalBarChart
            title="Total Meetings"
            description="June 2024"
            data={chartData}
            config={chartConfig}
          />
          <DonutChartComponent
            title="Total Meetings"
            description="June 2024"
            data={chartData}
            dataKey="meetings"
            config={chartConfig}
          />
        </div>
      </div> */}
    </section>
  );
}
