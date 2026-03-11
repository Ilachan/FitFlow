import { useState, useMemo } from "react";
import { COURSES } from "../lib/constants";
import Card from "../components/ui/Card";
import Button from "../components/ui/Button";
import { Icons } from "../lib/icons";
import CalendarView from "../components/CalendarView";

const Browse = () => {
  const [activeView, setActiveView] = useState<"grid" | "calendar">("grid");
  const [searchQuery, setSearchQuery] = useState("");

  const filteredCourses = useMemo(() => {
    return COURSES.filter(
      (course) =>
        course.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
        course.instructor.toLowerCase().includes(searchQuery.toLowerCase()) ||
        course.type.toLowerCase().includes(searchQuery.toLowerCase())
    );
  }, [searchQuery]);

  return (
    <div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-slate-100 pb-6">
        <div>
          <h2 className="text-2xl font-bold text-slate-900 tracking-tight">Browse Classes</h2>
          <p className="text-slate-500 text-sm">Book your spot in our premium sessions.</p>
        </div>
        
        <div className="flex flex-col sm:flex-row items-center gap-3 w-full md:w-auto">
          {/* Search Bar */}
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

          {/* View Switcher */}
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
        </div>
      </div>

      {activeView === "grid" ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredCourses.length > 0 ? (
            filteredCourses.map((course) => (
              <Card
                key={course.id}
                className="p-0 overflow-hidden hover:-translate-y-1 transition-transform duration-300"
              >
                <div className="h-32 bg-linear-to-br from-indigo-100 to-purple-100 flex items-center justify-center text-6xl relative">
                  {course.image}
                  <div className="absolute top-4 right-4 bg-white/90 backdrop-blur px-3 py-1 rounded-full text-xs font-bold text-indigo-600 shadow-sm">
                    {course.type}
                  </div>
                </div>
                <div className="p-6">
                  <h3 className="text-xl font-bold text-slate-800">
                    {course.title}
                  </h3>
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
                          course.spots < 5 ? "text-rose-500" : "text-emerald-600"
                        }`}
                      >
                        <span
                          className={`w-2 h-2 rounded-full ${
                            course.spots < 5 ? "bg-rose-500" : "bg-emerald-500"
                          }`}
                        ></span>
                        {course.spots} spots left
                      </span>
                    )}
                    <Button
                      variant={course.spots === 0 ? "outline" : "primary"}
                      className="text-sm px-6"
                    >
                      {course.spots === 0 ? "Join Waitlist" : "Book"}
                    </Button>
                  </div>
                </div>
              </Card>
            ))
          ) : (
            <div className="col-span-full py-20 text-center">
              <div className="inline-flex items-center justify-center w-12 h-12 rounded-full bg-slate-100 text-slate-400 mb-4">
                <Icons.Search className="w-6 h-6" />
              </div>
              <p className="text-slate-500">No classes found matching "{searchQuery}"</p>
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
        <CalendarView courses={filteredCourses} />
      )}
    </div>
  );
};

export default Browse;
