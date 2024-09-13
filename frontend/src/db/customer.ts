"use server"
import { cache } from "react";
import { Customer, customerSchema, employeeSchema } from "./schemas";
import { sendRequest, getSession } from "./db";

export const getCustomers = cache(async (): Promise<Customer[] | null> => {
  try {
    const session = await getSession();
    if (!session) return null;

    const data = await sendRequest("customers", session.token, "GET");
    return data ? customerSchema.array().parse(data) : null;
  } catch (error) {
    console.error("Error fetching customers:", error);
    return null;
  }
});

export const getCustomerById = cache(async (id: number) => {
  try {
    const session = await getSession();
    if (!session) return null;

    const data = await sendRequest(`customers/${id}`, session.token, "GET");
    return data ? customerSchema.parse(data) : null;
  } catch (error) {
    console.error("Error fetching customer:", error);
    return null;
  }
});

export const getCustomerEmployeeById = cache(async (id: number) => {
  try {
    const session = await getSession();
    if (!session) return null;
    let data = await sendRequest(`customers/${id}`, session.token, "GET");
    const customer = data ? customerSchema.parse(data) : null;
    if (!customer) return null;
    data = await sendRequest(`employees/${customer.Employee_Id}`, session.token, "GET");
    return data ? employeeSchema.parse(data) : null;
  } catch (error) {
    console.error("Error fetching customer employee:", error);
    return null;
  }
});

// export const getImagesByCustomerId = cache(async (customerId: number) => {
//   try {
//     const session = await getSession();
//     if (!session) return null;
//     const data = await sendRequest(`api/customers/${customerId}/image`, session.token, "GET");
//     console.log("LOG data", data);
//     return data;
//   } catch (error) {
//     console.error("Error fetching images:", error);
//     return null;
//   }
// });

// export const getImagesByCustomerId = cache(async (customerId: number) => {
//   try {
//     const session = await getSession();
//     if (!session) return null;
//     const data = await fetch(`http://localhost:8000/api/customers/${customerId}/image`, {
//       headers: {
//         Authorization: session.token,
//       },
//       method: "GET",
//     });
//     console.log("LOG data", data);
//     return await data.blob();
//   } catch (error) {
//     console.error("Error fetching images:", error);
//     return null;
//   }
// });

export const postCustomer = cache(async (customer: Partial<Customer>) => {
  try {
    const session = await getSession();
    if (!session) return null;
    const data = await sendRequest("customers", session.token, "POST", JSON.stringify(customer));
    return data ? customerSchema.parse(data) : null;
  } catch (error) {
    console.error("Error posting customer:", error);
    return null;
  }
});

export const updateCustomer = cache(async (customer: Partial<Customer>, id: number) => {
  try {
    const session = await getSession();
    if (!session) return null;
    const data = await sendRequest(`customers/${id}`, session.token, "PATCH", JSON.stringify(customer));
    return data ? customerSchema.parse(data) : null;
  } catch (error) {
    console.error("Error updating customer:", error);
    return null;
  }
});

export const deleteCustomer = cache(async (id: number) => {
  try {
    const session = await getSession();
    if (!session) return null;

    const data = await sendRequest(`customers/${id}`, session.token, "DELETE");
    return data ? customerSchema.parse(data) : null;
  } catch (error) {
    console.error("Error deleting customer:", error);
    return null;
  }
});