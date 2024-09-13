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
  TrendingUpIcon,
  AwardIcon,
  UserIcon,
} from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { getPayments } from "@/db/payment";
import { Employee, Encounter, Payment } from "@/db/schemas";
import { Cell, Pie, PieChart, ResponsiveContainer, Tooltip } from "recharts";
import { BestEmployeesChart } from "./bestEmployeesChart";

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

function MoneyEarnedPerYear(payments: Payment[], year: number) {
  return payments
    .filter((payment) => new Date(payment.Date).getFullYear() == year)
    .reduce((acc, payment) => acc + payment.Amount, 0);
}

function PaymentsThisYear(payments: Payment[], today: Date) {
  return payments.filter(
    (payment) => new Date(payment.Date).getFullYear() == today.getFullYear(),
  );
}

function CalculateEarningsIncreased(payments: Payment[], today: Date) {
  const thisYear = MoneyEarnedPerYear(payments, today.getFullYear());
  const lastYear = MoneyEarnedPerYear(payments, today.getFullYear() - 1);
  return ((thisYear - lastYear) / lastYear) * 100;
}

function GetMostUsedMethodPerYear(payments: Payment[], year: number): string {
  const paymentMethodCount: Record<string, number> = {};
  const sortedPayments = payments.filter((payments) => {
    return new Date(payments.Date).getFullYear() == year;
  });
  payments.forEach((payment) => {
    const method = payment.PaymentMethod;
    if (method) {
      if (!paymentMethodCount[method]) {
        paymentMethodCount[method] = 0;
      }
      paymentMethodCount[method]++;
    }
  });
  let mostPopularMethod = "None";
  let maxCount = 0;
  for (const method in paymentMethodCount) {
    if (paymentMethodCount[method] > maxCount) {
      mostPopularMethod = method;
      maxCount = paymentMethodCount[method];
    }
  }
  return mostPopularMethod;
}

export default function PaymentStatistics({
  employees,
  encounters,
}: {
  employees: Employee[];
  encounters: Encounter[];
}) {
  const today = new Date();
  const { data, isLoading, isError } = useQuery({
    queryFn: async () => {
      const [payments] = await Promise.all([getPayments()]);

      return { payments };
    },
    queryKey: ["CoachData"],
    gcTime: 1000 * 60,
  });

  if (isError) return <div>Error...</div>;

  if (!data || isLoading) {
    return <div>Loading...</div>;
  }

  if (!data.payments || !employees) return <div>Parsing Error...</div>;
  return (
    <div className="container mx-auto p-4">
      <h1 className="mb-6 text-3xl font-bold">Employee Payments Statistics</h1>
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
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Total Payments
            </CardTitle>
            <TrendingUpIcon className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{data.payments?.length}</div>
            <p className="text-xs text-muted-foreground">
              {PaymentsThisYear(data.payments, today).length + " "}payments this
              year
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Money Earned This Year
            </CardTitle>
            <AwardIcon className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {MoneyEarnedPerYear(data.payments, today.getFullYear()).toFixed(
                2,
              )}
              €
            </div>
            <p className="text-xs text-muted-foreground">
              {CalculateEarningsIncreased(data.payments, today).toFixed(2)}%
              from last year
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Most popular payment method
            </CardTitle>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              xmlSpace="preserve"
              width={16}
              height={16}
              viewBox="0 0 502.685 502.685"
            >
              <path
                d="M482.797 276.924c4.53-5.824 6.73-13.331 4.724-20.988L428.05 30.521c-3.451-13.029-16.847-20.837-29.854-17.386L18.184 113.331C5.22 116.761-2.61 130.2.798 143.207L60.269 368.6c3.408 13.007 16.868 20.816 29.876 17.408l134.278-35.419v75.476c0 42.214 69.954 64.303 139.11 64.303 69.113 0 139.153-22.089 139.153-64.302V311.61c-.001-13.741-7.529-25.303-19.889-34.686zm-43.034-77.698 6.212 23.469-75.541 19.953-6.169-23.512 75.498-19.91zM395.931 50.733l11.799 44.695-118.014 31.148-11.799-44.695 118.014-31.148zm-52.956 174.011 6.04 22.951c-27.934 1.251-55.113 6.126-76.943 14.452l-4.616-17.429 75.519-19.974zm-262.991 94.48-6.169-23.426 75.519-19.975 6.212 23.555-75.562 19.846zm90.641-48.987 75.476-19.953 5.716 21.506a95.81 95.81 0 0 0-5.242 3.473l-69.781 18.421-6.169-23.447zm306.866 153.972c0 24.612-50.993 44.544-113.958 44.544-62.9 0-113.937-19.953-113.937-44.544v-27.718c0-.928.539-1.769.69-2.653 3.602 23.34 52.654 41.847 113.247 41.847 60.614 0 109.687-18.508 113.268-41.847.151.884.69 1.726.69 2.653v27.718zm0-54.531c0 24.591-50.993 44.522-113.958 44.522-62.9 0-113.937-19.931-113.937-44.522V341.96c0-.906.539-1.769.69-2.653 3.602 23.318 52.654 41.869 113.247 41.869 60.614 0 109.687-18.551 113.268-41.869.151.884.69 1.747.69 2.653v27.718zM363.532 356.11c-62.9 0-113.937-19.931-113.937-44.501 0-24.569 51.036-44.5 113.937-44.5 62.965 0 113.958 19.931 113.958 44.5.001 24.57-50.993 44.501-113.958 44.501z"
                style={{
                  fill: "#010002",
                }}
              />
            </svg>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {GetMostUsedMethodPerYear(data.payments, today.getFullYear())}
            </div>
            <p className="text-xs text-muted-foreground">
              {GetMostUsedMethodPerYear(data.payments, today.getFullYear() - 1)}{" "}
              was last year most used method
            </p>
          </CardContent>
        </Card>
      </div>

      <div className="mb-8 grid grid-cols-1 gap-8 lg:grid-cols-2">
        <BestEmployeesChart employees={employees} encounters={encounters}/>
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
