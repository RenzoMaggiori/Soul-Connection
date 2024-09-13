"use client";

import { FC, useMemo } from "react";
import { RadialBar, RadialBarChart, PolarAngleAxis, ResponsiveContainer, Tooltip } from "recharts";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

interface ChartData {
    coach: string;
    meetings: number;
    fill: string;
}

interface RadialChartProps {
    title: string;
    description: string;
    data: ChartData[];
}

const RadialChart: FC<RadialChartProps> = ({ title, description, data }) => {
    const totalMeetings = useMemo(() => {
        return data.reduce((acc, curr) => acc + curr.meetings, 0);
    }, [data]);

    return (
        <Card className="flex flex-col">
            <CardHeader className="items-center pb-0">
                <CardTitle style={{ color: 'var(--color-generic)' }}>{title}</CardTitle>
                <CardDescription>{description}</CardDescription>
            </CardHeader>
            <CardContent className="flex flex-1 items-center justify-center pb-0">
                <ResponsiveContainer width="100%" height={250}>
                    <RadialBarChart
                        innerRadius="50%"
                        outerRadius="100%"
                        data={data}
                        startAngle={180}
                        endAngle={0}
                        barSize={10}
                    >
                        <PolarAngleAxis type="number" domain={[0, totalMeetings]} angleAxisId={0} tick={false} />
                        <Tooltip cursor={false} />
                        {data.map((entry, index) => (
                            <RadialBar
                                key={index}
                                background
                                dataKey="meetings"
                                fill={entry.fill}
                            />
                        ))}
                        <text
                            x="50%"
                            y="50%"
                            textAnchor="middle"
                            dominantBaseline="middle"
                            className="text-generic text-3xl font-bold"
                        >
                            {totalMeetings.toLocaleString()}
                        </text>
                        <text
                            x="50%"
                            y="60%"
                            textAnchor="middle"
                            dominantBaseline="middle"
                            className="fill-muted-foreground text-sm"
                        >
                            meetings
                        </text>
                    </RadialBarChart>
                </ResponsiveContainer>
            </CardContent>
        </Card>
    );
};

export default RadialChart;
