import { cache } from "react";
import { getSession, sendRequest } from "./db";
import { employeeSchema, tipSchema } from "./schemas";

export const getTips = cache(async () => {
    try {
      const session = await getSession();
      if (!session) return null;

      const data = await sendRequest(`tips`, session.token, "GET");
      return data ? tipSchema.array().parse(data) : null;
    } catch (error) {
      console.error("Error fetching employee:", error);
      return null;
    }
  });