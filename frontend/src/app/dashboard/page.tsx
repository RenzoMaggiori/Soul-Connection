"use client"
import React, { useState } from "react";
import { Separator } from "@/components/ui/separator";
import { useQuery } from "@tanstack/react-query";
import { getEvents } from "@/db/event";
import { getCustomers } from "@/db/customer";
import {DatePicker} from "antd"
import {Moment} from "moment"
import EventStatistics from "@/components/ui/eventStatistics";

export default function HomePage() {
  const [selectedDate, setSelectedDate] = useState<String | undefined>("2024-07");
  const { data, isLoading, isError } = useQuery({
    queryFn: async () => {
      const [events, customers] = await Promise.all([getEvents(), getCustomers()]);

      return { events, customers };
    },
    queryKey: ["EventChartData"],
    gcTime: 1000 * 60,
  });

  if (isError || !data) return <div>Loading...</div>;
  if (!data.events || !data.customers) return <div>Parsing Error...</div>;

  const handleDateChange = (date: Moment | null) => {
    if (date) {
      const formattedDate = date.format('YYYY-MM');
      setSelectedDate(formattedDate);
    }
  };

  return (
    <>
      <div className="flex flex-row justify-between">
        <h1 className="text-xl md:text-4xl font-semibold text-generic py-2">Dashboard</h1>
        <DatePicker picker="month" onChange={handleDateChange}/>
      </div>
      <h2 className="text-md text-generic text-muted-foreground py-2">Welcome!</h2>
      <Separator className="mt-2 md:mt-4" />
      {isLoading || !data || !data.events ? (
        <div>Loading</div>
      ) : (
        <EventStatistics customers={data.customers} events={data.events} month={selectedDate ? selectedDate : "2024-07"} />
      )}
    </>
  );
}
