"use client";

import { FC } from "react";
import { Bar, BarChart, XAxis, YAxis, ResponsiveContainer, Tooltip } from "recharts";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { ChartConfig, ChartContainer, ChartTooltip, ChartTooltipContent } from "@/components/ui/charts/chart";

interface ChartData {
    coach: string;
    meetings: number;
    fill: string;
}

interface HorizontalBarChartProps {
    title: string;
    description: string;
    data: ChartData[];
    config: ChartConfig;
}

const HorizontalBarChart: FC<HorizontalBarChartProps> = ({ title, description, data, config }) => {
    return (
        <Card className="flex flex-col">
            <CardHeader className="items-center pb-4">
                <CardTitle style={{ color: 'var(--color-generic)' }}>{title}</CardTitle>
                <CardDescription>{description}</CardDescription>
            </CardHeader>
            <CardContent>
                <ChartContainer config={config}>
                    <ResponsiveContainer width="100%" height="100%">
                        <BarChart
                            data={data}
                            layout="vertical"
                            margin={{ left: 0 }}
                        >
                            <YAxis
                                dataKey="coach"
                                type="category"
                                tickLine={false}
                                tickMargin={10}
                                axisLine={false}
                                tickFormatter={(value) => value.toString()}
                            />
                            <XAxis dataKey="meetings" type="number" hide />
                            <Tooltip cursor={false} />
                            <Bar dataKey="meetings" radius={5} />
                        </BarChart>
                    </ResponsiveContainer>
                </ChartContainer>
            </CardContent>
        </Card>
    );
};

export default HorizontalBarChart;
