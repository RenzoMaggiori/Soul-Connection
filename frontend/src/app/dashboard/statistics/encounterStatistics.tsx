"use client";
import React from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
} from "recharts";
import {
  ArrowUpIcon,
  ArrowDownIcon,
  UserIcon,
  TrendingUpIcon,
  AwardIcon,
} from "lucide-react";
import { Employee, Encounter } from "@/db/schemas";
import { useQuery } from "@tanstack/react-query";
import { getEncounters } from "@/db/encounter";

const employeeData = [
  { name: "John", successRate: 85, projectsCompleted: 23, rating: 4.5 },
  { name: "Emma", successRate: 92, projectsCompleted: 28, rating: 4.8 },
  { name: "Michael", successRate: 78, projectsCompleted: 19, rating: 4.2 },
  { name: "Sophia", successRate: 88, projectsCompleted: 25, rating: 4.6 },
  { name: "William", successRate: 81, projectsCompleted: 21, rating: 4.3 },
];

const departmentData = [
  { name: "Sales", value: 30 },
  { name: "Marketing", value: 25 },
  { name: "Engineering", value: 35 },
  { name: "HR", value: 10 },
];

const COLORS = ["#0088FE", "#00C49F", "#FFBB28", "#FF8042"];

function CalculateAverageEncounterRating(encounters: Encounter[], today: Date) {
  const encountersThisYear = encounters.filter(
    (encounter) =>
      new Date(encounter.Date).getFullYear() == today.getFullYear(),
  );
  return encountersThisYear.reduce((acc, encounter) => acc + encounter.Rating, 0) /
    encountersThisYear.length;
}
function CalculateRatingIncreased(encounters: Encounter[], today: Date) {

  const thisYear = CalculateAverageEncounterRating(encounters, today);
  const lastYearDate = new Date(today.getFullYear() - 1);
  const lastYear = CalculateAverageEncounterRating(encounters, lastYearDate);
  return ((thisYear - lastYear) / lastYear) * 100;
}

export default function EncounterStatistics({
  employees,
  encounters,
}: {
  employees: Employee[];
  encounters: Encounter[];
}) {
  const today = new Date();
  // const { data, isLoading, isError } = useQuery({
  //   queryFn: async () => {
  //     const [ encounters ] = await Promise.all([
  //       getEncounters(),
  //     ]);

  //     return { encounters };
  //   },
  //   queryKey: ["CoachData"],
  //   gcTime: 1000 * 60,
  // });

  // if (isError) return <div>Error...</div>;

  // if (!data || isLoading) {
  //   return <div>Loading...</div>;
  // }
  // if (!data.encounters) return <div>No data...</div>;

  return (
    <div className="container mx-auto p-4">
      <h1 className="mb-6 text-3xl font-bold">
        Employee Encounters Statistics
      </h1>

      <div className="mb-8 grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Total Employees
            </CardTitle>
            <UserIcon className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{employees.length}</div>
            <p className="text-xs text-muted-foreground">
              +12% from last month
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Average Success Rate
            </CardTitle>
            <TrendingUpIcon className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">84.8%</div>
            <p className="text-xs text-muted-foreground">
              +2.3% from last quarter
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Projects Completed
            </CardTitle>
            <AwardIcon className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">1,284</div>
            <p className="text-xs text-muted-foreground">+8% from last year</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Average Rating
            </CardTitle>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              className="h-4 w-4 text-muted-foreground"
            >
              <path d="M12 2L15.09 8.26L22 9.27L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9.27L8.91 8.26L12 2Z" />
            </svg>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {CalculateAverageEncounterRating(encounters, today).toFixed(2)}
            </div>
            <p className="text-xs text-muted-foreground">
              {CalculateAverageEncounterRating(encounters, today).toFixed(2)}% from last year
            </p>
          </CardContent>
        </Card>
      </div>

      <div className="mb-8 grid grid-cols-1 gap-8 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Employee Success Rates</CardTitle>
            <CardDescription>Top 5 employees by success rate</CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={employeeData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="name" />
                <YAxis />
                <Tooltip />
                <Legend />
                <Bar
                  dataKey="successRate"
                  fill="#8884d8"
                  name="Success Rate (%)"
                />
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Department Distribution</CardTitle>
            <CardDescription>
              Employee distribution across departments
            </CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={departmentData}
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
                  {departmentData.map((entry, index) => (
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
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Top Performing Employees</CardTitle>
          <CardDescription>
            Based on success rate and projects completed
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-8">
            {employeeData.map((employee) => (
              <div key={employee.name} className="flex items-center">
                <div className="space-y-1">
                  <p className="text-sm font-medium leading-none">
                    {employee.name}
                  </p>
                  <p className="text-sm text-muted-foreground">
                    {employee.projectsCompleted} projects completed
                  </p>
                </div>
                <div className="ml-auto font-medium">
                  {employee.successRate}%
                </div>
                <Progress
                  value={employee.successRate}
                  className="ml-4 w-[60px]"
                />
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
