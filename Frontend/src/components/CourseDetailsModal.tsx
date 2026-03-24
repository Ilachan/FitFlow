import Modal from "./ui/Modal";
import Button from "./ui/Button";
import { Icons } from "../lib/icons";

export type CourseCardItem = {
  id: number;
  title: string;
  instructor: string;
  time: string;
  day: string;
  spots: number;
  capacity: number;
  type: string;
  image: string;
};

interface CourseDetailsModalProps {
  course: CourseCardItem | null;
  isOpen: boolean;
  onClose: () => void;
  onBook?: (course: CourseCardItem) => void;
  onDrop?: (course: CourseCardItem) => void;
  booking?: boolean;
  dropping?: boolean;
  enrolled?: boolean;
}

const CourseDetailsModal = ({
  course,
  isOpen,
  onClose,
  onBook,
  onDrop,
  booking = false,
  dropping = false,
  enrolled = false,
}: CourseDetailsModalProps) => {
  if (!course) return null;

  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <div className="h-48 bg-linear-to-br from-indigo-100 to-purple-100 flex items-center justify-center text-8xl relative">
        {course.image}
        <div className="absolute top-4 left-4 bg-white/90 backdrop-blur px-3 py-1 rounded-full text-xs font-bold text-indigo-600 shadow-sm">
          {course.type}
        </div>
      </div>
      <div className="p-8">
        <div className="flex items-center justify-between gap-2 mb-2">
          <h3 className="text-2xl font-bold text-slate-900 truncate min-w-0">
            {course.title}
          </h3>
          {enrolled && (
            <div className="inline-flex items-center gap-1 px-2 py-1 rounded-full bg-emerald-50 text-emerald-700 text-xs font-semibold shrink-0">
              <Icons.Check /> Enrolled
            </div>
          )}
        </div>

        <div className="space-y-3 mb-8">
          <div className="flex items-center gap-3 text-slate-600">
            <div className="w-8 h-8 rounded-full bg-slate-100 flex items-center justify-center text-slate-500">
              <Icons.User className="w-4 h-4" />
            </div>
            <div>
              <p className="text-xs text-slate-400 font-medium uppercase tracking-wider">
                Instructor
              </p>
              <p className="font-semibold">{course.instructor}</p>
            </div>
          </div>

          <div className="flex items-center gap-3 text-slate-600">
            <div className="w-8 h-8 rounded-full bg-slate-100 flex items-center justify-center text-slate-500">
              <Icons.Clock className="w-4 h-4" />
            </div>
            <div>
              <p className="text-xs text-slate-400 font-medium uppercase tracking-wider">
                Schedule
              </p>
              <p className="font-semibold">
                {course.day}, {course.time}
              </p>
            </div>
          </div>

          <div className="flex items-center gap-3 text-slate-600">
            <div className="w-8 h-8 rounded-full bg-slate-100 flex items-center justify-center text-slate-500">
              <Icons.Info className="w-4 h-4" />
            </div>
            <div>
              <p className="text-xs text-slate-400 font-medium uppercase tracking-wider">
                Availability
              </p>
              <p className="font-semibold">
                {course.capacity - course.spots} / {course.capacity} spots
                filled
              </p>
            </div>
          </div>
        </div>

        <div className="flex flex-col gap-4">
          <div className="flex items-center justify-between p-4 bg-slate-50 rounded-xl border border-slate-100">
            <div>
              {course.spots === 0 ? (
                <span className="text-amber-600 text-sm font-bold flex items-center gap-2">
                  <span className="w-2.5 h-2.5 rounded-full bg-amber-500 animate-pulse"></span>{" "}
                  Waitlist Only
                </span>
              ) : (
                <span
                  className={`text-sm font-bold flex items-center gap-2 ${
                    course.spots < 5 ? "text-rose-500" : "text-emerald-600"
                  }`}
                >
                  <span
                    className={`w-2.5 h-2.5 rounded-full ${
                      course.spots < 5 ? "bg-rose-500" : "bg-emerald-500"
                    }`}
                  ></span>
                  {course.spots} spots remaining
                </span>
              )}
            </div>
            <p className="text-xs text-slate-500 italic">
              Limited availability
            </p>
          </div>

          <div className="grid grid-cols-2 gap-3">
            <Button variant="secondary" onClick={onClose} className="w-full">
              Cancel
            </Button>
            <Button
              variant={
                enrolled ? "danger" : course.spots === 0 ? "outline" : "primary"
              }
              className="w-full shadow-lg shadow-indigo-200"
              onClick={() => {
                if (enrolled) {
                  if (onDrop) {
                    onDrop(course);
                    return;
                  }
                  return;
                }
                if (onBook) {
                  onBook(course);
                  return;
                }
                onClose();
              }}
            >
              {enrolled
                ? dropping
                  ? "Dropping..."
                  : "Drop"
                : booking
                  ? "Booking..."
                  : course.spots === 0
                    ? "Join Waitlist"
                    : "Book Now"}
            </Button>
          </div>
        </div>
      </div>
    </Modal>
  );
};

export default CourseDetailsModal;
