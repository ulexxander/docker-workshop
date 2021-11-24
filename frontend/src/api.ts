const API_BASE_URL = "/api";

type APIResponse =
  | {
      Data: any;
    }
  | { Error: string };

export async function apiRequest<Data>(
  endpoint: string,
  fetchOptions: RequestInit
) {
  const response = await fetch(API_BASE_URL + endpoint, fetchOptions);
  const decoded: APIResponse = await response.json();
  if ("Error" in decoded) {
    throw new Error(`API ERROR: ${decoded.Error}`);
  }
  return decoded.Data as Data;
}

export type Note = {
  ID: number;
  Text: string;
  CreatedAt: string;
};

export type NoteCreateParams = {
  Text: string;
};
