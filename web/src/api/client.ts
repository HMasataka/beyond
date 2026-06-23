import createClient from "openapi-fetch";
import type { paths } from "./schema";

const baseUrl = import.meta.env.VITE_API_BASE_URL;
if (!baseUrl) {
  throw new Error("VITE_API_BASE_URL is not set");
}

export const apiClient = createClient<paths>({ baseUrl });
