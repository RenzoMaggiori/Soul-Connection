"use client";

import React, { FC, useMemo } from "react";
import { TrendingUp } from "lucide-react";
import { Label, Pie, PieChart, Tooltip, ResponsiveContainer } from "recharts";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle,
} from "@/components/ui/card";

interface ChartData {
    coach: string;
    meetings: number;
    fill: string;
}

interface DonutChartComponentProps {
    title: string;
    description: string;
    data: ChartData[];
    dataKey: string;
    config: Record<string, { label: string; color?: string }>;
}

const DonutChartComponent: FC<DonutChartComponentProps> = ({
    title,
    description,
    data,
    dataKey,
    config,
}) => {
    const totalmeetings = useMemo(() => {
        return data.reduce((acc, curr) => acc + curr.meetings, 0);
    }, [data]);

    return (
        <Card className="flex flex-col">
            <CardHeader className="items-center pb-0">
                <CardTitle style={{ color: 'var(--color-generic)' }}>{title}</CardTitle>
                <CardDescription>{description}</CardDescription>
            </CardHeader>
            <CardContent className="flex-1 pb-0">
                <ResponsiveContainer width="100%" height="100%">
                    <PieChart>
                        <Tooltip cursor={false} />
                        <Pie
                            data={data}
                            dataKey={dataKey}
                            nameKey="coach"
                            innerRadius={60}
                            strokeWidth={5}
                        >
                            <Label
                                content={({ viewBox }) => {
                                    if (viewBox && "cx" in viewBox && "cy" in viewBox) {
                                        return (
                                            <text
                                                x={viewBox.cx}
                                                y={viewBox.cy}
                                                textAnchor="middle"
                                                dominantBaseline="middle"
                                            >
                                                <tspan
                                                    x={viewBox.cx}
                                                    y={viewBox.cy}
                                                    className="text-generic text-3xl font-bold"
                                                >
                                                    {totalmeetings.toLocaleString()}
                                                </tspan>
                                                <tspan
                                                    x={viewBox.cx}
                                                    y={(viewBox.cy || 0) + 24}
                                                    className="text-generic"
                                                >
                                                    meetings
                                                </tspan>
                                            </text>
                                        );
                                    }
                                }}
                            />
                        </Pie>
                    </PieChart>
                </ResponsiveContainer>
            </CardContent>
        </Card>
    );
};

export default DonutChartComponent;
