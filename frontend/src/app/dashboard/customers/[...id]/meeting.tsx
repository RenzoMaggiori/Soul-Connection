import { z } from "zod";

export const meetingSchema = z.object({
    Id: z.number().min(1),
    Date: z.string().min(1),
    Rating: z.number().min(0).max(5),
    Comment: z.string().min(1),
    Source: z.string().min(1),
    Customer_Id: z.number().nullable()
});

export type Meeting = z.infer<typeof meetingSchema>;
