"use server"
import { cache } from "react";
import { getSession, sendRequest } from "./db";
import { paymentSchema } from "./schemas";

export const getPayments = cache(async () => {
  try {
    const session = await getSession();
    if (!session) return null;

    const data = await sendRequest("payments", session.token, "GET");
    return data ? paymentSchema.array().parse(data) : null;
  } catch (error) {
    console.error("Error fetching employees:", error);
    return null;
  }
});

export const getPaymentsByCustomerId = cache(async (customerId: number) => {
  try {
    const session = await getSession();
    if (!session) return null;
    const data = await sendRequest(`payments/customer/${customerId}`, session.token, "GET");

    if (Array.isArray(data)) {
      return paymentSchema.array().parse(data);
    } else if (data) {
      return [paymentSchema.parse(data)];
    }
    return null;
  } catch (error) {
    console.error("Error fetching payment:", error);
    return null;
  }
});
