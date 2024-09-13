"use client";
import React from "react";
import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from "@/components/ui/card";
import {
    TrendingUp,
} from "lucide-react";
import { Customer, Event } from "@/db/schemas";
import { EventsChart } from "./eventChart";
import CustomerGraph from "./customerGraph";
import CustomerSignChart from "./customerSignChart";
import EventTypeChart from "./eventTypePieChart";


export default function EventStatistics({
    customers,
    events,
    month,
}: {
    customers: Customer[];
    events: Event[];
    month: String;
}) {
    const [yearStr, monthStr] = month.split('-');
    const year = parseInt(yearStr, 10);
    const monthIndex = parseInt(monthStr, 10) - 1


    return (
        <div className="container mx-auto p-4">
            <div className="mb-8 grid grid-cols-1 gap-8 lg:grid-cols-2">
                <Card className="text-generic">
                    <CardHeader>
                        <CardTitle>Customers Overview</CardTitle>
                        <CardDescription>
                            When customers have joined in the time.
                        </CardDescription>
                    </CardHeader>
                    <CardContent>
                        <CustomerGraph customers={customers} monthIndex={monthIndex} year={year}/>
                    </CardContent>
                    <CardFooter className="flex-col items-start gap-2 text-sm">
                        <div className="flex gap-2 font-medium leading-none">
                            Showing customers count for September 2023 <TrendingUp className="h-4 w-4" />
                        </div>
                        <div className="leading-none text-muted-foreground">
                            Data is based on the selected month
                        </div>
                    </CardFooter>
                </Card>
                <EventsChart events={events} month={monthIndex} year={year}/>
            </div>
            <div className="mb-8 grid grid-cols-1 gap-8 lg:grid-cols-2">
                <CustomerSignChart customers={customers}/>
                <EventTypeChart events={events} year={year} month={monthIndex}/>
            </div>

        </div>
    );
}
