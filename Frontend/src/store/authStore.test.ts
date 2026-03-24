import { beforeEach, describe, expect, it } from "vitest";
import { useAuthStore } from "./authStore";

describe("authStore", () => {
  beforeEach(() => {
    localStorage.clear();
    useAuthStore.setState({
      isAuthenticated: false,
      token: null,
      userId: null,
      role: null,
    });
  });

  it("stores auth state and localStorage on login", () => {
    const { login } = useAuthStore.getState();

    login("token-123", "student", 88);

    const state = useAuthStore.getState();
    expect(state.isAuthenticated).toBe(true);
    expect(state.token).toBe("token-123");
    expect(state.role).toBe("student");
    expect(state.userId).toBe(88);

    expect(localStorage.getItem("auth_token")).toBe("token-123");
    expect(localStorage.getItem("user_role")).toBe("student");
    expect(localStorage.getItem("user_id")).toBe("88");
  });

  it("removes user_id when login has no userId", () => {
    localStorage.setItem("user_id", "123");
    const { login } = useAuthStore.getState();

    login("token-456", "manager");

    expect(localStorage.getItem("user_id")).toBeNull();
    expect(useAuthStore.getState().userId).toBeNull();
  });

  it("clears auth state and localStorage on logout", () => {
    const { login, logout } = useAuthStore.getState();
    login("token-789", "student", 7);

    logout();

    const state = useAuthStore.getState();
    expect(state.isAuthenticated).toBe(false);
    expect(state.token).toBeNull();
    expect(state.role).toBeNull();
    expect(state.userId).toBeNull();

    expect(localStorage.getItem("auth_token")).toBeNull();
    expect(localStorage.getItem("user_role")).toBeNull();
    expect(localStorage.getItem("user_id")).toBeNull();
  });
});
