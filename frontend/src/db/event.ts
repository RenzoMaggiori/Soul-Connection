"use server";
import { cache } from "react";
import { Event, eventSchema } from "@/db/schemas";
import { getSession, sendRequest } from "@/db/db";

export const getEvents = cache(async (): Promise<Event[] | null> => {
  try {
    const session = await getSession();
    if (!session) {
      return null;
    }
    const data = await sendRequest("events", session.token, "GET");
    return data ? eventSchema.array().parse(data) : null;
  } catch (error) {
    console.error("Error fetching events");
    return null;
  }
});
