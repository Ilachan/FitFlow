import { describe, expect, it, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import InstructorProfile from "./InstructorProfile";

vi.mock("react-hot-toast", () => {
  const toastFn = vi.fn();
  return {
    default: Object.assign(toastFn, {
      success: vi.fn(),
      error: vi.fn(),
    }),
  };
});

describe("InstructorProfile page", () => {
  it("updates bio through quick action modal", async () => {
    render(<InstructorProfile />);

    await userEvent.click(
      screen.getByRole("button", { name: "Edit Profile Bio" }),
    );

    const bioInput = screen.getByRole("textbox");
    await userEvent.clear(bioInput);
    await userEvent.type(bioInput, "Updated instructor bio for sprint 3 test.");
    await userEvent.click(screen.getByRole("button", { name: "Save" }));

    expect(
      await screen.findByText("Updated instructor bio for sprint 3 test."),
    ).toBeInTheDocument();
  });

  it("saves updated availability days", async () => {
    render(<InstructorProfile />);

    await userEvent.click(
      screen.getByRole("button", { name: "Manage Availability" }),
    );
    await userEvent.click(screen.getByRole("button", { name: "Tue" }));
    await userEvent.click(screen.getByRole("button", { name: "Mon" }));
    await userEvent.click(screen.getByRole("button", { name: "Save" }));

    expect(await screen.findByText("Tue")).toBeInTheDocument();
    expect(screen.queryByText("Mon")).not.toBeInTheDocument();
  });
});
