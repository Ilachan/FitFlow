import { useMemo } from "react";
import { useCalendarApp, ScheduleXCalendar } from "@schedule-x/react";
import {
  createViewMonthGrid,
  createViewWeek,
  createViewDay,
  type CalendarEventExternal,
} from "@schedule-x/calendar";
import { Temporal } from "temporal-polyfill";
import "@schedule-x/theme-default/dist/index.css";
import { Icons } from "../lib/icons";
import type { CourseCardItem } from "./CourseDetailsModal";

const COLOR_STYLES = [
  { bg: "#f5f3ff", border: "#6366f1", text: "#4338ca" }, // Indigo
  { bg: "#fff1f2", border: "#f43f5e", text: "#9f1239" }, // Rose
  { bg: "#ecfdf5", border: "#10b981", text: "#065f46" }, // Emerald
  { bg: "#fffbeb", border: "#f59e0b", text: "#92400e" }, // Amber
  { bg: "#f0f9ff", border: "#0ea5e9", text: "#075985" }, // Sky
];

interface CalendarViewProps {
  courses?: CourseCardItem[];
  onEventClick?: (course: CourseCardItem) => void;
  enrolledCourseIds?: number[];
}

type CalendarEventData = {
  id: string;
  title: string;
  instructor: string;
  spots: number;
  image: string;
  time: string;
  _isPast?: boolean;
  _enrolled?: boolean;
  _style?: { bg: string; border: string; text: string };
};

type CalendarEventProps = {
  calendarEvent: CalendarEventData;
};

const CustomEvent = (props: CalendarEventProps) => {
  const event = props.calendarEvent;
  const spotsColor = event.spots === 0 ? "#ef4444" : event.spots < 5 ? "#f59e0b" : "#10b981";
  const colors = event._style || COLOR_STYLES[0];

  return (
    <div
      className={`sx__event-card-wrapper ${event._isPast ? "sx__event--past" : ""}`}
      style={{
        backgroundColor: colors.bg,
        borderLeft: `4px solid ${colors.border}`,
        height: "100%",
        width: "100%",
        display: "flex",
        flexDirection: "column",
      }}
    >
      <div className="sx__event-card-inner">
        <div className="sx__event-card-header">
          <span className="sx__event-card-emoji">{event.image}</span>
          <span className="sx__event-card-title" style={{ color: colors.text }}>
            {event.title}
          </span>
          {event._enrolled && <span className="sx__event-enrolled-dot" />}
        </div>

        <div className="sx__event-card-instructor">
          <Icons.User className="w-3.5 h-3.5" />
          <span className="truncate">{event.instructor}</span>
        </div>

        <div className="sx__event-card-footer">
          <div
            className="sx__event-status-badge"
            style={{ backgroundColor: `${spotsColor}20`, color: spotsColor }}
          >
            <div
              className="sx__event-status-dot"
              style={{ backgroundColor: spotsColor }}
            />
            <span>
              {event.spots === 0 ? "Full" : `${event.spots} spots left`}
            </span>
          </div>
        </div>
      </div>
    </div>
  );
};

const MonthEvent = (props: CalendarEventProps) => {
  const event = props.calendarEvent;
  const colors = event._style || COLOR_STYLES[0];

  return (
    <div
      className="sx__month-event-wrapper"
      style={{
        backgroundColor: colors.bg,
        borderLeft: `2px solid ${colors.border}`,
        padding: "2px 6px",
        borderRadius: "4px",
        display: "flex",
        alignItems: "center",
        gap: "4px",
        fontSize: "0.7rem",
        fontWeight: 600,
        color: colors.text,
        overflow: "hidden",
        whiteSpace: "nowrap",
        width: "100%",
      }}
    >
      <span className="shrink-0">{event.image}</span>
      <span
        className="shrink-0 opacity-70 font-bold"
        style={{ fontSize: "0.6rem" }}
      >
        {event.time.split(" ")[0]}
      </span>
      <span className="truncate">{event.title}</span>
      {event._enrolled && <span className="sx__month-enrolled-dot" />}
    </div>
  );
};

const CalendarView = ({
  courses = [],
  onEventClick,
  enrolledCourseIds = [],
}: CalendarViewProps) => {
  const now = useMemo(() => Temporal.Now.zonedDateTimeISO("UTC"), []);
  const today = useMemo(() => now.toPlainDate(), [now]);

  // Generate dynamic day map for the current week
  const dayMap = useMemo(() => {
    const map: Record<string, string> = {};
    const days = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"];

    // Find the Monday of the current week
    // dayOfWeek: 1 (Mon) to 7 (Sun)
    const currentDayOfWeek = today.dayOfWeek;
    const diffToMonday = currentDayOfWeek - 1;
    const monday = today.subtract({ days: diffToMonday });

    days.forEach((day, index) => {
      map[day] = monday.add({ days: index }).toString();
    });

    return map;
  }, [today]);

  const events = useMemo(() => {
    return courses
      .map((course, index) => {
        const dateStr = dayMap[course.day.trim()];
        if (!dateStr) return null;

        try {
          const [timePart, period] = course.time.split(" ");
          const [parsedHours, minutes] = timePart.split(":").map(Number);
          let hours = parsedHours;
          if (period === "PM" && hours !== 12) hours += 12;
          if (period === "AM" && hours === 12) hours = 0;

          const startTimeStr = `${String(hours).padStart(2, "0")}:${String(minutes).padStart(2, "0")}:00`;
          const start = Temporal.ZonedDateTime.from(
            `${dateStr}T${startTimeStr}[UTC]`,
          );
          const end = start.add({ hours: 1 });
          const isPast = Temporal.ZonedDateTime.compare(start, now) < 0;

          return {
            id: String(course.id),
            title: course.title,
            instructor: course.instructor,
            spots: course.spots,
            image: course.image,
            time: course.time,
            start,
            end,
            calendarId: course.type.toLowerCase(),
            _isPast: isPast,
            _style: COLOR_STYLES[index % COLOR_STYLES.length],
            _enrolled: enrolledCourseIds.includes(course.id),
          };
        } catch {
          return null;
        }
      })
      .filter((event) => event !== null);
  }, [courses, dayMap, now, enrolledCourseIds]);

  const calendarApp = useCalendarApp(
    {
      views: [createViewWeek(), createViewMonthGrid(), createViewDay()],
      events: events as CalendarEventExternal[],
      defaultView: createViewWeek().name,
      selectedDate: today,
      dayBoundaries: {
        start: "05:00",
        end: "23:00",
      },
      callbacks: {
        onEventClick(calendarEvent: CalendarEventExternal, e: UIEvent) {
          void e;
          const course = courses.find(
            (c) => String(c.id) === String(calendarEvent.id),
          );
          if (course && onEventClick) {
            onEventClick(course);
          }
        },
      },
    },
    [events, onEventClick, courses, today],
  );

  return (
    <div className="bg-white rounded-2xl border border-slate-200 overflow-hidden sx-react-calendar-wrapper">
      <style>{`
        .sx-react-calendar-wrapper { height: 750px; }
        .sx__calendar { 
          --sx-color-primary: #6366f1; 
          --sx-border-radius-large: 1rem; 
          font-family: inherit; 
          border: none !important;
        }
        
        .sx__event-time, 
        .sx__event-title,
        .sx__time-grid-event-inner > *:not(.sx__event-card-wrapper),
        .sx__month-grid-event-inner > *:not(.sx__month-event-wrapper) { 
            display: none !important; 
        }

        .sx__month-grid-event {
          background-color: transparent !important;
          border: none !important;
          padding: 2px 4px !important;
        }

        .sx__event {
          cursor: pointer !important; 
          border-radius: 8px !important; 
          padding: 0 !important;
          border: none !important; 
          overflow: hidden !important;
          box-shadow: 0 1px 2px rgba(0,0,0,0.05) !important;
          background: transparent !important;
        }

        .sx__event-card-inner {
          display: flex;
          flex-direction: column;
          padding: 10px 12px;
          height: 100%;
          width: 100%;
          overflow: hidden;
          gap: 6px;
        }

        .sx__event--past {
          opacity: 0.6;
          filter: grayscale(0.4);
        }

        .sx__event-card-header {
          display: flex;
          justify-content: flex-start;
          align-items: center;
          gap: 8px;
        }

        .sx__event-enrolled-dot {
          width: 7px;
          height: 7px;
          border-radius: 9999px;
          background: #16a34a;
          flex-shrink: 0;
          margin-left: auto;
        }

        .sx__event-card-title {
          font-weight: 700 !important;
          font-size: 0.8rem !important;
          line-height: 1.2;
          display: -webkit-box;
          -webkit-line-clamp: 2;
          -webkit-box-orient: vertical;
          overflow: hidden;
        }

        .sx__event-card-emoji {
          font-size: 1.1rem;
          flex-shrink: 0;
        }

        .sx__event-card-instructor {
          display: flex;
          align-items: center;
          gap: 4px;
          font-size: 0.7rem !important;
          font-weight: 500 !important;
          color: #64748b;
        }

        .sx__event-card-footer {
          margin-top: auto;
          display: flex;
          align-items: center;
          justify-content: space-between;
          gap: 6px;
        }

        .sx__month-enrolled-dot {
          width: 6px;
          height: 6px;
          border-radius: 9999px;
          background: #16a34a;
          margin-left: auto;
          flex-shrink: 0;
        }

        .sx__event-status-badge {
          display: flex;
          align-items: center;
          gap: 4px;
          padding: 2px 6px;
          border-radius: 9999px;
          font-size: 0.65rem !important;
          font-weight: 600 !important;
        }

        .sx__event-status-dot {
          width: 5px;
          height: 5px;
          border-radius: 50%;
        }

        .sx__event:hover {
          transform: translateY(-1px);
          box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1) !important;
          z-index: 50 !important;
        }
        
        .sx__event:hover .sx__event-card-wrapper {
          filter: brightness(0.98);
        }
      `}</style>
      <ScheduleXCalendar
        calendarApp={calendarApp}
        customComponents={{
          timeGridEvent: CustomEvent,
          monthGridEvent: MonthEvent,
        }}
      />
    </div>
  );
};

export default CalendarView;
