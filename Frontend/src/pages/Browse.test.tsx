import { beforeEach, afterEach, describe, expect, it, vi } from "vitest";
import { MemoryRouter } from "react-router-dom";
import { render, screen } from "@testing-library/react";
import Browse from "./Browse";

const mockListClassesRequest = vi.fn();

vi.mock("../lib/api", () => ({
  listClassesRequest: (...args: unknown[]) => mockListClassesRequest(...args),
  getUserEnrollmentsRequest: vi.fn(),
  registerClassRequest: vi.fn(),
  dropClassRequest: vi.fn(),
  createClassRequest: vi.fn(),
  updateClassRequest: vi.fn(),
  deleteClassRequest: vi.fn(),
  listClassEnrollmentsRequest: vi.fn(),
}));

vi.mock("../store/authStore", () => ({
  useAuthStore: () => ({
    isAuthenticated: false,
    token: null,
    userId: null,
    role: "student",
  }),
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

vi.mock("../components/CalendarView", () => ({
  default: () => <div>calendar-view</div>,
}));

vi.mock("../components/CourseDetailsModal", () => ({
  default: () => null,
}));

const makeBackendClass = (overrides?: Partial<Record<string, unknown>>) => ({
  id: 1,
  name: "Morning Flow Yoga",
  course_code: "YOG101",
  description: "Core flow",
  start_time: "10:00:00",
  end_time: "11:00:00",
  capacity: 20,
  duration: 60,
  category: "Yoga",
  weekday: "Thu",
  spot: 8,
  ...overrides,
});

describe("Browse page enrollment window display", () => {
  beforeEach(() => {
    vi.useFakeTimers({ shouldAdvanceTime: true });
    vi.setSystemTime(new Date("2026-04-13T08:00:00"));
    mockListClassesRequest.mockResolvedValue({
      page: 1,
      page_size: 10,
      total: 1,
      classes: [makeBackendClass()],
    });
  });

  afterEach(() => {
    vi.useRealTimers();
    vi.clearAllMocks();
  });

  it("shows 'Not open yet' instead of large countdown when open time is over 30 hours away", async () => {
    render(
      <MemoryRouter>
        <Browse />
      </MemoryRouter>,
    );

    await screen.findByText("Morning Flow Yoga");

    expect(
      screen.getByRole("button", { name: "Not open yet" }),
    ).toBeInTheDocument();
    expect(screen.queryByText(/Opens in \d+h/)).not.toBeInTheDocument();
  });

  it("still shows countdown when open time is within 30 hours", async () => {
    mockListClassesRequest.mockResolvedValueOnce({
      page: 1,
      page_size: 10,
      total: 1,
      classes: [makeBackendClass({ weekday: "Tue" })],
    });

    render(
      <MemoryRouter>
        <Browse />
      </MemoryRouter>,
    );

    await screen.findByText("Morning Flow Yoga");

    expect(
      screen.getByRole("button", { name: /Opens in 1h/i }),
    ).toBeInTheDocument();
  });
});
