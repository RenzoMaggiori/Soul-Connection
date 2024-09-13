"use server";

import { cache } from "react";
import { clotheSchema } from "./schemas";
import { sendRequest, getSession } from "./db";

export const getClothesByCustomerId = cache(async (customerId: number) => {
    try {
        const session = await getSession();
        if (!session) return null;

        console.log("LOG customerId", customerId);
        const data = await sendRequest(`clothes/customer/${customerId}`, session.token, "GET");

        console.log("LOG data", data);
        return data ? clotheSchema.array().parse(data) : null;
    } catch (error) {
        console.error("Error fetching clothes:", error);
        return null;
    }
});
