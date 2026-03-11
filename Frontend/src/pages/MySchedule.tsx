import React, { useState, useMemo } from "react";
import { ENROLLED, HISTORY } from "../lib/constants";
import Card from "../components/ui/Card";
import Button from "../components/ui/Button";
import { Icons } from "../lib/icons";
import { cn } from "../lib/utils";
import { 
  BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, 
  PieChart, Pie, Cell 
} from 'recharts';
import { 
  startOfWeek, endOfWeek, startOfMonth, endOfMonth, startOfYear, 
  isWithinInterval, parseISO, format, eachDayOfInterval, eachMonthOfInterval 
} from 'date-fns';

type TimeRange = "week" | "month" | "year";

// Map theme variables to hex for Recharts
const BRAND_COLORS = [
    'var(--color-brand-indigo)', 
    'var(--color-brand-violet)', 
    'var(--color-brand-pink)', 
    'var(--color-brand-rose)', 
    'var(--color-brand-amber)'
];

const MySchedule = () => {
    const [showModal, setShowModal] = useState(false);
    const [activeTab, setActiveTab] = useState<"upcoming" | "history">("upcoming");
    const [timeRange, setTimeRange] = useState<TimeRange>("month");

    // Use current date based on mock data timeline (March 11, 2026)
    const referenceDate = new Date(2026, 2, 11);

    const filteredHistory = useMemo(() => {
        let start, end;
        if (timeRange === "week") {
            start = startOfWeek(referenceDate);
            end = endOfWeek(referenceDate);
        } else if (timeRange === "month") {
            start = startOfMonth(referenceDate);
            end = endOfMonth(referenceDate);
        } else {
            start = startOfYear(referenceDate);
            end = referenceDate;
        }

        return HISTORY.filter(item => 
            isWithinInterval(parseISO(item.date), { start, end })
        );
    }, [timeRange]);

    // Process data for charts
    const totalSessions = filteredHistory.length;
    const totalHours = filteredHistory.reduce((acc, curr) => acc + curr.duration, 0);
    
    const typeCounts = filteredHistory.reduce((acc: Record<string, number>, curr) => {
        acc[curr.type] = (acc[curr.type] || 0) + 1;
        return acc;
    }, {});
    
    const typeData = Object.keys(typeCounts).map(type => ({
        name: type,
        value: typeCounts[type]
    }));

    const frequencyData = useMemo(() => {
        if (timeRange === "year") {
            const months = eachMonthOfInterval({
                start: startOfYear(referenceDate),
                end: referenceDate
            });
            return months.map(m => {
                const monthStr = format(m, 'MMM');
                const count = filteredHistory.filter(item => 
                    format(parseISO(item.date), 'MMM') === monthStr
                ).length;
                return { label: monthStr, count };
            });
        } else {
            const intervalStart = timeRange === "week" ? startOfWeek(referenceDate) : startOfMonth(referenceDate);
            const days = eachDayOfInterval({
                start: intervalStart,
                end: referenceDate
            });
            return days.map(d => {
                const dateStr = format(d, 'yyyy-MM-dd');
                const count = filteredHistory.filter(item => item.date === dateStr).length;
                return { 
                    label: format(d, timeRange === 'week' ? 'EEE' : 'd'), 
                    count 
                };
            });
        }
    }, [filteredHistory, timeRange]);

  return (
    <div className="max-w-2xl mx-auto space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div className="flex flex-col gap-6">
        <div className="flex items-center justify-between">
            <h2 className="text-2xl font-bold text-slate-800">
                My Schedule
            </h2>
        </div>

        {/* Tab Navigation */}
        <div className="flex gap-1 p-1 bg-slate-100 rounded-xl w-fit">
            <button
                onClick={() => setActiveTab("upcoming")}
                className={cn(
                    "px-4 py-2 text-sm font-medium rounded-lg transition-all",
                    activeTab === "upcoming" 
                        ? "bg-white text-indigo-600 shadow-sm" 
                        : "text-slate-500 hover:text-slate-700"
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
                        : "text-slate-500 hover:text-slate-700"
                )}
            >
                History & Charts
            </button>
        </div>
      </div>

      {activeTab === "upcoming" ? (
        <div className="space-y-4">
            {ENROLLED.length > 0 ? (
                ENROLLED.map((course) => (
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
                        onClick={() => setShowModal(true)}
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
                    <p className="text-slate-500">No upcoming sessions. Time to book some!</p>
                </Card>
            )}
        </div>
      ) : (
        <div className="space-y-6">
            {/* Filter Dropdown */}
            <div className="flex justify-end">
                <div className="relative inline-block text-left">
                    <select 
                        value={timeRange}
                        onChange={(e) => setTimeRange(e.target.value as TimeRange)}
                        className="appearance-none bg-white border border-slate-200 rounded-lg px-4 py-2 pr-8 text-sm font-medium text-slate-700 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 cursor-pointer shadow-sm hover:border-slate-300 transition-colors"
                    >
                        <option value="week">This Week</option>
                        <option value="month">This Month</option>
                        <option value="year">This Year</option>
                    </select>
                    <div className="absolute inset-y-0 right-0 flex items-center px-2 pointer-events-none text-slate-400">
                        <Icons.ChevronDown />
                    </div>
                </div>
            </div>

            {/* Stats Overview */}
            <div className="grid grid-cols-2 gap-4">
                <Card className="p-4 flex items-center gap-4">
                    <div className="w-10 h-10 rounded-xl bg-indigo-50 flex items-center justify-center text-indigo-600">
                        <Icons.Check />
                    </div>
                    <div>
                        <p className="text-xs text-slate-500 font-medium">Sessions</p>
                        <p className="text-xl font-bold text-slate-800">{totalSessions}</p>
                    </div>
                </Card>
                <Card className="p-4 flex items-center gap-4">
                    <div className="w-10 h-10 rounded-xl bg-emerald-50 flex items-center justify-center text-emerald-600">
                        <Icons.Clock />
                    </div>
                    <div>
                        <p className="text-xs text-slate-500 font-medium">Hours</p>
                        <p className="text-xl font-bold text-slate-800">{totalHours}h</p>
                    </div>
                </Card>
            </div>

            {/* Charts Section */}
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
                                    contentStyle={{ borderRadius: '12px', border: 'none', boxShadow: '0 4px 6px -1px rgb(0 0 0 / 0.1)' }}
                                    cursor={{ fill: '#f1f5f9' }}
                                />
                                <Bar dataKey="count" fill="var(--color-brand-indigo)" radius={[4, 4, 0, 0]} />
                            </BarChart>
                        </ResponsiveContainer>
                    </div>
                </Card>
                <Card className="p-6">
                    <h4 className="text-sm font-bold text-slate-800 mb-4">Class Types</h4>
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
                                            <Cell key={`cell-${index}`} fill={BRAND_COLORS[index % BRAND_COLORS.length]} />
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

            {/* Detailed History List */}
            <div className="space-y-4">
                <h3 className="text-lg font-bold text-slate-800 px-1">Detailed History</h3>
                {filteredHistory.length > 0 ? (
                    [...filteredHistory].reverse().map((session) => (
                        <Card key={session.id} className="p-4 flex items-center justify-between hover:bg-slate-50/50 transition-colors">
                            <div className="flex items-center gap-4">
                                <div className="w-12 h-12 rounded-xl bg-slate-50 flex items-center justify-center text-xl">
                                    {session.type === 'Yoga' ? '🧘‍♀️' : session.type === 'Cardio' ? '🔥' : '💪'}
                                </div>
                                <div>
                                    <h4 className="font-bold text-slate-800">{session.title}</h4>
                                    <p className="text-xs text-slate-500">
                                        {format(parseISO(session.date), 'MMM d, yyyy')} • {session.instructor}
                                    </p>
                                </div>
                            </div>
                            <div className="text-right">
                                <span className="text-xs font-semibold bg-emerald-50 text-emerald-600 px-2 py-1 rounded-full">
                                    Completed
                                </span>
                                <p className="text-xs text-slate-400 mt-1">{session.duration}h</p>
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

       {showModal && (
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
              <Button variant="secondary" onClick={() => setShowModal(false)}>
                Keep Spot
              </Button>
              <Button variant="danger" onClick={() => setShowModal(false)}>
                Yes, Drop
              </Button>
            </div>
          </Card>
        </div>
      )}
    </div>
  );
};

export default MySchedule;
