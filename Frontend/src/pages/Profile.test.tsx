import { beforeEach, describe, expect, it, vi } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import Profile from "./Profile";

const mockGetProfileRequest = vi.fn();
const mockUpdateProfileRequest = vi.fn();
const mockCreateManagerInviteCodeRequest = vi.fn();

vi.mock("../store/authStore", () => ({
  useAuthStore: () => ({
    role: "student",
    token: "fake-token",
  }),
}));

vi.mock("../lib/api", () => ({
  getProfileRequest: (...args: unknown[]) => mockGetProfileRequest(...args),
  updateProfileRequest: (...args: unknown[]) =>
    mockUpdateProfileRequest(...args),
  createManagerInviteCodeRequest: (...args: unknown[]) =>
    mockCreateManagerInviteCodeRequest(...args),
}));

vi.mock("react-hot-toast", () => {
  const toastFn = vi.fn();
  return {
    default: Object.assign(toastFn, {
      success: vi.fn(),
      error: vi.fn(),
    }),
  };
});

const baseProfile = {
  name: "Alex Johnson",
  email: "alex@example.com",
  avatar_url: "",
  date_of_birth: "2000-01-01",
  gender: "Male",
  phone_number: "(555) 111-2222",
  address: "123 Main St",
};

describe("Profile page", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockGetProfileRequest.mockResolvedValue(baseProfile);
    mockUpdateProfileRequest.mockResolvedValue({ message: "ok" });
  });

  it("shows validation error for invalid avatar URL", async () => {
    render(<Profile />);

    const avatarInput = await screen.findByPlaceholderText(
      "https://example.com/avatar.png",
    );

    await userEvent.clear(avatarInput);
    await userEvent.type(avatarInput, "invalid-avatar-url");

    const saveButtons = screen.getAllByRole("button", { name: "Save Changes" });
    await userEvent.click(saveButtons[saveButtons.length - 1]);

    expect(
      screen.getByText("Avatar URL must start with http:// or https://."),
    ).toBeInTheDocument();
    expect(mockUpdateProfileRequest).not.toHaveBeenCalled();
  });

  it("restores original values when cancel changes is clicked", async () => {
    render(<Profile />);

    const nameInput = (await screen.findByPlaceholderText(
      "Enter full name",
    )) as HTMLInputElement;

    expect(nameInput.value).toBe("Alex Johnson");

    await userEvent.clear(nameInput);
    await userEvent.type(nameInput, "Changed Name");
    expect(nameInput.value).toBe("Changed Name");

    await userEvent.click(
      screen.getByRole("button", { name: "Cancel Changes" }),
    );

    expect(nameInput.value).toBe("Alex Johnson");
  });

  it("opens save confirmation and submits update on confirm", async () => {
    mockGetProfileRequest
      .mockResolvedValueOnce(baseProfile)
      .mockResolvedValueOnce({ ...baseProfile, name: "Updated Alex" });

    render(<Profile />);

    const nameInput = await screen.findByPlaceholderText("Enter full name");
    await userEvent.clear(nameInput);
    await userEvent.type(nameInput, "Updated Alex");

    const saveButtons = screen.getAllByRole("button", { name: "Save Changes" });
    await userEvent.click(saveButtons[saveButtons.length - 1]);

    expect(
      screen.getByText(
        "Are you sure you want to save these changes to your profile?",
      ),
    ).toBeInTheDocument();

    await userEvent.click(screen.getByRole("button", { name: "Confirm Save" }));

    await waitFor(() => {
      expect(mockUpdateProfileRequest).toHaveBeenCalledTimes(1);
    });

    expect(mockUpdateProfileRequest).toHaveBeenCalledWith(
      "fake-token",
      expect.objectContaining({ name: "Updated Alex" }),
    );
  });
});
