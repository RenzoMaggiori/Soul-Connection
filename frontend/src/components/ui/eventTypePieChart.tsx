import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { TrendingUp } from 'lucide-react';
import React from 'react';
import { Event } from '@/db/schemas';
import { Cell, Pie, PieChart, ResponsiveContainer, Tooltip } from 'recharts';

const COLORS = ["#0088FE", "#00C49F", "#FFBB28", "#FF8042", "#FF6347", "#FF4500"];

const countEventsByType = (events: Event[]) => {
    const typeCount: Record<string, number> = {};

    events.forEach((event) => {
        const type = event.Type;
        if (type) {
            if (!typeCount[type]) {
                typeCount[type] = 0;
            }
            typeCount[type]++;
        }
    });

    return Object.keys(typeCount).map((type) => ({
        name: type,
        value: typeCount[type],
    }));
};

const EventTypeChart = ({ events, year, month }: { events: Event[], year: number, month: number }) => {
    const filteredEvents = events.filter((event) => {
        const eventDate = new Date(event.Date);
        return (
            eventDate.getFullYear() === year && eventDate.getMonth() === month
        );
    });
    const eventData = countEventsByType(filteredEvents);
    const monthName = [
        "January", "February", "March", "April", "May", "June",
        "July", "August", "September", "October", "November", "December"
    ];

    return (
        <Card className="text-generic">
            <CardHeader>
                <CardTitle>Event Type Distribution in {monthName[month]}</CardTitle>
                <CardDescription>
                    Distribution of different event types
                </CardDescription>
            </CardHeader>
            <CardContent>
                <ResponsiveContainer width="100%" height={300}>
                    <PieChart>
                        <Pie
                            data={eventData}
                            cx="50%"
                            cy="50%"
                            labelLine={false}
                            outerRadius={80}
                            fill="#8884d8"
                            dataKey="value"
                            label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                        >
                            {eventData.map((entry, index) => (
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
                    Showing event type distribution <TrendingUp className="h-4 w-4" />
                </div>
                <div className="leading-none text-muted-foreground">
                    Data is based on the selected month
                </div>
            </CardFooter>
        </Card>
    );
};

export default EventTypeChart;
