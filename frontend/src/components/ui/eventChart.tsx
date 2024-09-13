"use client";

import { TrendingUp } from "lucide-react";
import { Bar, BarChart, CartesianGrid, LabelList, XAxis } from "recharts";

import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from "@/components/ui/card";
import {
    ChartConfig,
    ChartContainer,
    ChartTooltip,
    ChartTooltipContent,
} from "@/components/ui/chart";
import { Event } from "@/db/schemas";

export const description = "A bar chart with a label";

const chartConfig = {
    desktop: {
        label: "Rating",
        color: "hsl(var(--chart-1))",
    },
} satisfies ChartConfig;

export function EventsChart({ events, month, year }: { events: Event[], month: number, year: number }) {
    const formatDate = (dateString: string) =>
        new Date(dateString).toLocaleDateString("en-GB", {
            day: "2-digit",
            month: "2-digit",
        });

    const filteredEvents = events.filter((event) => {
        const eventDate = new Date(event.Date);
        return (
            eventDate.getFullYear() === year && eventDate.getMonth() === month
        );
    });

    const totalEvents = filteredEvents.length;
    const daysInMonth = new Date(year, month + 1, 0).getDate();
    const dailyAverage = (totalEvents / daysInMonth).toFixed(2);

    const firstDayOfMonth = new Date(year, month, 1);
    const lastDayOfMonth = new Date(year, month + 1, 0);
    const numberOfWeeks =
        Math.ceil((lastDayOfMonth.getDate() - firstDayOfMonth.getDay()) / 7) + 1;

    const weeklyAverage = (totalEvents / numberOfWeeks).toFixed(2);
    const sortedEvents = filteredEvents.sort(
        (a, b) => new Date(a.Date).getTime() - new Date(b.Date).getTime()
    );

    const eventsCountByDate = sortedEvents.reduce<{ [key: string]: number }>(
        (acc, event) => {
            const date = formatDate(event.Date);
            acc[date] = (acc[date] || 0) + 1;
            return acc;
        },
        {}
    );

    const chartData = Object.entries(eventsCountByDate).map(([date, count]) => ({
        date,
        count,
    }));

    const monthName = [
        "January", "February", "March", "April", "May", "June",
        "July", "August", "September", "October", "November", "December"
    ];

    const metrics = [
        { label: "Monthly", value: totalEvents, unit: "" },
        { label: "Weekly (avg)", value: weeklyAverage, unit: "" },
        { label: "Daily (avg)", value: dailyAverage, unit: "" },
    ];

    return (
        <Card className="text-generic">
            <CardHeader>
                <CardTitle>Events in {monthName[month]} </CardTitle>
                <CardDescription>Our events and their status</CardDescription>
            </CardHeader>
            <CardContent>
                <div className="text-md text-center flex justify-between p-2">
                    {metrics.map((metric, index) => (
                        <div key={index} className="flex-col items-center justify-center">
                            <p className="text-slate-600">{metric.label}</p>
                            <strong className="text-lg">
                                {metric.value}
                                {metric.unit}
                            </strong>
                        </div>
                    ))}
                </div>
                <ChartContainer config={chartConfig}>
                    <BarChart
                        accessibilityLayer
                        data={chartData}
                        margin={{
                            top: 20,
                        }}
                    >
                        <CartesianGrid vertical={false} />
                        <XAxis
                            dataKey="date"
                            tickLine={false}
                            tickMargin={10}
                            axisLine={false}
                        />
                        <ChartTooltip
                            cursor={false}
                            content={<ChartTooltipContent hideLabel />}
                        />
                        <Bar dataKey="count" fill="var(--color-desktop)" radius={8}>
                            <LabelList
                                position="top"
                                offset={12}
                                className="fill-foreground"
                                fontSize={12}
                            />
                        </Bar>
                    </BarChart>
                </ChartContainer>
            </CardContent>
            <CardFooter className="flex-col items-start gap-2 text-sm">
                <div className="flex gap-2 font-medium leading-none">
                    Showing events count for September 2023 <TrendingUp className="h-4 w-4" />
                </div>
                <div className="leading-none text-muted-foreground">
                    Data is based on the selected month
                </div>
            </CardFooter>
        </Card>
    );
}
