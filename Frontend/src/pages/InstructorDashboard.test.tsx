import { describe, expect, it, vi } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import InstructorDashboard from "./InstructorDashboard";

vi.mock("react-hot-toast", () => {
  const toastFn = vi.fn();
  return {
    default: Object.assign(toastFn, {
      success: vi.fn(),
      error: vi.fn(),
    }),
  };
});

describe("InstructorDashboard page", () => {
  it("marks all students present for the selected class", async () => {
    render(<InstructorDashboard />);

    await userEvent.click(screen.getByText("Power Vinyasa"));
    await userEvent.click(
      screen.getByRole("button", { name: "Mark All Present" }),
    );

    await waitFor(() => {
      expect(screen.queryAllByText("Pending")).toHaveLength(0);
    });
    expect(screen.getAllByText("Present").length).toBeGreaterThanOrEqual(4);
  });

  it("adds a walk-in student to the selected class roster", async () => {
    render(<InstructorDashboard />);

    await userEvent.click(screen.getByText("Power Vinyasa"));
    await userEvent.click(screen.getByRole("button", { name: "Add Walk-in" }));

    await userEvent.type(
      screen.getByPlaceholderText("e.g. John Doe"),
      "Test Walkin",
    );
    await userEvent.type(
      screen.getByPlaceholderText("john@example.com"),
      "walkin@example.com",
    );
    await userEvent.click(screen.getByRole("button", { name: "Add Student" }));

    expect(await screen.findByText("Test Walkin")).toBeInTheDocument();
  });
});
