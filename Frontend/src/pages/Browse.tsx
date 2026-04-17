import { useState, useMemo, useEffect, useCallback } from "react";
import Card from "../components/ui/Card";
import Button from "../components/ui/Button";
import Badge from "../components/ui/Badge";
import Modal from "../components/ui/Modal";
import { Icons } from "../lib/icons";
import CalendarView from "../components/CalendarView";
import CourseDetailsModal, {
  type CourseCardItem,
} from "../components/CourseDetailsModal";
import {
  createClassRequest,
  deleteClassRequest,
  listClassesRequest,
  listClassEnrollmentsRequest,
  registerClassRequest,
  dropClassRequest,
  getUserEnrollmentsRequest,
  updateClassRequest,
  type ClassEnrollmentItem,
  type BackendClass,
  type ClassUpsertRequest,
} from "../lib/api";
import toast from "react-hot-toast";
import { useAuthStore } from "../store/authStore";
import { useNavigate } from "react-router-dom";

const getCourseEmoji = (category: string) => {
  const key = (category || "").toLowerCase();
  if (key.includes("yoga")) return "🧘‍♀️";
  if (key.includes("cardio") || key.includes("hiit")) return "🔥";
  if (key.includes("pilates") || key.includes("strength")) return "💪";
  if (key.includes("dance") || key.includes("zumba")) return "💃";
  return "🏋️";
};

const normalizeWeekday = (weekday: string) => {
  const day = (weekday || "").trim();
  return day.length <= 3 ? day : day.slice(0, 3);
};

const formatTime = (time: string) => {
  const [hourText = "0", minuteText = "00"] = (time || "").split(":");
  const hour24 = Number(hourText);
  if (!Number.isFinite(hour24)) {
    return "--:--";
  }

  const suffix = hour24 >= 12 ? "PM" : "AM";
  const hour12 = hour24 % 12 === 0 ? 12 : hour24 % 12;
  return `${String(hour12).padStart(2, "0")}:${minuteText.slice(0, 2)} ${suffix}`;
};

const toInputTime = (time: string) => {
  const [hourText = "00", minuteText = "00"] = (time || "").split(":");
  const hour = hourText.padStart(2, "0").slice(0, 2);
  const minute = minuteText.padStart(2, "0").slice(0, 2);
  return `${hour}:${minute}`;
};

const toMinutes = (time: string) => {
  const [hourText = "0", minuteText = "0"] = time.split(":");
  const hour = Number(hourText);
  const minute = Number(minuteText);
  if (!Number.isFinite(hour) || !Number.isFinite(minute)) {
    return null;
  }
  return hour * 60 + minute;
};

const fromMinutes = (totalMinutes: number) => {
  const normalized = Math.max(
    0,
    Math.min(23 * 60 + 59, Math.floor(totalMinutes)),
  );
  const hour = Math.floor(normalized / 60);
  const minute = normalized % 60;
  return `${String(hour).padStart(2, "0")}:${String(minute).padStart(2, "0")}`;
};

const WEEKDAY_TO_INDEX: Record<string, number> = {
  sun: 0,
  mon: 1,
  tue: 2,
  wed: 3,
  thu: 4,
  fri: 5,
  sat: 6,
};

const LONG_COUNTDOWN_THRESHOLD_HOURS = 30;

const formatDuration = (totalMinutes: number) => {
  if (!Number.isFinite(totalMinutes) || totalMinutes < 0) {
    return "";
  }

  const hours = Math.floor(totalMinutes / 60);
  const minutes = Math.floor(totalMinutes % 60);
  if (hours === 0) {
    return `${minutes}m`;
  }
  if (minutes === 0) {
    return `${hours}h`;
  }
  return `${hours}h ${minutes}m`;
};

const parseTimeParts = (time: string) => {
  const [hourText = "0", minuteText = "0"] = (time || "").split(":");
  const hour = Number(hourText);
  const minute = Number(minuteText);
  if (!Number.isFinite(hour) || !Number.isFinite(minute)) {
    return null;
  }
  return { hour, minute };
};

const getNextSessionStart = (course: CourseCardItem) => {
  const dayKey = normalizeWeekday(course.day).toLowerCase();
  const dayIndex = WEEKDAY_TO_INDEX[dayKey];
  const parsedTime = parseTimeParts(course.startTimeRaw || "");

  if (typeof dayIndex !== "number" || !parsedTime) {
    return null;
  }

  const now = new Date();
  const nextStart = new Date(now);
  nextStart.setHours(parsedTime.hour, parsedTime.minute, 0, 0);

  const daysAhead = (dayIndex - nextStart.getDay() + 7) % 7;
  nextStart.setDate(nextStart.getDate() + daysAhead);

  if (daysAhead === 0 && now.getTime() > nextStart.getTime()) {
    nextStart.setDate(nextStart.getDate() + 7);
  }

  return nextStart;
};

const getEnrollmentWindow = (course: CourseCardItem, now = new Date()) => {
  const nextStart = getNextSessionStart(course);
  if (!nextStart) {
    return {
      canBook: true,
      message: "",
      countdown: "",
      opensAt: null as Date | null,
    };
  }

  const opensAt = new Date(nextStart.getTime() - 25 * 60 * 60 * 1000);

  if (now.getTime() >= nextStart.getTime()) {
    return {
      canBook: false,
      message: "Registration is closed for the current session.",
      countdown: "",
      opensAt,
    };
  }

  if (now.getTime() < opensAt.getTime()) {
    const diffMs = opensAt.getTime() - now.getTime();
    const hours = Math.floor(diffMs / (60 * 60 * 1000));
    const minutes = Math.floor((diffMs % (60 * 60 * 1000)) / (60 * 1000));

    if (hours > LONG_COUNTDOWN_THRESHOLD_HOURS) {
      return {
        canBook: false,
        message: "Enrollment is not open yet.",
        countdown: "",
        opensAt,
      };
    }

    return {
      canBook: false,
      message: `Enrollment opens in ${hours}h ${minutes}m.`,
      countdown: formatDuration(Math.max(0, diffMs / 60000)),
      opensAt,
    };
  }

  return {
    canBook: true,
    message: "",
    countdown: "",
    opensAt,
  };
};

type ClassFormState = {
  name: string;
  courseCode: string;
  description: string;
  startTime: string;
  endTime: string;
  capacity: string;
  duration: string;
  category: string;
  weekday: string;
};

const DEFAULT_CLASS_FORM: ClassFormState = {
  name: "",
  courseCode: "",
  description: "",
  startTime: "09:00",
  endTime: "10:00",
  capacity: "20",
  duration: "60",
  category: "",
  weekday: "Mon",
};

const WEEKDAY_OPTIONS = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"];

const mapClassToCard = (course: BackendClass): CourseCardItem => ({
  id: course.id,
  title: course.name,
  code: course.course_code,
  description: course.description,
  instructor: "TBD Coach",
  time: formatTime(course.start_time),
  startTimeRaw: toInputTime(course.start_time),
  endTimeRaw: toInputTime(course.end_time),
  day: normalizeWeekday(course.weekday),
  spots: course.spot,
  capacity: course.capacity,
  duration: course.duration,
  type: course.category || "General",
  image: getCourseEmoji(course.category),
});

const Browse = () => {
  const navigate = useNavigate();
  const { isAuthenticated, token, userId, role } = useAuthStore();
  const [activeView, setActiveView] = useState<"grid" | "calendar">("grid");
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedCourse, setSelectedCourse] = useState<CourseCardItem | null>(
    null,
  );
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [loading, setLoading] = useState(true);
  const [booking, setBooking] = useState(false);
  const [dropping, setDropping] = useState(false);
  const [courses, setCourses] = useState<CourseCardItem[]>([]);
  const [enrolledCourseIds, setEnrolledCourseIds] = useState<number[]>([]);
  const [courseToDrop, setCourseToDrop] = useState<CourseCardItem | null>(null);
  const [classModalOpen, setClassModalOpen] = useState(false);
  const [classModalMode, setClassModalMode] = useState<"create" | "edit">(
    "create",
  );
  const [classEditingTarget, setClassEditingTarget] =
    useState<CourseCardItem | null>(null);
  const [classSubmitting, setClassSubmitting] = useState(false);
  const [classFormError, setClassFormError] = useState<string | null>(null);
  const [classForm, setClassForm] =
    useState<ClassFormState>(DEFAULT_CLASS_FORM);
  const [classToDelete, setClassToDelete] = useState<CourseCardItem | null>(
    null,
  );
  const [classDeleting, setClassDeleting] = useState(false);
  const [rosterCourse, setRosterCourse] = useState<CourseCardItem | null>(null);
  const [rosterLoading, setRosterLoading] = useState(false);
  const [rosterItems, setRosterItems] = useState<ClassEnrollmentItem[]>([]);
  const [clockTick, setClockTick] = useState(Date.now());

  const canManageClasses = role === "manager" || role === "supermanager";

  const loadClasses = async () => {
    setLoading(true);
    try {
      const data = await listClassesRequest(1);
      setCourses((data.classes || []).map(mapClassToCard));
    } catch (error) {
      const message =
        error instanceof Error ? error.message : "Failed to load classes";
      toast.error(message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void loadClasses();
  }, []);

  const loadEnrollments = useCallback(async () => {
    if (!isAuthenticated || !token || !userId) {
      setEnrolledCourseIds([]);
      return;
    }

    try {
      const data = await getUserEnrollmentsRequest(token, userId);
      setEnrolledCourseIds((data.courses || []).map((course) => course.id));
    } catch {
      setEnrolledCourseIds([]);
    }
  }, [isAuthenticated, token, userId]);

  useEffect(() => {
    void loadEnrollments();
  }, [loadEnrollments]);

  useEffect(() => {
    const timer = window.setInterval(() => {
      setClockTick(Date.now());
    }, 30000);

    return () => window.clearInterval(timer);
  }, []);

  const handleCourseSelect = (course: CourseCardItem) => {
    setSelectedCourse(course);
    setIsModalOpen(true);
  };

  const handleBookCourse = async (course: CourseCardItem) => {
    if (booking) {
      return;
    }

    if (!isAuthenticated || !token) {
      toast.error("Please login first to register classes.");
      navigate("/login");
      return;
    }

    if (enrolledCourseIds.includes(course.id)) {
      toast("You are already enrolled in this class.", { icon: "ℹ️" });
      return;
    }

    const windowState = getEnrollmentWindow(course, new Date(clockTick));
    if (!windowState.canBook) {
      toast.error(windowState.message || "Enrollment is not open yet.");
      return;
    }

    const scrollY = window.scrollY;
    setBooking(true);
    try {
      await registerClassRequest(token, course.id);
      toast.success(
        course.spots === 0 ? "Added to waitlist." : "Booking successful!",
      );
      setIsModalOpen(false);
      await Promise.all([loadClasses(), loadEnrollments()]);
      window.scrollTo(0, scrollY);
    } catch (error) {
      const message =
        error instanceof Error ? error.message : "Failed to book class";
      toast.error(message);
    } finally {
      setBooking(false);
    }
  };

  const handleQuickActionOnCard = async (course: CourseCardItem) => {
    if (isEnrolledCourse(course.id)) {
      requestDropCourse(course);
      return;
    }

    await handleBookCourse(course);
  };

  const requestDropCourse = (course: CourseCardItem) => {
    setCourseToDrop(course);
  };

  const handleDropCourse = async () => {
    if (!courseToDrop || !token || dropping) {
      return;
    }

    const scrollY = window.scrollY;
    setDropping(true);
    try {
      await dropClassRequest(token, courseToDrop.id);
      toast.success("Class dropped successfully.");
      setCourseToDrop(null);
      setIsModalOpen(false);
      await Promise.all([loadClasses(), loadEnrollments()]);
      window.scrollTo(0, scrollY);
    } catch (error) {
      const message =
        error instanceof Error ? error.message : "Failed to drop class";
      toast.error(message);
    } finally {
      setDropping(false);
    }
  };

  const handleViewRoster = async (course: CourseCardItem) => {
    if (!token) {
      toast.error("Please login first.");
      navigate("/login");
      return;
    }

    setRosterCourse(course);
    setRosterLoading(true);
    setRosterItems([]);

    try {
      const data = await listClassEnrollmentsRequest(token, course.id);
      setRosterItems(data.enrollments || []);
    } catch (error) {
      const message =
        error instanceof Error ? error.message : "Failed to load roster";
      toast.error(message);
      setRosterItems([]);
    } finally {
      setRosterLoading(false);
    }
  };

  const filteredCourses = useMemo(() => {
    return courses.filter(
      (course) =>
        course.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
        (course.code || "").toLowerCase().includes(searchQuery.toLowerCase()) ||
        course.instructor.toLowerCase().includes(searchQuery.toLowerCase()) ||
        course.type.toLowerCase().includes(searchQuery.toLowerCase()),
    );
  }, [courses, searchQuery]);

  const isEnrolledCourse = (courseId: number) =>
    enrolledCourseIds.includes(courseId);

  const currentPreviewTime = new Date(clockTick);

  const openCreateClassModal = () => {
    if (!canManageClasses) {
      return;
    }
    setClassModalMode("create");
    setClassEditingTarget(null);
    setClassForm(DEFAULT_CLASS_FORM);
    setClassFormError(null);
    setClassModalOpen(true);
  };

  const openEditClassModal = (course: CourseCardItem) => {
    if (!canManageClasses) {
      return;
    }

    setClassModalMode("edit");
    setClassEditingTarget(course);
    setClassForm({
      name: course.title,
      courseCode: course.code || "",
      description: course.description || "",
      startTime: course.startTimeRaw || "09:00",
      endTime: course.endTimeRaw || "10:00",
      capacity: String(course.capacity),
      duration: String(course.duration ?? 60),
      category: course.type || "",
      weekday: WEEKDAY_OPTIONS.includes(course.day) ? course.day : "Mon",
    });
    setClassFormError(null);
    setClassModalOpen(true);
  };

  const validateClassForm = (): string | null => {
    if (!classForm.name.trim()) {
      return "Class name is required.";
    }
    if (!classForm.courseCode.trim()) {
      return "Course code is required.";
    }
    if (!classForm.startTime || !classForm.endTime) {
      return "Start and end time are required.";
    }

    const capacity = Number(classForm.capacity);
    if (!Number.isInteger(capacity) || capacity < 1) {
      return "Capacity must be an integer greater than 0.";
    }

    const duration = Number(classForm.duration);
    if (!Number.isInteger(duration) || duration < 0) {
      return "Duration must be a non-negative integer.";
    }

    if (!WEEKDAY_OPTIONS.includes(classForm.weekday)) {
      return "Weekday is invalid.";
    }

    if (classForm.endTime <= classForm.startTime) {
      return "End time must be later than start time.";
    }

    return null;
  };

  const handleStartTimeChange = (nextStartTime: string) => {
    setClassForm((prev) => {
      const startMinutes = toMinutes(nextStartTime);
      const endMinutes = toMinutes(prev.endTime);

      if (startMinutes === null || endMinutes === null) {
        return { ...prev, startTime: nextStartTime };
      }

      const nextDuration = Math.max(0, endMinutes - startMinutes);
      return {
        ...prev,
        startTime: nextStartTime,
        duration: String(nextDuration),
      };
    });
  };

  const handleEndTimeChange = (nextEndTime: string) => {
    setClassForm((prev) => {
      const startMinutes = toMinutes(prev.startTime);
      const endMinutes = toMinutes(nextEndTime);

      if (startMinutes === null || endMinutes === null) {
        return { ...prev, endTime: nextEndTime };
      }

      const nextDuration = Math.max(0, endMinutes - startMinutes);
      return {
        ...prev,
        endTime: nextEndTime,
        duration: String(nextDuration),
      };
    });
  };

  const handleDurationChange = (nextDurationText: string) => {
    setClassForm((prev) => {
      const startMinutes = toMinutes(prev.startTime);
      const nextDuration = Number(nextDurationText);

      if (
        startMinutes === null ||
        !Number.isFinite(nextDuration) ||
        nextDuration < 0
      ) {
        return { ...prev, duration: nextDurationText };
      }

      const nextEndTime = fromMinutes(startMinutes + nextDuration);
      return {
        ...prev,
        duration: String(Math.floor(nextDuration)),
        endTime: nextEndTime,
      };
    });
  };

  const toClassPayload = (): ClassUpsertRequest => ({
    name: classForm.name.trim(),
    course_code: classForm.courseCode.trim(),
    description: classForm.description.trim(),
    start_time: classForm.startTime,
    end_time: classForm.endTime,
    capacity: Number(classForm.capacity),
    duration: Number(classForm.duration),
    category: classForm.category.trim(),
    weekday: classForm.weekday,
  });

  const handleSaveClass = async () => {
    if (!canManageClasses) {
      return;
    }

    if (!token) {
      toast.error("Please login first.");
      navigate("/login");
      return;
    }

    const validationError = validateClassForm();
    if (validationError) {
      setClassFormError(validationError);
      return;
    }

    const scrollY = window.scrollY;
    setClassSubmitting(true);
    setClassFormError(null);
    try {
      const payload = toClassPayload();

      if (classModalMode === "create") {
        await createClassRequest(token, payload);
        toast.success("Class created successfully.");
      } else {
        if (!classEditingTarget) {
          throw new Error("No class selected for editing.");
        }
        await updateClassRequest(token, classEditingTarget.id, payload);
        toast.success("Class updated successfully.");
      }

      setClassModalOpen(false);
      setClassEditingTarget(null);
      await loadClasses();
      window.scrollTo(0, scrollY);
    } catch (error) {
      const message =
        error instanceof Error ? error.message : "Failed to save class";
      setClassFormError(message);
      toast.error(message);
    } finally {
      setClassSubmitting(false);
    }
  };

  const handleDeleteClass = async () => {
    if (!classToDelete || !token || classDeleting) {
      return;
    }

    const scrollY = window.scrollY;
    setClassDeleting(true);
    try {
      await deleteClassRequest(token, classToDelete.id);
      toast.success("Class deleted successfully.");
      setClassToDelete(null);
      if (selectedCourse?.id === classToDelete.id) {
        setSelectedCourse(null);
        setIsModalOpen(false);
      }
      await loadClasses();
      window.scrollTo(0, scrollY);
    } catch (error) {
      const message =
        error instanceof Error ? error.message : "Failed to delete class";
      toast.error(message);
    } finally {
      setClassDeleting(false);
    }
  };

  return (
    <div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-slate-100 pb-6">
        <div>
          <h2 className="text-2xl font-bold text-slate-900 tracking-tight">
            Browse Classes
          </h2>
          <p className="text-slate-500 text-sm">
            Book your spot in our premium sessions.
          </p>
        </div>

        <div className="flex flex-col sm:flex-row items-center gap-3 w-full md:w-auto">
          <div className="relative w-full sm:w-64">
            <Icons.Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
            <input
              type="text"
              placeholder="Search classes..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full pl-9 pr-4 py-2 text-sm rounded-full border border-slate-200 bg-white/50 focus:bg-white focus:ring-2 focus:ring-indigo-500/20 transition-all outline-hidden"
            />
          </div>

          <div className="flex bg-slate-100 p-1 rounded-full w-full sm:w-auto">
            <button
              onClick={() => setActiveView("grid")}
              className={`flex-1 sm:flex-none flex items-center justify-center gap-2 px-4 py-1.5 rounded-full text-xs font-bold transition-all ${
                activeView === "grid"
                  ? "bg-white text-indigo-600 shadow-sm"
                  : "text-slate-500 hover:text-slate-700"
              }`}
            >
              <Icons.Grid className="w-3.5 h-3.5" />
              Cards
            </button>
            <button
              onClick={() => setActiveView("calendar")}
              className={`flex-1 sm:flex-none flex items-center justify-center gap-2 px-4 py-1.5 rounded-full text-xs font-bold transition-all ${
                activeView === "calendar"
                  ? "bg-white text-indigo-600 shadow-sm"
                  : "text-slate-500 hover:text-slate-700"
              }`}
            >
              <Icons.Calendar className="w-3.5 h-3.5" />
              Calendar
            </button>
          </div>

          {canManageClasses && (
            <Button
              className="w-full sm:w-auto text-xs px-4 py-1.5 rounded-full shadow-sm"
              onClick={openCreateClassModal}
            >
              Add Class
            </Button>
          )}
        </div>
      </div>

      {loading ? (
        <Card className="p-8 text-center border-dashed border-2 bg-transparent text-slate-500">
          Loading classes...
        </Card>
      ) : activeView === "grid" ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredCourses.length > 0 ? (
            filteredCourses.map((course) => {
              const windowState = getEnrollmentWindow(
                course,
                currentPreviewTime,
              );
              return (
                <Card
                  key={course.id}
                  className="p-0 overflow-hidden hover:-translate-y-1 transition-transform duration-300 cursor-pointer group"
                  onClick={() => handleCourseSelect(course)}
                >
                  <div className="h-32 bg-linear-to-br from-indigo-100 to-purple-100 flex items-center justify-center text-6xl relative">
                    {course.image}
                    <div className="absolute top-4 right-4 bg-white/90 backdrop-blur px-3 py-1 rounded-full text-xs font-bold text-indigo-600 shadow-sm">
                      {course.type}
                    </div>
                  </div>
                  <div className="p-6">
                    <div className="flex items-center justify-between gap-2">
                      <h3 className="text-xl font-bold text-slate-800 group-hover:text-indigo-600 transition-colors truncate min-w-0">
                        {course.title}
                      </h3>
                      {isEnrolledCourse(course.id) && (
                        <div className="inline-flex items-center gap-1 px-2 py-1 rounded-full bg-emerald-50 text-emerald-700 text-xs font-semibold shrink-0">
                          <Icons.Check /> Enrolled
                        </div>
                      )}
                    </div>
                    {course.code && (
                      <p className="text-xs font-semibold uppercase tracking-wide text-indigo-500 mt-1">
                        {course.code}
                      </p>
                    )}
                    <div className="flex items-center gap-4 text-sm text-slate-500 mt-2 mb-6">
                      <div className="flex items-center gap-1">
                        <Icons.User /> {course.instructor}
                      </div>
                      <div className="flex items-center gap-1">
                        <Icons.Clock /> {course.day}, {course.time}
                      </div>
                    </div>
                    <div className="flex items-center justify-between">
                      {course.spots === 0 ? (
                        <span className="text-amber-600 text-sm font-semibold flex items-center gap-1">
                          <span className="w-2 h-2 rounded-full bg-amber-500"></span>{" "}
                          Waitlist Only
                        </span>
                      ) : (
                        <span
                          className={`text-sm font-semibold flex items-center gap-1 ${
                            course.spots < 5
                              ? "text-rose-500"
                              : "text-emerald-600"
                          }`}
                        >
                          <span
                            className={`w-2 h-2 rounded-full ${
                              course.spots < 5
                                ? "bg-rose-500"
                                : "bg-emerald-500"
                            }`}
                          ></span>
                          {course.spots} spots left
                        </span>
                      )}
                      <Button
                        variant={
                          isEnrolledCourse(course.id)
                            ? "danger"
                            : course.spots === 0 || !windowState.canBook
                              ? "outline"
                              : "primary"
                        }
                        className="text-sm px-6"
                        disabled={
                          !isEnrolledCourse(course.id) && !windowState.canBook
                        }
                        onClick={async (e) => {
                          e.stopPropagation();
                          await handleQuickActionOnCard(course);
                        }}
                      >
                        {isEnrolledCourse(course.id)
                          ? "Drop"
                          : !windowState.canBook
                            ? windowState.countdown
                              ? `Opens in ${windowState.countdown}`
                              : "Not open yet"
                            : course.spots === 0
                              ? "Join Waitlist"
                              : "Book"}
                      </Button>
                    </div>
                    {canManageClasses && (
                      <div className="flex items-center justify-start gap-2 mt-3 pt-3 border-t border-slate-100">
                        <Button
                          variant="ghost"
                          className="text-xs px-4 py-1.5 rounded-full bg-indigo-50! text-indigo-700! border border-indigo-200! hover:bg-indigo-100!"
                          onClick={(e) => {
                            e.stopPropagation();
                            openEditClassModal(course);
                          }}
                        >
                          Edit
                        </Button>
                        <Button
                          variant="danger"
                          className="text-xs px-4 py-1.5 rounded-full border border-red-200"
                          onClick={(e) => {
                            e.stopPropagation();
                            setClassToDelete(course);
                          }}
                        >
                          Delete
                        </Button>
                      </div>
                    )}
                  </div>
                </Card>
              );
            })
          ) : (
            <div className="col-span-full py-20 text-center">
              <div className="inline-flex items-center justify-center w-12 h-12 rounded-full bg-slate-100 text-slate-400 mb-4">
                <Icons.Search className="w-6 h-6" />
              </div>
              <p className="text-slate-500">
                No classes found matching "{searchQuery}"
              </p>
              <button
                onClick={() => setSearchQuery("")}
                className="text-indigo-600 font-semibold text-sm mt-2 hover:underline"
              >
                Clear search
              </button>
            </div>
          )}
        </div>
      ) : (
        <CalendarView
          courses={filteredCourses}
          onEventClick={handleCourseSelect}
          enrolledCourseIds={enrolledCourseIds}
        />
      )}

      <CourseDetailsModal
        course={selectedCourse}
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onBook={handleBookCourse}
        onDrop={requestDropCourse}
        onViewRoster={canManageClasses ? handleViewRoster : undefined}
        booking={booking}
        dropping={dropping}
        enrolled={selectedCourse ? isEnrolledCourse(selectedCourse.id) : false}
        bookingLocked={
          selectedCourse
            ? !isEnrolledCourse(selectedCourse.id) &&
              !getEnrollmentWindow(selectedCourse, currentPreviewTime).canBook
            : false
        }
        bookingHint={
          selectedCourse && !isEnrolledCourse(selectedCourse.id)
            ? getEnrollmentWindow(selectedCourse, currentPreviewTime).message ||
              undefined
            : undefined
        }
      />

      <Modal
        isOpen={classModalOpen}
        onClose={() => setClassModalOpen(false)}
        panelClassName="max-w-[56rem]"
      >
        <div className="p-6 sm:p-8">
          <h3 className="text-2xl font-bold text-slate-900 mb-1">
            {classModalMode === "create" ? "Add Class" : "Edit Class"}
          </h3>
          <p className="text-sm text-slate-500 mb-6">
            {classModalMode === "create"
              ? "Create a new class for members."
              : "Update class information and schedule."}
          </p>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
            <label className="block">
              <span className="text-xs font-semibold uppercase tracking-wide text-slate-500">
                Name
              </span>
              <input
                value={classForm.name}
                onChange={(e) =>
                  setClassForm((prev) => ({ ...prev, name: e.target.value }))
                }
                className="mt-1 w-full px-3 py-2.5 rounded-xl border border-slate-200 bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500/20 outline-hidden"
              />
            </label>

            <label className="block md:col-span-2 md:row-span-2">
              <span className="text-xs font-semibold uppercase tracking-wide text-slate-500">
                Description
              </span>
              <textarea
                value={classForm.description}
                onChange={(e) =>
                  setClassForm((prev) => ({
                    ...prev,
                    description: e.target.value,
                  }))
                }
                rows={5}
                className="mt-1 w-full h-33 px-3 py-2.5 rounded-xl border border-slate-200 bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500/20 outline-hidden resize-none"
              />
            </label>

            <label className="block">
              <span className="text-xs font-semibold uppercase tracking-wide text-slate-500">
                Course Code
              </span>
              <input
                value={classForm.courseCode}
                onChange={(e) =>
                  setClassForm((prev) => ({
                    ...prev,
                    courseCode: e.target.value,
                  }))
                }
                className="mt-1 w-full px-3 py-2.5 rounded-xl border border-slate-200 bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500/20 outline-hidden"
              />
            </label>

            <label className="block">
              <span className="text-xs font-semibold uppercase tracking-wide text-slate-500">
                Category
              </span>
              <input
                value={classForm.category}
                onChange={(e) =>
                  setClassForm((prev) => ({
                    ...prev,
                    category: e.target.value,
                  }))
                }
                className="mt-1 w-full px-3 py-2.5 rounded-xl border border-slate-200 bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500/20 outline-hidden"
              />
            </label>

            <label className="block">
              <span className="text-xs font-semibold uppercase tracking-wide text-slate-500">
                Capacity
              </span>
              <input
                type="number"
                min={1}
                value={classForm.capacity}
                onChange={(e) =>
                  setClassForm((prev) => ({
                    ...prev,
                    capacity: e.target.value,
                  }))
                }
                className="mt-1 w-full px-3 py-2.5 rounded-xl border border-slate-200 bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500/20 outline-hidden"
              />
            </label>

            <label className="block">
              <span className="text-xs font-semibold uppercase tracking-wide text-slate-500">
                Weekday
              </span>
              <select
                value={classForm.weekday}
                onChange={(e) =>
                  setClassForm((prev) => ({ ...prev, weekday: e.target.value }))
                }
                className="mt-1 w-full px-3 py-2.5 rounded-xl border border-slate-200 bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500/20 outline-hidden"
              >
                {WEEKDAY_OPTIONS.map((day) => (
                  <option key={day} value={day}>
                    {day}
                  </option>
                ))}
              </select>
            </label>

            <label className="block">
              <span className="text-xs font-semibold uppercase tracking-wide text-slate-500">
                Start Time
              </span>
              <input
                type="time"
                value={classForm.startTime}
                onChange={(e) => handleStartTimeChange(e.target.value)}
                className="mt-1 w-full px-3 py-2.5 rounded-xl border border-slate-200 bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500/20 outline-hidden"
              />
            </label>

            <label className="block">
              <span className="text-xs font-semibold uppercase tracking-wide text-slate-500">
                End Time
              </span>
              <input
                type="time"
                value={classForm.endTime}
                onChange={(e) => handleEndTimeChange(e.target.value)}
                className="mt-1 w-full px-3 py-2.5 rounded-xl border border-slate-200 bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500/20 outline-hidden"
              />
            </label>

            <label className="block">
              <span className="text-xs font-semibold uppercase tracking-wide text-slate-500">
                Duration (minutes)
              </span>
              <input
                type="number"
                min={0}
                value={classForm.duration}
                onChange={(e) => handleDurationChange(e.target.value)}
                className="mt-1 w-full px-3 py-2.5 rounded-xl border border-slate-200 bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500/20 outline-hidden"
              />
            </label>
          </div>

          {classFormError && (
            <p className="mt-4 text-sm font-medium text-rose-600">
              {classFormError}
            </p>
          )}

          <div className="grid grid-cols-2 gap-3 mt-6">
            <Button
              variant="secondary"
              onClick={() => {
                setClassModalOpen(false);
                setClassFormError(null);
              }}
            >
              Cancel
            </Button>
            <Button onClick={handleSaveClass}>
              {classSubmitting
                ? classModalMode === "create"
                  ? "Creating..."
                  : "Saving..."
                : classModalMode === "create"
                  ? "Create Class"
                  : "Save Changes"}
            </Button>
          </div>
        </div>
      </Modal>

      {courseToDrop && (
        <div className="fixed inset-0 z-120 flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-sm animate-in fade-in duration-200">
          <Card className="max-w-sm w-full p-6 text-center animate-in zoom-in-95 duration-200">
            <div className="w-12 h-12 bg-red-100 text-red-600 rounded-full flex items-center justify-center mx-auto mb-4">
              <Icons.Alert />
            </div>
            <h3 className="text-lg font-bold text-slate-900">
              Cancel Registration?
            </h3>
            <p className="text-slate-500 text-sm mt-2 mb-6">
              Are you sure you want to drop this class? You will lose your spot
              to the next person on the waitlist.
            </p>
            <div className="grid grid-cols-2 gap-3">
              <Button variant="secondary" onClick={() => setCourseToDrop(null)}>
                Keep Spot
              </Button>
              <Button variant="danger" onClick={handleDropCourse}>
                {dropping ? "Dropping..." : "Yes, Drop"}
              </Button>
            </div>
          </Card>
        </div>
      )}

      {classToDelete && (
        <div className="fixed inset-0 z-120 flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-sm animate-in fade-in duration-200">
          <Card className="max-w-sm w-full p-6 text-center animate-in zoom-in-95 duration-200">
            <div className="w-12 h-12 bg-red-100 text-red-600 rounded-full flex items-center justify-center mx-auto mb-4">
              <Icons.Alert />
            </div>
            <h3 className="text-lg font-bold text-slate-900">Delete Class?</h3>
            <p className="text-slate-500 text-sm mt-2 mb-6">
              This action will permanently delete "{classToDelete.title}" and
              cannot be undone.
            </p>
            <div className="grid grid-cols-2 gap-3">
              <Button
                variant="secondary"
                onClick={() => {
                  if (!classDeleting) {
                    setClassToDelete(null);
                  }
                }}
              >
                Cancel
              </Button>
              <Button variant="danger" onClick={handleDeleteClass}>
                {classDeleting ? "Deleting..." : "Yes, Delete"}
              </Button>
            </div>
          </Card>
        </div>
      )}

      {rosterCourse && (
        <Modal
          isOpen={Boolean(rosterCourse)}
          onClose={() => {
            setRosterCourse(null);
            setRosterItems([]);
          }}
          panelClassName="max-w-3xl"
        >
          <div className="p-6 sm:p-8">
            <div className="flex items-start justify-between gap-4 mb-6">
              <div>
                <h3 className="text-2xl font-bold text-slate-900">
                  Session Roster
                </h3>
                <p className="text-sm text-slate-500 mt-1">
                  {rosterCourse.title} • {rosterCourse.day}, {rosterCourse.time}
                </p>
              </div>
              <Badge className="bg-slate-900 text-white">
                {rosterItems.length} students
              </Badge>
            </div>

            {rosterLoading ? (
              <Card className="p-6 text-center text-slate-500 bg-transparent border-dashed border-2">
                Loading roster...
              </Card>
            ) : rosterItems.length > 0 ? (
              <div className="space-y-3 max-h-[32rem] overflow-y-auto pr-1">
                {rosterItems.map((item) => (
                  <Card
                    key={item.id}
                    className="p-4 flex items-center justify-between"
                  >
                    <div>
                      <p className="font-semibold text-slate-800">
                        {item.user?.name || `User ${item.user_id}`}
                      </p>
                      <p className="text-sm text-slate-500">
                        {item.user?.email || "No email available"}
                      </p>
                    </div>
                    <Badge
                      className={
                        item.status === "attended" || item.status === "present"
                          ? "bg-emerald-100 text-emerald-700"
                          : item.status === "missed" || item.status === "absent"
                            ? "bg-rose-100 text-rose-700"
                            : "bg-slate-100 text-slate-500"
                      }
                    >
                      {item.status}
                    </Badge>
                  </Card>
                ))}
              </div>
            ) : (
              <Card className="p-6 text-center text-slate-500 bg-transparent border-dashed border-2">
                No enrolled students were returned for this session.
              </Card>
            )}
          </div>
        </Modal>
      )}
    </div>
  );
};

export default Browse;
