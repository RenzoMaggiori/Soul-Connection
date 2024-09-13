import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Customer } from '@/db/schemas'
import { TrendingUp } from 'lucide-react';
import React from 'react'
import { Cell, Pie, PieChart, ResponsiveContainer, Tooltip } from 'recharts'

const COLORS = ["#0088FE", "#00C49F", "#FFBB28", "#FF8042"];

const countCustomersBySign = (customers: Customer[]) => {
    const signCount: Record<string, number> = {};

    customers.forEach((customer) => {
        const sign = customer.Astrological_Sign;
        if (sign) {
            if (!signCount[sign]) {
                signCount[sign] = 0;
            }
            signCount[sign]++;
        }
    });

    return Object.keys(signCount).map((sign) => ({
        name: sign,
        value: signCount[sign],
    }));
};

const customerSignChart = ({ customers }: { customers: Customer[] }) => {
    const astrologicalData = countCustomersBySign(customers);

    return (
        <Card className="text-generic">
            <CardHeader>
                <CardTitle>Astrological Distribution</CardTitle>
                <CardDescription>
                    Customers astrological distribution
                </CardDescription>
            </CardHeader>
            <CardContent>
                <ResponsiveContainer width="100%" height={300}>
                    <PieChart>
                        <Pie
                            data={astrologicalData}
                            cx="50%"
                            cy="50%"
                            labelLine={false}
                            outerRadius={80}
                            fill="#8884d8"
                            dataKey="value"
                            label={({ name, percent }) =>
                                `${name} ${(percent * 100).toFixed(0)}%`
                            }
                        >
                            {astrologicalData.map((entry, index) => (
                                <Cell
                                    key={`cell-${index}`}
                                    fill={COLORS[index % COLORS.length]}
                                />
                            ))}
                        </Pie>
                        <Tooltip />
                    </PieChart>
                </ResponsiveContainer>
            </CardContent>
            <CardFooter className="flex-col items-start gap-2 text-sm">
                <div className="flex gap-2 font-medium leading-none">
                    Showing customers astrological signs <TrendingUp className="h-4 w-4" />
                </div>
                <div className="leading-none text-muted-foreground">
                    Data is based on the information provided by the customers
                </div>
            </CardFooter>
        </Card>
    )
}

export default customerSignChart