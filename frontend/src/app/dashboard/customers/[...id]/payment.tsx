import { z } from "zod";

export const paymentSchema = z.object({
    Id: z.number(),
    Soul_Connection_Id: z.number(),
    Date: z.string(),
    PaymentMethod: z.string(),
    Amount: z.number(),
    Comment: z.string(),
    CustomerId: z.number(),
});

export type Payment = z.infer<typeof paymentSchema>;
