const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";

type ApiErrorBody = {
  error?: string;
  message?: string;
};

async function parseApiError(response: Response): Promise<string> {
  try {
    const data = (await response.json()) as ApiErrorBody;
    return data.error || data.message || `Request failed with status ${response.status}`;
  } catch {
    return `Request failed with status ${response.status}`;
  }
}

async function postJson<TResponse>(path: string, body: unknown): Promise<TResponse> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(body),
  });

  if (!response.ok) {
    throw new Error(await parseApiError(response));
  }

  return (await response.json()) as TResponse;
}

export type LoginResponse = {
  message: string;
  token: string;
  role_id: number;
};

export const loginRequest = (email: string, password: string) =>
  postJson<LoginResponse>("/auth/login", { email, password });

export const registerStudentRequest = (name: string, email: string, password: string) =>
  postJson<{ message: string }>("/auth/register", { name, email, password });

export const registerManagerRequest = (
  name: string,
  email: string,
  password: string,
  inviteCode: string,
) =>
  postJson<{ message: string }>("/auth/manager/register", {
    name,
    email,
    password,
    invite_code: inviteCode,
  });

export const roleIdToFrontendRole = (roleId: number): "student" | "manager" => {
  // Treat manager and super manager as manager in current frontend role model.
  if (roleId === 2 || roleId === 3) {
    return "manager";
  }
  return "student";
};

export const extractUserIdFromToken = (token: string): number | null => {
  try {
    const parts = token.split(".");
    if (parts.length < 2) {
      return null;
    }

    const payloadBase64Url = parts[1];
    const payloadBase64 = payloadBase64Url.replace(/-/g, "+").replace(/_/g, "/");
    const padded = payloadBase64 + "=".repeat((4 - (payloadBase64.length % 4)) % 4);
    const decoded = atob(padded);
    const payload = JSON.parse(decoded) as { id?: number };

    return typeof payload.id === "number" ? payload.id : null;
  } catch {
    return null;
  }
};
