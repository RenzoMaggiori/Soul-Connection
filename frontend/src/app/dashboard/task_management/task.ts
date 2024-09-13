import { z } from "zod";

export const taskSchema = z.object({
    title: z.string().min(1),
    description: z.string().min(1),
    status: z.union([z.literal("todo"), z.literal("doing"), z.literal("done")]),
    dueDate: z.date(),
});

export type Task = z.infer<typeof taskSchema>;
