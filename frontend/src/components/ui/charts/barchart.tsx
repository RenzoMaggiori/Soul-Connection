"use client";

import { FC } from "react";
import { Bar, BarChart, CartesianGrid, LabelList, XAxis, Tooltip, ResponsiveContainer } from "recharts";
import { ChartContainer, ChartConfig } from "@/components/ui/charts/chart";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";

interface ChartData {
    year: number;
    [key: number]: number;
}

interface BarChartComponentProps {
    title: string;
    description: string;
    data: ChartData[];
    config: ChartConfig;
    dataKey: string;
}

const BarChartComponent: FC<BarChartComponentProps> = ({
    title,
    description,
    data,
    config,
    dataKey,
}) => {
    return (
        <Card>
            <CardHeader>
                <CardTitle style={{ color: 'var(--color-generic)' }}>{title}</CardTitle>
                <CardDescription>{description}</CardDescription>
            </CardHeader>
            <CardContent>
                <ChartContainer config={config}>
                    <ResponsiveContainer width="100%" height="100%">
                        <BarChart data={data} margin={{ top: 20 }}>
                            <CartesianGrid vertical={false} />
                            <XAxis
                                dataKey="year"
                                tickLine={false}
                                tickMargin={10}
                                axisLine={false}
                                tickFormatter={(value) => value.toString()}
                            />
                            <Tooltip cursor={{ fill: 'transparent' }} />
                            <Bar dataKey={dataKey} fill={`var(--color-generic)`} radius={8}>
                                <LabelList position="top" offset={12} className="fill-foreground" fontSize={12} />
                            </Bar>
                        </BarChart>
                    </ResponsiveContainer>
                </ChartContainer>
            </CardContent>
        </Card>
    );
};

export default BarChartComponent;
