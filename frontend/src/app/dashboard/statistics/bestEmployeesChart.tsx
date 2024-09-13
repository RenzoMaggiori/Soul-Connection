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
import { Employee, Encounter } from "@/db/schemas";

export const description = "A bar chart with a label";

const chartConfig = {
  desktop: {
    label: "Rating",
    color: "hsl(var(--chart-1))",
  },
} satisfies ChartConfig;

function getBestEmployees(employees: Employee[], encounters: Encounter[]) {
  const today = new Date();
  const filteredEncounters = encounters.filter(
    (encounter) =>
      new Date(encounter.Date).getFullYear() == new Date().getFullYear(),
  );
  const employeeRatings = employees.map((employee) => {
    const employeeRelatedEncounters = filteredEncounters.filter(
      (encounter) => encounter.Customer_Id === employee.Id,
    );
    const meanRating = getMeanRating(employeeRelatedEncounters);
    return { employee, meanRating };
  });
  employeeRatings.sort((a, b) => b.meanRating - a.meanRating);
  return employeeRatings;
}
function getMeanRating(encounters: Encounter[]): number {
  if (encounters.length === 0) return 0;
  const totalRating = encounters.reduce(
    (sum, encounter) => sum + encounter.Rating,
    0,
  );
  return totalRating / encounters.length;
}

export function BestEmployeesChart({
  employees,
  encounters,
}: {
  employees: Employee[];
  encounters: Encounter[];
}) {
  const bestEmployees = getBestEmployees(employees, encounters);
  const chartData = bestEmployees
    .map((employee) => ({
      name: employee.employee.Name,
      rating: employee.meanRating.toFixed(1),
    }))
    .slice(0, 5);
  return (
    <Card>
      <CardHeader>
        <CardTitle>Best employees of the year</CardTitle>
      </CardHeader>
      <CardContent>
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
              dataKey="name"
              tickLine={false}
              tickMargin={10}
              axisLine={false}
            />
            <ChartTooltip
              cursor={false}
              content={<ChartTooltipContent hideLabel />}
            />
            <Bar dataKey="rating" fill="var(--color-desktop)" radius={8}>
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
          Trending up by 5.2% this month <TrendingUp className="h-4 w-4" />
        </div>
        <div className="leading-none text-muted-foreground">
          Showing total visitors for the last 6 months
        </div>
      </CardFooter>
    </Card>
  );
}
