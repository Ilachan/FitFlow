import { useState, useMemo, useEffect, useCallback } from "react";
import Card from "../components/ui/Card";
import Button from "../components/ui/Button";
import { Icons } from "../lib/icons";
import { cn } from "../lib/utils";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
} from "recharts";
import { format, parseISO } from "date-fns";
import {
  dropClassRequest,
  getUserAnalyticsRequest,
  getUserEnrollmentsRequest,
  type BackendClass,
  type UserAnalyticsResponse,
} from "../lib/api";
import { useAuthStore } from "../store/authStore";
import toast from "react-hot-toast";

type TimeRange = "7days" | "30days" | "90days";

type ScheduleCourse = {
  id: number;
  title: string;
  time: string;
};

const BRAND_COLORS = [
  "var(--color-brand-indigo)",
  "var(--color-brand-violet)",
  "var(--color-brand-pink)",
  "var(--color-brand-rose)",
  "var(--color-brand-amber)",
];

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

const mapClassToSchedule = (course: BackendClass): ScheduleCourse => ({
  id: course.id,
  title: course.name,
  time: `${normalizeWeekday(course.weekday)}, ${formatTime(course.start_time)}`,
});

const toBackendRange = (timeRange: TimeRange): "7d" | "1m" | "3m" => {
  if (timeRange === "30days") {
    return "1m";
  }
  if (timeRange === "90days") {
    return "3m";
  }
  return "7d";
};

const MySchedule = () => {
  const { token, userId, isAuthenticated } = useAuthStore();
  const [activeTab, setActiveTab] = useState<"upcoming" | "history">(
    "upcoming",
  );
  const [timeRange, setTimeRange] = useState<TimeRange>("30days");
  const [loading, setLoading] = useState(true);
  const [dropping, setDropping] = useState(false);
  const [upcomingCourses, setUpcomingCourses] = useState<ScheduleCourse[]>([]);
  const [analytics, setAnalytics] = useState<
    UserAnalyticsResponse["analytics"] | null
  >(null);
  const [courseToDrop, setCourseToDrop] = useState<ScheduleCourse | null>(null);

  const loadScheduleData = useCallback(async () => {
    if (!isAuthenticated || !token || !userId) {
      setUpcomingCourses([]);
      setAnalytics(null);
      setLoading(false);
      return;
    }

    setLoading(true);
    try {
      const [enrollmentData, analyticsData] = await Promise.all([
        getUserEnrollmentsRequest(token, userId),
        getUserAnalyticsRequest(token, userId, toBackendRange(timeRange)),
      ]);

      setUpcomingCourses(
        (enrollmentData.courses || []).map(mapClassToSchedule),
      );
      setAnalytics(analyticsData.analytics);
    } catch (error) {
      const message =
        error instanceof Error ? error.message : "Failed to load schedule";
      toast.error(message);
    } finally {
      setLoading(false);
    }
  }, [isAuthenticated, token, userId, timeRange]);

  useEffect(() => {
    void loadScheduleData();
  }, [loadScheduleData]);

  const totalSessions = analytics?.total_classes || 0;
  const totalHours = totalSessions;

  const typeData = useMemo(() => {
    return (analytics?.categories || []).map((item) => ({
      name: item.category || "General",
      value: item.classes,
    }));
  }, [analytics]);

  const frequencyData = useMemo(() => {
    const daily = analytics?.daily || [];

    if (timeRange === "90days") {
      const monthTotals: Record<string, number> = {};
      daily.forEach((item) => {
        const label = format(parseISO(item.date), "MMM");
        monthTotals[label] = (monthTotals[label] || 0) + item.classes;
      });

      return Object.entries(monthTotals).map(([label, count]) => ({
        label,
        count,
      }));
    }

    return daily.map((item) => ({
      label: format(parseISO(item.date), timeRange === "week" ? "EEE" : "d"),
      count: item.classes,
      date: item.date,
    }));
  }, [analytics, timeRange]);

  const detailRows = useMemo(() => {
    return [...(analytics?.daily || [])]
      .sort((a, b) => (a.date > b.date ? -1 : 1))
      .map((item, index) => ({
        id: `${item.date}-${index}`,
        title:
          item.classes > 1
            ? `${item.classes} classes completed`
            : "1 class completed",
        date: item.date,
      }));
  }, [analytics]);

  const handleDropCourse = async () => {
    if (!courseToDrop || !token || dropping) {
      return;
    }

    setDropping(true);
    try {
      await dropClassRequest(token, courseToDrop.id);
      toast.success("Class dropped successfully.");
      setCourseToDrop(null);
      await loadScheduleData();
    } catch (error) {
      const message =
        error instanceof Error ? error.message : "Failed to drop class";
      toast.error(message);
    } finally {
      setDropping(false);
    }
  };

  return (
    <div className="max-w-2xl mx-auto space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div className="flex flex-col gap-6">
        <div className="flex items-center justify-between">
          <h2 className="text-2xl font-bold text-slate-800">My Schedule</h2>
        </div>

        <div className="flex gap-1 p-1 bg-slate-100 rounded-xl w-fit">
          <button
            onClick={() => setActiveTab("upcoming")}
            className={cn(
              "px-4 py-2 text-sm font-medium rounded-lg transition-all",
              activeTab === "upcoming"
                ? "bg-white text-indigo-600 shadow-sm"
                : "text-slate-500 hover:text-slate-700",
            )}
          >
            Upcoming Sessions
          </button>
          <button
            onClick={() => setActiveTab("history")}
            className={cn(
              "px-4 py-2 text-sm font-medium rounded-lg transition-all",
              activeTab === "history"
                ? "bg-white text-indigo-600 shadow-sm"
                : "text-slate-500 hover:text-slate-700",
            )}
          >
            History & Charts
          </button>
        </div>
      </div>

      {loading ? (
        <Card className="p-8 text-center border-dashed border-2 bg-transparent text-slate-500">
          Loading schedule...
        </Card>
      ) : activeTab === "upcoming" ? (
        <div className="space-y-4">
          {upcomingCourses.length > 0 ? (
            upcomingCourses.map((course) => (
              <Card
                key={course.id}
                className="p-6 flex items-center justify-between group"
              >
                <div className="flex items-center gap-4">
                  <div className="w-16 h-16 rounded-2xl bg-indigo-50 flex items-center justify-center text-2xl text-indigo-600">
                    <Icons.Calendar />
                  </div>
                  <div>
                    <h3 className="text-lg font-bold text-slate-800">
                      {course.title}
                    </h3>
                    <p className="text-slate-500 text-sm flex items-center gap-2 mt-1">
                      <span className="bg-indigo-50 text-indigo-700 px-2 py-0.5 rounded text-xs font-semibold">
                        Confirmed
                      </span>
                      • {course.time}
                    </p>
                  </div>
                </div>
                <Button
                  variant="danger"
                  onClick={() => setCourseToDrop(course)}
                  className="opacity-0 group-hover:opacity-100 transition-opacity"
                >
                  Drop
                </Button>
              </Card>
            ))
          ) : (
            <Card className="p-12 text-center border-dashed border-2 bg-transparent">
              <div className="w-16 h-16 bg-slate-50 text-slate-300 rounded-full flex items-center justify-center mx-auto mb-4">
                <Icons.Calendar />
              </div>
              <p className="text-slate-500">
                No upcoming sessions. Time to book some!
              </p>
            </Card>
          )}
        </div>
      ) : (
        <div className="space-y-6">
          <div className="flex justify-end">
            <div className="relative inline-block text-left">
              <select
                value={timeRange}
                onChange={(e) => setTimeRange(e.target.value as TimeRange)}
                className="appearance-none bg-white border border-slate-200 rounded-lg px-4 py-2 pr-8 text-sm font-medium text-slate-700 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 cursor-pointer shadow-sm hover:border-slate-300 transition-colors"
              >
                <option value="7days">Last 7 Days</option>
                <option value="30days">Last 30 Days</option>
                <option value="90days">Last 90 Days</option>
              </select>
              <div className="absolute inset-y-0 right-0 flex items-center px-2 pointer-events-none text-slate-400">
                <Icons.ChevronDown />
              </div>
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <Card className="p-4 flex items-center gap-4">
              <div className="w-10 h-10 rounded-xl bg-indigo-50 flex items-center justify-center text-indigo-600">
                <Icons.Check />
              </div>
              <div>
                <p className="text-xs text-slate-500 font-medium">Sessions</p>
                <p className="text-xl font-bold text-slate-800">
                  {totalSessions}
                </p>
              </div>
            </Card>
            <Card className="p-4 flex items-center gap-4">
              <div className="w-10 h-10 rounded-xl bg-emerald-50 flex items-center justify-center text-emerald-600">
                <Icons.Clock />
              </div>
              <div>
                <p className="text-xs text-slate-500 font-medium">Hours</p>
                <p className="text-xl font-bold text-slate-800">
                  {totalHours}h
                </p>
              </div>
            </Card>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Card className="p-6">
              <h4 className="text-sm font-bold text-slate-800 mb-4">
                Activity Frequency ({timeRange})
              </h4>
              <div className="h-48 w-full">
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart data={frequencyData}>
                    <XAxis
                      dataKey="label"
                      fontSize={10}
                      axisLine={false}
                      tickLine={false}
                    />
                    <YAxis hide />
                    <Tooltip
                      contentStyle={{
                        borderRadius: "12px",
                        border: "none",
                        boxShadow: "0 4px 6px -1px rgb(0 0 0 / 0.1)",
                      }}
                      cursor={{ fill: "#f1f5f9" }}
                    />
                    <Bar
                      dataKey="count"
                      fill="var(--color-brand-indigo)"
                      radius={[4, 4, 0, 0]}
                    />
                  </BarChart>
                </ResponsiveContainer>
              </div>
            </Card>
            <Card className="p-6">
              <h4 className="text-sm font-bold text-slate-800 mb-4">
                Class Types
              </h4>
              <div className="h-48 w-full">
                {typeData.length > 0 ? (
                  <ResponsiveContainer width="100%" height="100%">
                    <PieChart>
                      <Pie
                        data={typeData}
                        innerRadius={50}
                        outerRadius={70}
                        paddingAngle={5}
                        dataKey="value"
                      >
                        {typeData.map((_, index) => (
                          <Cell
                            key={`cell-${index}`}
                            fill={BRAND_COLORS[index % BRAND_COLORS.length]}
                          />
                        ))}
                      </Pie>
                      <Tooltip />
                    </PieChart>
                  </ResponsiveContainer>
                ) : (
                  <div className="h-full flex items-center justify-center text-slate-400 text-xs text-center p-4">
                    No sessions completed yet for this period.
                  </div>
                )}
              </div>
            </Card>
          </div>

          <div className="space-y-4">
            <h3 className="text-lg font-bold text-slate-800 px-1">
              Detailed History
            </h3>
            {detailRows.length > 0 ? (
              detailRows.map((session) => (
                <Card
                  key={session.id}
                  className="p-4 flex items-center justify-between hover:bg-slate-50/50 transition-colors"
                >
                  <div className="flex items-center gap-4">
                    <div className="w-12 h-12 rounded-xl bg-slate-50 flex items-center justify-center text-xl">
                      ✅
                    </div>
                    <div>
                      <h4 className="font-bold text-slate-800">
                        {session.title}
                      </h4>
                      <p className="text-xs text-slate-500">
                        {format(parseISO(session.date), "MMM d, yyyy")}
                      </p>
                    </div>
                  </div>
                  <div className="text-right">
                    <span className="text-xs font-semibold bg-emerald-50 text-emerald-600 px-2 py-1 rounded-full">
                      Completed
                    </span>
                  </div>
                </Card>
              ))
            ) : (
              <Card className="p-8 text-center text-slate-500 border-dashed border-2 bg-transparent">
                No sessions found for this period.
              </Card>
            )}
          </div>
        </div>
      )}

      {courseToDrop && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-sm animate-in fade-in duration-200">
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
    </div>
  );
};

export default MySchedule;
