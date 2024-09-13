"use server";
import { cache } from "react";
import { sendRequest, getSession } from "./db";
import { encounterSchema } from "./schemas";

export const getCustomerEncountersById = cache(async (id: number) => {
  try {
    const session = await getSession();
    if (!session) return null;
    const data = await sendRequest(`/encounters/customer/${id}`, session.token, "GET");
    return data ? encounterSchema.array().parse(data) : null;
  } catch (error) {
    console.error("Error fetching customer encounters:", error);
    return null;
  }
});

export const getEncounters = cache(async () => {
  try {
    const session = await getSession();
    if (!session) return null;
    const data = await sendRequest("encounters", session.token, "GET");
    return data ? encounterSchema.array().parse(data) : null;
  } catch (error) {
    console.error("Error fetching encounters:", error);
    return null;
  }
});

export const getEncountersByCustomerId = cache(async (customerId: number) => {
  try {
    const session = await getSession();
    if (!session) return null;
    const data = await sendRequest(`encounters/customer/${customerId}`, session.token, "GET");

    if (Array.isArray(data)) {
      return encounterSchema.array().parse(data);
    } else if (data) {
      return [encounterSchema.parse(data)];
    }
    return null;
  } catch (error) {
    console.error("Error fetching encounters:", error);
    return null;
  }
});
