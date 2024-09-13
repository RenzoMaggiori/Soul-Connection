"use server";
import { FormState, LoginFormSchema } from "@/lib/definitions";
import { createSession, deleteSession } from "@/lib/session";
import { redirect } from 'next/navigation';

export async function postLogin(email: string, password: string) {
  try {
    const response = await fetch("http://localhost:8000/api/auth/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email, password }),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    // Extract relevant information from headers
    const authToken = response.headers.get('Authorization');

    if (!authToken) {
      throw new Error('Missing authentication information in response headers');
    }

    return { authToken };
  } catch (error) {
    console.error("Error sending login to API");
    throw error;
  }
}

export async function login(
  email: string,
  password: string,
): Promise<FormState> {

  const validatedFields = LoginFormSchema.safeParse({
    email: email,
    password: password,
  });

  if (!validatedFields.success) {
    return {
      errors: validatedFields.error.flatten().fieldErrors,
    };
  }
  let result;
  try {
    result = await postLogin(
      validatedFields.data.email,
      validatedFields.data.password,
    );
  } catch (error) {
    return { message: "Invalid email or password." };
  }
  await createSession(result.authToken);
}

export async function logout() {
  deleteSession();
}