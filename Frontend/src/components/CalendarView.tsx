import { useMemo, useEffect } from "react";
import { useCalendarApp, ScheduleXCalendar } from "@schedule-x/react";
import {
  createViewMonthGrid,
  createViewWeek,
  createViewDay,
} from "@schedule-x/calendar";
import { Temporal } from "temporal-polyfill";
import "@schedule-x/theme-default/dist/index.css";
import { COURSES } from "../lib/constants";
import { Icons } from "../lib/icons";

const DAY_MAP: Record<string, string> = {
  Mon: "2026-03-09",
  Tue: "2026-03-10",
  Wed: "2026-03-11",
  Thu: "2026-03-12",
  Fri: "2026-03-13",
  Sat: "2026-03-14",
  Sun: "2026-03-15",
};

const COLOR_STYLES = [
  { bg: "#f5f3ff", border: "#6366f1", text: "#4338ca" }, // Indigo
  { bg: "#fff1f2", border: "#f43f5e", text: "#9f1239" }, // Rose
  { bg: "#ecfdf5", border: "#10b981", text: "#065f46" }, // Emerald
  { bg: "#fffbeb", border: "#f59e0b", text: "#92400e" }, // Amber
  { bg: "#f0f9ff", border: "#0ea5e9", text: "#075985" }, // Sky
];

interface CalendarViewProps {
  courses?: typeof COURSES;
}

// Stable Custom Event Component
const CustomEvent = (props: any) => {
  const event = props.calendarEvent;
  const spotsColor = event.spots === 0 ? "#ef4444" : event.spots < 5 ? "#f59e0b" : "#10b981";
  const colors = event._style || COLOR_STYLES[0];

  return (
    <div 
      className={`sx__event-card-wrapper ${event._isPast ? 'sx__event--past' : ''}`}
      style={{ 
        backgroundColor: colors.bg, 
        borderLeft: `4px solid ${colors.border}`,
        height: '100%',
        width: '100%',
        display: 'flex',
        flexDirection: 'column'
      }}
    >
      <div className="sx__event-card-inner">
        <div className="sx__event-card-header">
          <span className="sx__event-card-emoji">{event.image}</span>
          <span className="sx__event-card-title" style={{ color: colors.text }}>{event.title}</span>
        </div>

        <div className="sx__event-card-instructor">
          <Icons.User className="w-3.5 h-3.5" />
          <span className="truncate">{event.instructor}</span>
        </div>

        <div className="sx__event-card-footer">
          <div className="sx__event-status-badge" style={{ backgroundColor: `${spotsColor}20`, color: spotsColor }}>
            <div className="sx__event-status-dot" style={{ backgroundColor: spotsColor }} />
            <span>{event.spots === 0 ? 'Full' : `${event.spots} spots left`}</span>
          </div>
        </div>
      </div>
    </div>
  );
};

const CalendarView = ({ courses = COURSES }: CalendarViewProps) => {
  const now = Temporal.Now.zonedDateTimeISO("UTC");

  const events = useMemo(() => {
    return courses.map((course, index) => {
      const dateStr = DAY_MAP[course.day.trim()];
      if (!dateStr) return null;

      try {
        const [timePart, period] = course.time.split(" ");
        let [hours, minutes] = timePart.split(":").map(Number);
        if (period === "PM" && hours !== 12) hours += 12;
        if (period === "AM" && hours === 12) hours = 0;

        const startStr = `${dateStr}T${String(hours).padStart(2, "0")}:${String(minutes).padStart(2, "0")}:00[UTC]`;
        const startZdt = Temporal.ZonedDateTime.from(startStr);
        const endZdt = startZdt.add({ hours: 1 });
        const isPast = Temporal.ZonedDateTime.compare(startZdt, now) < 0;

        return {
          id: String(course.id),
          title: course.title,
          instructor: course.instructor,
          spots: course.spots,
          image: course.image,
          start: startZdt,
          end: endZdt,
          calendarId: course.type.toLowerCase(),
          _isPast: isPast,
          _style: COLOR_STYLES[index % COLOR_STYLES.length]
        };
      } catch (e) {
        return null;
      }
    }).filter(event => event !== null);
  }, [courses, now]);

  const calendarApp = useCalendarApp({
    views: [createViewWeek(), createViewMonthGrid(), createViewDay()],
    events: events as any[],
    defaultView: "week",
    dayBoundaries: { start: "07:00", end: "22:00" },
  }, []);

  useEffect(() => {
    if (calendarApp && events) calendarApp.events.set(events as any[]);
  }, [events, calendarApp]);

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
        
        /* HIDE DEFAULT CONTENT */
        .sx__event-time, 
        .sx__event-title,
        .sx__time-grid-event-inner > *:not(.sx__event-card-wrapper) { 
            display: none !important; 
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

        /* Hover Effects */
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
          monthGridEvent: CustomEvent
        }}
      />
    </div>
  );
};

export default CalendarView;
