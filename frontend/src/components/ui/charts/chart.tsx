"use client";

import * as React from "react";
import { Bar, BarChart, CartesianGrid, LabelList, XAxis, Tooltip, ResponsiveContainer } from "recharts";

export interface ChartConfig {
    [key: string]: {
        label?: React.ReactNode;
        color?: string;
    };
}

const ChartContainer: React.FC<{ config: ChartConfig; children: React.ReactNode }> = ({ config, children }) => {
    return <div style={{ width: '100%', height: '400px' }}>{children}</div>;
};

const ChartTooltip: React.FC = (props) => {
    return <Tooltip {...props} />;
};

const ChartTooltipContent: React.FC<{ payload?: any; label?: string }> = ({ payload, label }) => {
    if (!payload || !payload.length) return null;

    return (
        <div className="bg-white p-2 shadow-md rounded-md">
            <p className="font-semibold">{label}</p>
            {payload.map((entry: any, index: number) => (
                <div key={`item-${index}`} className="flex items-center justify-between">
                    <span style={{ color: entry.color }}>{entry.name}:</span>
                    <span>{entry.value}</span>
                </div>
            ))}
        </div>
    );
};

export { ChartContainer, ChartTooltip, ChartTooltipContent };
