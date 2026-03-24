import { describe, expect, it } from "vitest";
import { MemoryRouter } from "react-router-dom";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import Login from "./Login";

describe("Login page", () => {
  it("shows validation errors when submitting empty form", async () => {
    render(
      <MemoryRouter>
        <Login />
      </MemoryRouter>,
    );

    await userEvent.click(screen.getByRole("button", { name: "Sign In" }));

    expect(screen.getByText("Email is required.")).toBeInTheDocument();
    expect(screen.getByText("Password is required.")).toBeInTheDocument();
  });

  it("toggles password visibility", async () => {
    render(
      <MemoryRouter>
        <Login />
      </MemoryRouter>,
    );

    const passwordInput = screen.getByPlaceholderText(
      "Password",
    ) as HTMLInputElement;
    const toggleButton = screen.getByRole("button", { name: "Show" });

    expect(passwordInput.type).toBe("password");
    await userEvent.click(toggleButton);
    expect(passwordInput.type).toBe("text");
  });
});
