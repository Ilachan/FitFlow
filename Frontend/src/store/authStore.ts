import { create } from "zustand";

type Role = "student" | "manager";

type AuthStore = {
  isAuthenticated: boolean;
  token: string | null;
  userId: number | null;
  role: Role | null;
  login: (token: string, role: Role, userId?: number | null) => void;
  logout: () => void;
};

const getInitialState = () => {
  const token = localStorage.getItem("auth_token");
  const role = localStorage.getItem("user_role") as Role;
  const rawUserId = localStorage.getItem("user_id");
  const userId = rawUserId ? Number(rawUserId) : null;
  return {
    isAuthenticated: !!token,
    token,
    userId: Number.isFinite(userId) ? userId : null,
    role: role || null,
  };
};

export const useAuthStore = create<AuthStore>((set) => ({
  ...getInitialState(),
  login: (token: string, role: Role, userId?: number | null) => {
    localStorage.setItem("auth_token", token);
    localStorage.setItem("user_role", role);
    if (typeof userId === "number") {
      localStorage.setItem("user_id", String(userId));
    } else {
      localStorage.removeItem("user_id");
    }
    set({ isAuthenticated: true, token, role, userId: typeof userId === "number" ? userId : null });
  },
  logout: () => {
    localStorage.removeItem("auth_token");
    localStorage.removeItem("user_role");
    localStorage.removeItem("user_id");
    set({ isAuthenticated: false, token: null, userId: null, role: null });
  },
}));
