import { ChartConfig, ChartContainer, ChartTooltip, ChartTooltipContent } from '@/components/ui/chart';
import { Customer } from '@/db/schemas';
import React from 'react'
import { Area, Cell, AreaChart, Tooltip, XAxis, YAxis } from "recharts";

const chartConfig = {
    desktop: {
        label: "Rating",
        color: "hsl(var(--chart-1))",
    },
} satisfies ChartConfig;

type AggregatedData = {
    name: string;
    "new customers": number;
};

type CustomerWithJoinDate = Customer & {
    joinDate?: Date;
};

const COLORS = ["#0088FE", "#00C49F", "#FFBB28", "#FF8042"];


function generateJoinDates(customers: Customer[]) {
    const storedData = localStorage.getItem("customersWithJoinDates");

    if (storedData) {
        return JSON.parse(storedData).map((customer: CustomerWithJoinDate) => ({
            ...customer,
            joinDate: customer.joinDate ? new Date(customer.joinDate) : undefined,
        }));
    }
    const startDate = new Date("2023-07-01");
    const endDate = new Date("2024-07-31");
    const timeDiff = endDate.getTime() - startDate.getTime();
    const customersWithDates = customers.map((customer) => {
        const randomTime = startDate.getTime() + Math.random() * timeDiff;
        return { ...customer, joinDate: new Date(randomTime) };
    });

    localStorage.setItem("customersWithJoinDates", JSON.stringify(customersWithDates));

    return customersWithDates;
}

function aggregateCustomerData(customers: CustomerWithJoinDate[]): AggregatedData[] {
    const aggregatedData: Record<string, number> = {};

    customers.forEach((customer) => {
        if (customer.joinDate) {
            const date = customer.joinDate.toISOString().split("T")[0];
            if (!aggregatedData[date]) {
                aggregatedData[date] = 0;
            }
            aggregatedData[date]++;
        }
    });

    return Object.keys(aggregatedData).map((date) => ({
        name: date,
        "new customers": aggregatedData[date],
    }));
}



const customerGraph = ({ customers, monthIndex, year }: { customers: Customer[], monthIndex: number, year: number }) => {
    const customersWithDates = generateJoinDates(customers);
    const customerData = aggregateCustomerData(customersWithDates);
    const filteredCustomerData = customerData.filter((customer) => {
        const eventDate = new Date(customer.name);
        return (
            eventDate.getFullYear() === year && eventDate.getMonth() === monthIndex
        );
    });
    const totalCustomers = filteredCustomerData.length;

    const firstDayOfMonth = new Date(year, monthIndex, 1);
    const lastDayOfMonth = new Date(year, monthIndex + 1, 0);
    const numberOfWeeks =
        Math.ceil((lastDayOfMonth.getDate() - firstDayOfMonth.getDay()) / 7) + 1;

    const weeklyAverage = (totalCustomers / numberOfWeeks).toFixed(2);
    const metrics = [
        { label: "Total", value: customers.length, unit: "" },
        { label: "Monthly", value: totalCustomers, unit: "" },
        { label: "Weekly (avg)", value: weeklyAverage, unit: "" },
    ];
    return (
        <>
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
                <AreaChart
                    accessibilityLayer
                    data={filteredCustomerData}
                    margin={{
                        top: 20,
                    }}
                >
                    <XAxis
                        dataKey="name"
                        tickLine={false}
                        tickMargin={10}
                        axisLine={false}
                    />
                    <YAxis />
                    <Area
                        fill="#8884d8"
                        type="monotone"
                        dataKey="new customers"
                    >
                        <ChartTooltip
                            cursor={false}
                            content={<ChartTooltipContent hideLabel />}
                        />
                        {filteredCustomerData.map((entry, index) => (
                            <Cell
                                key={`cell-${index}`}
                                fill={COLORS[index % COLORS.length]}
                            />
                        ))}
                    </Area>
                    <Tooltip />
                </AreaChart>
            </ChartContainer>
        </>
    )
}

export default customerGraph