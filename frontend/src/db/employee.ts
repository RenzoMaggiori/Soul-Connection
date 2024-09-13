"use server";
import { cache } from "react";
import { sendRequest, getSession } from "./db";
import { Employee, employeeSchema } from "./schemas";

export const getEmployees = cache(async () => {
  try {
    const session = await getSession();
    if (!session) return null;

    const data = await sendRequest("employees", session.token, "GET");
    return data ? employeeSchema.array().parse(data) : null;
  } catch (error) {
    console.error("Error fetching employees:", error);
    return null;
  }
});

export const getEmployeeById = cache(async (id: number) => {
  try {
    const session = await getSession();
    if (!session) return null;

    const data = await sendRequest(`employees/${id}`, session.token, "GET");
    return data ? employeeSchema.parse(data) : null;
  } catch (error) {
    console.error("Error fetching employee:", error);
    return null;
  }
});

export const postEmployee = cache(async (employee: Partial<Employee>) => {
  try {
    const session = await getSession();
    if (!session) return null;

    const data = await sendRequest("employees", session.token, "POST", JSON.stringify(employee));
    return data ? employeeSchema.parse(data) : null;
  } catch (error) {
    console.error("Error posting employee:", error);
    return null;
  }
});

export const updateEmployee = cache(async (employee: Partial<Employee>, id: number) => {
  try {
    const session = await getSession();
    if (!session) return null;

    const data = await sendRequest(`employees/${id}`, session.token, "PATCH", JSON.stringify(employee));
    return data ? employeeSchema.parse(data) : null;
  } catch (error) {
    console.error("Error updating employee:", error);
    return null;
  }
});

export const deleteEmployee = cache(async (id: number) => {
  try {
    const session = await getSession();
    if (!session) return null;

    const data = await sendRequest(`employees/${id}`, session.token, "DELETE");
    return data ? employeeSchema.parse(data) : null;
  } catch (error) {
    console.error("Error deleting employee:", error);
    return null;
  }
});