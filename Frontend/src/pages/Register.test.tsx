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

  it("shows password policy error for invalid password", async () => {
    render(
      <MemoryRouter>
        <Register />
      </MemoryRouter>,
    );

    await userEvent.type(screen.getByPlaceholderText("Full Name"), "Test User");
    await userEvent.type(
      screen.getByPlaceholderText("Email Address"),
      "test@example.com",
    );
    await userEvent.type(screen.getByPlaceholderText("Password"), "abc");
    await userEvent.click(
      screen.getByRole("button", { name: "Create Account" }),
    );

    expect(
      screen.getByText(
        "Password must be 8+ chars and include uppercase, lowercase, and a number.",
      ),
    ).toBeInTheDocument();
  });
});
