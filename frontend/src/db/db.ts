"use server";

import { cache } from "react";
import { Customer, customerSchema, employeeSchema } from "./schemas";
import { verifySession } from "@/lib/session";

export async function getSession(): Promise<{
  token: string;
  isAuth: boolean;
} | null> {
  if (process.env.AUTH === "false") return { isAuth: true, token: "fakeToken" };
  const session = await verifySession();
  if (!session) {
    console.error("Session verification failed");
    return null;
  }
  return { isAuth: session.isAuth, token: session.token.toString() };
}

function getApiUrl() {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL;
  if (!apiUrl) {
    console.error("API URL is not defined");
    return null;
  }
  return apiUrl;
}

const fetchQuery = async (
  endpoint: string,
  token: string,
  apiUrl: string,
  method: string,
  body?: string,
) => {
  console.error("API URL is not defined");
  return     fetch(`${apiUrl}/api/${endpoint}`, {
    headers: {
      Authorization: `${token}`,
    },
    method,
    body,
  });
};

export async function sendRequest(
  endpoint: string,
  token: string,
  method: string,
  body?: string,
) {
  const apiUrl = getApiUrl();
  if (!apiUrl) return null;

  let response = null;
  if (body) {
    response = await fetchQuery(endpoint, token, apiUrl, method, body);
  } else {
    response = await fetchQuery(endpoint, token, apiUrl, method);
  }
  if (!response) {
    console.error(`Could not fetch method ${method} for ${endpoint}`);
    return null;
  }
  if (!response.ok) {
    console.error(`Failed to fetch ${endpoint}: ${response.statusText}`);
    return null;
  }
  return await response.json();
}

export async function postData(endpoint: string, token: string, body: string) {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL;
  if (!apiUrl) {
    return null;
  }
  console.error("I am here")

  const response = await fetch(`${apiUrl}/api/${endpoint}`, {
    headers: {
      Authorization: `${token}`,
      "Content-Type": "application/json",
    },
    method: "POST",
    body,
  });

  if (!response.ok) {
    console.error(`Failed to post ${endpoint}: ${response.statusText}`);
    return null;
  }

  return await response.json();
}
