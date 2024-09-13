import { z } from 'zod';

// Manager Schema
export const managerSchema = z.object({
  id: z.number(),
  email: z.string().email(),
  password: z.string(),
});

// Employee Schema
export const employeeSchema = z.object({
  Id: z.number(),
  Soul_Connection_Id: z.number().nullable(),
  Name: z.string(),
  Email: z.string(),
  Password: z.string(),
  Surname: z.string(),
  Birth_Date: z.string(),
  Gender: z.string(),
  Work: z.string(),
});

// Event Schema
export const eventSchema = z.object({
  Id: z.number(),
  Name: z.string(),
  Date: z.string(),
  Max_Participants: z.number().int(),
  Location_X: z.string(),
  Location_Y: z.string(),
  Type: z.string(),
  Employee_Id: z.number().nullable(),
});

// Customer Schema
export const customerSchema = z.object({
  Id: z.number(),
  Soul_Connection_Id: z.number().nullable(),
  Email: z.string(),
  Name: z.string(),
  Surname: z.string(),
  Birth_Date: z.string(),
  Phone_Number: z.string(),
  Gender: z.string(),
  Description: z.string(),
  Address: z.string(),
  Astrological_Sign: z.string(),
  Employee_Id: z.number().nullable(),
});

// Payment Schema
export const paymentSchema = z.object({
  Id: z.number(),
  Soul_Connection_Id: z.number(),
  Date: z.string(),
  PaymentMethod: z.string(),
  Amount: z.number(),
  Comment: z.string(),
  CustomerId: z.number(),
});

// Encounter Schema
export const encounterSchema = z.object({
  Id: z.number(),
  Date: z.string(),
  Rating: z.number().int(),
  Comment: z.string(),
  Source: z.string(),
  Customer_Id: z.number().nullable(),
});

// Clothe Schema
export const clotheSchema = z.object({
  Id: z.number(),
  Soul_Connection_Id: z.number(),
  Type: z.string(),
  Image_Id: z.string(),
  CreatedAt: z.string(),
  CustomerId: z.number(),
});

// Tip Schema
export const tipSchema = z.object({
  Id: z.number(),
  Title: z.string(),
  Tip: z.string(),
});

export type Manager = z.infer<typeof managerSchema>;
export type Employee = z.infer<typeof employeeSchema>;
export type Event = z.infer<typeof eventSchema>;
export type Customer = z.infer<typeof customerSchema>;
export type Payment = z.infer<typeof paymentSchema>;
export type Encounter = z.infer<typeof encounterSchema>;
export type Clothe = z.infer<typeof clotheSchema>;
export type Tip = z.infer<typeof tipSchema>;
