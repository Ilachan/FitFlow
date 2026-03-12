import { useState } from "react";
import { useNavigate } from "react-router-dom";
import Button from "../components/ui/Button";
import Card from "../components/ui/Card";
import Input from "../components/ui/Input";
import { Icons } from "../lib/icons";
import { useAuthStore } from "../store/authStore";
import BackgroundBlobs from "../components/ui/BackgroundBlobs";
import toast from "react-hot-toast";

const Register = () => {
  const navigate = useNavigate();
  const { login } = useAuthStore();
  const [role, setRole] = useState<"student" | "manager">("student");
  const [invitationCode, setInvitationCode] = useState("");
  const [formData, setFormData] = useState({
    fullName: "",
    email: "",
    password: "",
  });

  const handleRegister = () => {
    if (!formData.fullName || !formData.email || !formData.password) {
      toast.error("Please fill in all fields.");
      return;
    }

    if (role === "manager" && !invitationCode) {
      toast.error("Invitation code is required for managers.");
      return;
    }

    // Mock invitation code validation
    if (role === "manager" && invitationCode !== "FITFLOW2026") {
      toast.error("Invalid invitation code.");
      return;
    }

    // In a real app, you would make an API call to register here.
    // fetch('/api/register', { method: 'POST', body: JSON.stringify({ ...formData, role, invitationCode }) })

    // Mock registration + auto-login
    const mockToken =
      "mock-jwt-token-" + Math.random().toString(36).substring(7);
    login(mockToken, role);

    toast.success("Registration successful! Welcome to FitFlow.");

    // Redirect to dashboard after successful registration
    navigate("/");
  };

  return (
    <div className="min-h-screen flex items-center justify-center p-4 relative">
      <BackgroundBlobs />
      <Card className="w-full max-w-md p-8 relative z-10">
        <Button
          variant="ghost"
          onClick={() => navigate("/login")}
          className="absolute top-4 left-4 p-2 text-slate-400 hover:text-slate-600"
        >
          <Icons.ArrowLeft />
        </Button>

        <div className="text-center mb-8">
          <div className="w-16 h-16 bg-linear-to-tr from-violet-600 to-indigo-600 rounded-2xl mx-auto flex items-center justify-center mb-4 shadow-lg shadow-indigo-500/30">
            <span className="text-3xl">💪</span>
          </div>
          <h1 className="text-3xl font-bold text-slate-800 mb-2">FitFlow</h1>
          <p className="text-slate-500">Start your journey today.</p>
        </div>

        <div className="space-y-4">
          {/* Role Selection */}
          <div className="flex p-1 bg-slate-100 rounded-xl mb-6">
            <button
              onClick={() => setRole("student")}
              className={`flex-1 py-2 text-sm font-medium rounded-lg transition-all ${
                role === "student"
                  ? "bg-white text-indigo-600 shadow-sm"
                  : "text-slate-500 hover:text-slate-700"
              }`}
            >
              Student
            </button>
            <button
              onClick={() => setRole("manager")}
              className={`flex-1 py-2 text-sm font-medium rounded-lg transition-all ${
                role === "manager"
                  ? "bg-white text-indigo-600 shadow-sm"
                  : "text-slate-500 hover:text-slate-700"
              }`}
            >
              Manager
            </button>
          </div>

          <div className="relative group">
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-400 group-focus-within:text-indigo-500 transition-colors">
              <Icons.User />
            </div>
            <Input
              type="text"
              placeholder="Full Name"
              value={formData.fullName}
              onChange={(e) =>
                setFormData({ ...formData, fullName: e.target.value })
              }
            />
          </div>
          <div className="relative group">
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-400 group-focus-within:text-indigo-500 transition-colors">
              <Icons.Mail />
            </div>
            <Input
              type="email"
              placeholder="Email Address"
              value={formData.email}
              onChange={(e) =>
                setFormData({ ...formData, email: e.target.value })
              }
            />
          </div>

          <div className="relative group">
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-400 group-focus-within:text-indigo-500 transition-colors">
              <Icons.Lock />
            </div>
            <Input
              type="password"
              placeholder="Password"
              value={formData.password}
              onChange={(e) =>
                setFormData({ ...formData, password: e.target.value })
              }
            />
          </div>

          {/* Invitation Code for Managers */}
          {role === "manager" && (
            <div className="relative group animate-in fade-in slide-in-from-top-2 duration-300">
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-400 group-focus-within:text-indigo-500 transition-colors">
                <Icons.Key />
              </div>
              <Input
                type="text"
                placeholder="Invitation Code"
                value={invitationCode}
                onChange={(e) => setInvitationCode(e.target.value)}
                className="border-indigo-100 focus:ring-indigo-500"
              />
              <p className="mt-1.5 text-xs text-slate-400 pl-1">
                Required for manager accounts.
              </p>
            </div>
          )}

          <Button className="w-full py-3.5 text-lg mt-2" onClick={handleRegister}>
            Create Account
          </Button>
        </div>

        <div className="mt-8 text-center text-sm text-slate-500">
          <button
            onClick={() => navigate("/login")}
            className="font-semibold text-indigo-600 hover:text-indigo-700 hover:underline"
          >
            Already have an account? Log in
          </button>
        </div>
      </Card>
    </div>
  );
};


export default Register;
