import { create } from "zustand";

type Role = "student" | "manager";

type AuthStore = {
  isAuthenticated: boolean;
  token: string | null;
  role: Role | null;
  login: (token: string, role: Role) => void;
  logout: () => void;
};

const getInitialState = () => {
  const token = localStorage.getItem("auth_token");
  const role = localStorage.getItem("user_role") as Role;
  return {
    isAuthenticated: !!token,
    token,
    role: role || null,
  };
};

export const useAuthStore = create<AuthStore>((set) => ({
  ...getInitialState(),
  login: (token: string, role: Role) => {
    localStorage.setItem("auth_token", token);
    localStorage.setItem("user_role", role);
    set({ isAuthenticated: true, token, role });
  },
  logout: () => {
    localStorage.removeItem("auth_token");
    localStorage.removeItem("user_role");
    set({ isAuthenticated: false, token: null, role: null });
  },
}));
