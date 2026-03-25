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

async function getJson<TResponse>(path: string): Promise<TResponse> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  });

  if (!response.ok) {
    throw new Error(await parseApiError(response));
  }

  return (await response.json()) as TResponse;
}

async function authRequest<TResponse>(
  path: string,
  method: "GET" | "PUT" | "POST" | "DELETE",
  token: string,
  body?: unknown,
): Promise<TResponse> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    method,
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: body ? JSON.stringify(body) : undefined,
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

export const registerStudentRequest = (
  name: string,
  email: string,
  password: string,
) => postJson<{ message: string }>("/auth/register", { name, email, password });

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

export const roleIdToFrontendRole = (
  roleId: number,
): "student" | "manager" | "supermanager" => {
  if (roleId === 2) {
    return "supermanager";
  }
  if (roleId === 3) {
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
    const payloadBase64 = payloadBase64Url
      .replace(/-/g, "+")
      .replace(/_/g, "/");
    const padded =
      payloadBase64 + "=".repeat((4 - (payloadBase64.length % 4)) % 4);
    const decoded = atob(padded);
    const payload = JSON.parse(decoded) as { id?: number };

    return typeof payload.id === "number" ? payload.id : null;
  } catch {
    return null;
  }
};

export type UserProfile = {
  name: string;
  email: string;
  avatar_url: string;
  date_of_birth: string;
  gender: string;
  phone_number: string;
  address: string;
};

export type UpdateUserProfilePayload = {
  name: string;
  email: string;
  avatar_url?: string;
  date_of_birth?: string;
  phone_number?: string;
  address?: string;
  gender?: string;
};

export const getProfileRequest = (token: string) =>
  authRequest<UserProfile>("/auth/profile", "GET", token);

export const updateProfileRequest = (
  token: string,
  profile: UpdateUserProfilePayload,
) => authRequest<{ message: string }>("/auth/profile", "PUT", token, profile);

export type CreateManagerInviteCodeRequest = {
  expire_hours: number;
  invitee_email?: string;
};

export type CreateManagerInviteCodeResponse = {
  message: string;
  code: string;
};

export const createManagerInviteCodeRequest = (
  token: string,
  payload: CreateManagerInviteCodeRequest,
) =>
  authRequest<CreateManagerInviteCodeResponse>(
    "/auth/manager/invite-codes",
    "POST",
    token,
    payload,
  );

export type BackendClass = {
  id: number;
  name: string;
  course_code: string;
  description: string;
  start_time: string;
  end_time: string;
  capacity: number;
  duration: number;
  category: string;
  weekday: string;
  spot: number;
};

export type ListClassesResponse = {
  page: number;
  page_size: number;
  total: number;
  classes: BackendClass[];
};

export const listClassesRequest = (page = 1) =>
  getJson<ListClassesResponse>(`/classes?page=${page}`);

export const registerClassRequest = (token: string, courseId: number) =>
  authRequest<{ message: string }>("/classes/register", "POST", token, {
    course_id: courseId,
  });

export const dropClassRequest = (token: string, courseId: number) =>
  authRequest<{ message: string }>("/classes/drop", "POST", token, {
    course_id: courseId,
  });

export type UserEnrollmentsResponse = {
  courses: BackendClass[];
};

export const getUserEnrollmentsRequest = (token: string, userId: number) =>
  authRequest<UserEnrollmentsResponse>(
    `/users/${userId}/enrollments`,
    "GET",
    token,
  );

export type ClassUpsertRequest = {
  name: string;
  course_code: string;
  description: string;
  start_time: string;
  end_time: string;
  capacity: number;
  duration: number;
  category: string;
  weekday: string;
};

export type ClassUpsertResponse = {
  class: BackendClass;
};

export const createClassRequest = (
  token: string,
  payload: ClassUpsertRequest,
) => authRequest<ClassUpsertResponse>("/classes", "POST", token, payload);

export const updateClassRequest = (
  token: string,
  classId: number,
  payload: ClassUpsertRequest,
) =>
  authRequest<ClassUpsertResponse>(
    `/classes/${classId}`,
    "PUT",
    token,
    payload,
  );

export const deleteClassRequest = (token: string, classId: number) =>
  authRequest<{ message: string }>(`/classes/${classId}`, "DELETE", token);

export type UserAnalyticsResponse = {
  analytics: {
    user_id: number;
    range: "7d" | "1m" | "3m";
    from_date: string;
    to_date: string;
    total_classes: number;
    active_days: number;
    daily: Array<{
      date: string;
      classes: number;
    }>;
    categories: Array<{
      category: string;
      classes: number;
      percentage: number;
    }>;
  };
};

export const getUserAnalyticsRequest = (
  token: string,
  userId: number,
  range: "7d" | "1m" | "3m",
) =>
  authRequest<UserAnalyticsResponse>(
    `/users/${userId}/analytics?range=${range}`,
    "GET",
    token,
  );
