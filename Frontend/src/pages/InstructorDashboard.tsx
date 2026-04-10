import { useState } from "react";
import Card from "../components/ui/Card";
import Button from "../components/ui/Button";
import Badge from "../components/ui/Badge";
import Modal from "../components/ui/Modal";
import Input from "../components/ui/Input";
import { Icons } from "../lib/icons";
import BackgroundBlobs from "../components/ui/BackgroundBlobs";
import toast from "react-hot-toast";

type Student = {
  id: number;
  name: string;
  email: string;
  status: "pending" | "present" | "absent" | "walk-in";
};

type ClassData = {
  id: number;
  name: string;
  time: string;
  capacity: number;
  status: "upcoming" | "active" | "completed";
  students: Student[];
};

const InstructorDashboard = () => {
  const isPreviewInstructor =
    sessionStorage.getItem("instructor_preview") === "true";
  const [selectedClassId, setSelectedClassId] = useState<number | null>(null);
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [newStudent, setNewStudent] = useState({ name: "", email: "" });

  // Mock data with student lists
  const [classes, setClasses] = useState<ClassData[]>([
    {
      id: 1,
      name: "Morning Flow Yoga",
      time: "08:00 AM - 09:00 AM",
      capacity: 20,
      status: "completed",
      students: [
        {
          id: 101,
          name: "Alice Johnson",
          email: "alice@example.com",
          status: "present",
        },
        {
          id: 102,
          name: "Bob Smith",
          email: "bob@example.com",
          status: "present",
        },
        {
          id: 103,
          name: "Charlie Davis",
          email: "charlie@example.com",
          status: "absent",
        },
      ],
    },
    {
      id: 2,
      name: "Power Vinyasa",
      time: "11:30 AM - 12:30 PM",
      capacity: 15,
      status: "active",
      students: [
        {
          id: 201,
          name: "David Miller",
          email: "david@example.com",
          status: "pending",
        },
        {
          id: 202,
          name: "Emma Wilson",
          email: "emma@example.com",
          status: "pending",
        },
        {
          id: 203,
          name: "Frank Wright",
          email: "frank@example.com",
          status: "pending",
        },
        {
          id: 204,
          name: "Grace Lee",
          email: "grace@example.com",
          status: "pending",
        },
      ],
    },
    {
      id: 3,
      name: "Gentle Hatha",
      time: "05:00 PM - 06:00 PM",
      capacity: 10,
      status: "upcoming",
      students: [
        {
          id: 301,
          name: "Henry Ford",
          email: "henry@example.com",
          status: "pending",
        },
        {
          id: 302,
          name: "Ivy Chen",
          email: "ivy@example.com",
          status: "pending",
        },
      ],
    },
  ]);

  const selectedClass = classes.find((c) => c.id === selectedClassId);

  const toggleAttendance = (
    studentId: number,
    status: "present" | "absent",
  ) => {
    if (!selectedClassId) return;
    setClasses((prev) =>
      prev.map((c) => {
        if (c.id !== selectedClassId) return c;
        return {
          ...c,
          students: c.students.map((s) =>
            s.id === studentId ? { ...s, status } : s,
          ),
        };
      }),
    );
  };

  const markAllPresent = () => {
    if (!selectedClassId) return;
    setClasses((prev) =>
      prev.map((c) => {
        if (c.id !== selectedClassId) return c;
        return {
          ...c,
          students: c.students.map((s) => ({ ...s, status: "present" })),
        };
      }),
    );
  };

  const handleAddWalkIn = () => {
    if (!selectedClassId || !newStudent.name || !newStudent.email) {
      toast.error("Please provide both name and email.");
      return;
    }

    if (
      selectedClass &&
      selectedClass.students.length >= selectedClass.capacity
    ) {
      const confirmed = window.confirm(
        "Class is already at full capacity. Add walk-in anyway?",
      );
      if (!confirmed) return;
    }

    setClasses((prev) =>
      prev.map((c) => {
        if (c.id !== selectedClassId) return c;
        const walkIn: Student = {
          id: Date.now(),
          name: newStudent.name,
          email: newStudent.email,
          status: "present",
        };
        return {
          ...c,
          students: [...c.students, walkIn],
        };
      }),
    );

    setIsAddModalOpen(false);
    setNewStudent({ name: "", email: "" });
    toast.success(`${newStudent.name} added and marked as present!`);
  };

  const handleSubmitAttendance = () => {
    if (!selectedClass) {
      return;
    }

    toast.success("Attendance submitted in demo mode.");
  };

  return (
    <div className="min-h-screen bg-slate-50 p-6 relative">
      <BackgroundBlobs />
      <div className="max-w-7xl mx-auto relative z-10">
        {isPreviewInstructor && (
          <Card className="mb-6 p-4 border-dashed border-indigo-200 bg-indigo-50/70">
            <div className="flex items-center justify-between gap-4 flex-wrap">
              <div>
                <Badge className="bg-indigo-100 text-indigo-700 mb-2">
                  Demo Mode
                </Badge>
              </div>
              <p className="text-xs text-slate-500">
                Mark Present, set No-show, then submit attendance.
              </p>
            </div>
          </Card>
        )}

        <header className="mb-8 flex flex-col md:flex-row md:items-end justify-between gap-4">
          <div className="flex items-center gap-4">
            <div className="w-14 h-14 bg-linear-to-tr from-indigo-600 to-violet-600 rounded-2xl flex items-center justify-center text-white shadow-lg shadow-indigo-200">
              <Icons.User className="w-8 h-8" />
            </div>
            <div>
              <h1 className="text-3xl font-bold text-slate-800 tracking-tight">
                Instructor Dashboard
              </h1>
              <p className="text-slate-500 font-medium">
                Welcome back, Instructor One
              </p>
            </div>
          </div>
          <div className="flex gap-3">
            {isPreviewInstructor && (
              <Badge className="bg-slate-900 text-white px-4 py-1.5 text-sm">
                Preview
              </Badge>
            )}
            <Badge className="bg-indigo-50 text-indigo-700 px-4 py-1.5 text-sm">
              Yoga Specialist
            </Badge>
            <Badge className="bg-emerald-50 text-emerald-700 px-4 py-1.5 text-sm">
              Level 4 Certified
            </Badge>
          </div>
        </header>

        <div className="grid grid-cols-1 xl:grid-cols-12 gap-8">
          {/* Class List */}
          <div className="xl:col-span-4 space-y-4">
            <h2 className="text-lg font-bold text-slate-700 flex items-center gap-2 mb-4">
              <Icons.Calendar className="w-5 h-5 text-indigo-500" />
              Today's Schedule
            </h2>
            {classes.map((cls) => (
              <Card
                key={cls.id}
                onClick={() => setSelectedClassId(cls.id)}
                className={`group relative overflow-hidden p-4 cursor-pointer transition-all border-2 ${
                  selectedClassId === cls.id
                    ? "border-indigo-500 ring-4 ring-indigo-50 shadow-md bg-indigo-50/30"
                    : "border-transparent hover:border-slate-200"
                }`}
              >
                <span
                  className={`absolute left-0 top-0 h-full w-1 transition-colors ${
                    selectedClassId === cls.id
                      ? "bg-indigo-500"
                      : "bg-transparent group-hover:bg-slate-200"
                  }`}
                />
                <div className="flex justify-between items-start mb-2">
                  <h3 className="font-bold text-slate-800">{cls.name}</h3>
                  {cls.status === "active" && (
                    <span className="flex h-2 w-2 rounded-full bg-emerald-500 animate-ping" />
                  )}
                </div>
                <div className="space-y-1.5 text-sm text-slate-500">
                  <div className="flex items-center gap-2">
                    <Icons.Clock className="w-4 h-4" />
                    {cls.time}
                  </div>
                  <div className="flex items-center gap-2">
                    <Icons.User className="w-4 h-4" />
                    {cls.students.length} / {cls.capacity} Students
                  </div>
                </div>
                <div className="mt-3 flex items-center justify-between">
                  <Badge
                    className={
                      cls.status === "active"
                        ? "bg-emerald-100 text-emerald-700"
                        : cls.status === "completed"
                          ? "bg-slate-100 text-slate-600"
                          : "bg-blue-100 text-blue-700"
                    }
                  >
                    {cls.status.charAt(0).toUpperCase() + cls.status.slice(1)}
                  </Badge>
                  <span
                    className={`flex h-8 w-8 items-center justify-center rounded-full border transition-all ${
                      selectedClassId === cls.id
                        ? "border-indigo-200 bg-indigo-100 text-indigo-600 shadow-sm"
                        : "border-slate-200 bg-white text-slate-400"
                    }`}
                  >
                    <Icons.ChevronRight className="w-4 h-4" />
                  </span>
                </div>
              </Card>
            ))}
          </div>

          {/* Student Roster / Attendance Management */}
          <div className="xl:col-span-8">
            {selectedClass ? (
              <div className="space-y-6 animate-in fade-in slide-in-from-right-4 duration-300">
                <div className="bg-white rounded-3xl p-6 shadow-sm border border-slate-100">
                  <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
                    <div>
                      <h2 className="text-2xl font-bold text-slate-800 mb-1">
                        {selectedClass.name}
                      </h2>
                      <p className="text-slate-500 flex items-center gap-2">
                        <Icons.Clock className="w-4 h-4" />
                        {selectedClass.time} • Student Roster
                      </p>
                    </div>
                    <div className="flex flex-wrap gap-2">
                      <Button
                        variant="secondary"
                        onClick={markAllPresent}
                        className="border-slate-200 bg-white text-slate-700 hover:bg-slate-50 hover:border-slate-300"
                      >
                        Mark All Present
                      </Button>
                      <Button
                        variant="primary"
                        className="bg-slate-900 text-white shadow-md shadow-slate-200 hover:bg-slate-800 hover:shadow-slate-300"
                        onClick={() => setIsAddModalOpen(true)}
                      >
                        <Icons.Plus className="w-4 h-4 mr-2" />
                        Add Walk-in
                      </Button>
                      <Button
                        className="bg-gradient-to-r from-indigo-600 to-cyan-600 text-white shadow-lg shadow-indigo-200 hover:from-indigo-500 hover:to-cyan-500"
                        onClick={handleSubmitAttendance}
                      >
                        Submit Attendance
                      </Button>
                    </div>
                  </div>

                  <div className="overflow-hidden border border-slate-100 rounded-2xl">
                    <table className="w-full text-left border-collapse">
                      <thead>
                        <tr className="bg-slate-50 text-slate-400 text-xs uppercase tracking-wider font-semibold">
                          <th className="px-6 py-4">Student</th>
                          <th className="px-6 py-4 text-center">Status</th>
                          <th className="px-6 py-4 text-right">Actions</th>
                        </tr>
                      </thead>
                      <tbody className="divide-y divide-slate-50">
                        {selectedClass.students.map((student) => (
                          <tr
                            key={student.id}
                            className="hover:bg-slate-50/50 transition-colors"
                          >
                            <td className="px-6 py-4">
                              <div className="flex items-center gap-3">
                                <div className="w-10 h-10 rounded-full bg-slate-100 flex items-center justify-center text-slate-500 font-bold border border-slate-200">
                                  {student.name.charAt(0)}
                                </div>
                                <div>
                                  <div className="flex items-center gap-2">
                                    <p className="font-bold text-slate-700">
                                      {student.name}
                                    </p>
                                    {student.id > 1000 && (
                                      <Badge className="bg-amber-50 text-amber-700 text-[10px] py-0 px-1.5">
                                        WALK-IN
                                      </Badge>
                                    )}
                                  </div>
                                  <p className="text-xs text-slate-400">
                                    {student.email}
                                  </p>
                                </div>
                              </div>
                            </td>
                            <td className="px-6 py-4 text-center">
                              {student.status === "present" ? (
                                <Badge className="bg-emerald-100 text-emerald-700">
                                  Present
                                </Badge>
                              ) : student.status === "absent" ? (
                                <Badge className="bg-rose-100 text-rose-700">
                                  No-show
                                </Badge>
                              ) : (
                                <Badge className="bg-slate-100 text-slate-500">
                                  Pending
                                </Badge>
                              )}
                            </td>
                            <td className="px-6 py-4">
                              <div className="flex justify-end gap-2">
                                <button
                                  onClick={() =>
                                    toggleAttendance(student.id, "present")
                                  }
                                  className={`p-2 rounded-lg transition-all ${
                                    student.status === "present"
                                      ? "bg-emerald-500 text-white shadow-md shadow-emerald-100"
                                      : "bg-white text-slate-300 border border-slate-200 hover:border-emerald-500 hover:text-emerald-500"
                                  }`}
                                  title="Mark Present"
                                >
                                  <Icons.Check />
                                </button>
                                <button
                                  onClick={() =>
                                    toggleAttendance(student.id, "absent")
                                  }
                                  className={`p-2 rounded-lg transition-all ${
                                    student.status === "absent"
                                      ? "bg-rose-500 text-white shadow-md shadow-rose-100"
                                      : "bg-white text-slate-300 border border-slate-200 hover:border-rose-500 hover:text-rose-500"
                                  }`}
                                  title="Mark No-show"
                                >
                                  <Icons.X className="w-5 h-5" />
                                </button>
                              </div>
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </div>
              </div>
            ) : (
              <div className="h-full flex flex-col items-center justify-center text-center p-12 bg-white rounded-3xl border-2 border-dashed border-slate-200 text-slate-400">
                <div className="w-20 h-20 bg-slate-50 rounded-full flex items-center justify-center mb-4">
                  <Icons.Calendar className="w-10 h-10" />
                </div>
                <h3 className="text-xl font-bold text-slate-600 mb-2">
                  Select a Class
                </h3>
                <p className="max-w-xs mx-auto">
                  Choose a class from your schedule to view the student roster
                  and manage attendance.
                </p>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Add Walk-in Modal */}
      <Modal isOpen={isAddModalOpen} onClose={() => setIsAddModalOpen(false)}>
        <div className="p-8">
          <div className="text-center mb-8">
            <div className="w-16 h-16 bg-slate-100 rounded-2xl mx-auto flex items-center justify-center mb-4">
              <Icons.User className="w-8 h-8 text-slate-600" />
            </div>
            <h2 className="text-2xl font-bold text-slate-800">
              Add Walk-in Student
            </h2>
            <p className="text-slate-500">Adding to {selectedClass?.name}</p>
          </div>

          <div className="space-y-6">
            <div className="space-y-2">
              <label className="text-sm font-semibold text-slate-700 ml-1">
                Full Name
              </label>
              <Input
                placeholder="e.g. John Doe"
                value={newStudent.name}
                onChange={(e) =>
                  setNewStudent({ ...newStudent, name: e.target.value })
                }
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-semibold text-slate-700 ml-1">
                Email Address
              </label>
              <Input
                type="email"
                placeholder="john@example.com"
                value={newStudent.email}
                onChange={(e) =>
                  setNewStudent({ ...newStudent, email: e.target.value })
                }
              />
            </div>

            <div className="pt-4 flex gap-3">
              <Button
                variant="ghost"
                className="flex-1"
                onClick={() => setIsAddModalOpen(false)}
              >
                Cancel
              </Button>
              <Button
                className="flex-1 bg-indigo-600 hover:bg-indigo-700"
                onClick={handleAddWalkIn}
              >
                Add Student
              </Button>
            </div>
          </div>
        </div>
      </Modal>
    </div>
  );
};

export default InstructorDashboard;
