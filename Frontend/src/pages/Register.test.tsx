import { describe, expect, it } from "vitest";
import { MemoryRouter } from "react-router-dom";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import Register from "./Register";

describe("Register page", () => {
  it("shows validation errors for empty required fields", async () => {
    render(
      <MemoryRouter>
        <Register />
      </MemoryRouter>,
    );

    await userEvent.click(
      screen.getByRole("button", { name: "Create Account" }),
    );

    expect(screen.getByText("Full name is required.")).toBeInTheDocument();
    expect(screen.getByText("Email is required.")).toBeInTheDocument();
    expect(screen.getByText("Password is required.")).toBeInTheDocument();
  });

  it("updates password strength indicator as user types", async () => {
    render(
      <MemoryRouter>
        <Register />
      </MemoryRouter>,
    );

    const passwordInput = screen.getByPlaceholderText("Password");

    expect(screen.getByText("weak")).toBeInTheDocument();

    await userEvent.type(passwordInput, "Abcdefg1");

    expect(screen.getByText("strong")).toBeInTheDocument();
  });
});
