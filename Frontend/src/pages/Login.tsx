import { useState } from "react";
import { useNavigate } from "react-router-dom";
import Button from "../components/ui/Button";
import Card from "../components/ui/Card";
import Input from "../components/ui/Input";
import { Icons } from "../lib/icons";
import { useAuthStore } from "../store/authStore";
import toast from "react-hot-toast";
import {
  extractUserIdFromToken,
  loginRequest,
  roleIdToFrontendRole,
} from "../lib/api";
import { validateEmail } from "../lib/validation";

import BackgroundBlobs from "../components/ui/BackgroundBlobs";

const Login = () => {
  const navigate = useNavigate();
  const { login } = useAuthStore();
  const [formData, setFormData] = useState({
    email: "",
    password: "",
  });
  const [selectedRole, setSelectedRole] = useState<"student" | "manager">("student");
  const [loading, setLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [errors, setErrors] = useState<{ email?: string; password?: string }>({});

  const validateForm = () => {
    const nextErrors: { email?: string; password?: string } = {};

    if (!formData.email.trim()) {
      nextErrors.email = "Email is required.";
    } else if (!validateEmail(formData.email)) {
      nextErrors.email = "Please enter a valid email address.";
    }

    if (!formData.password.trim()) {
      nextErrors.password = "Password is required.";
    }

    setErrors(nextErrors);
    return Object.keys(nextErrors).length === 0;
  };

  const handleLogin = async () => {
    if (loading) {
      return;
    }

    if (!validateForm()) {
      toast.error("Please correct the highlighted fields.");
      return;
    }

    setLoading(true);
    try {
      const data = await loginRequest(formData.email, formData.password);
      const frontendRole = roleIdToFrontendRole(data.role_id);
      const userId = extractUserIdFromToken(data.token);

      login(data.token, frontendRole, userId);
      toast.success("Welcome back! Successfully logged in.");
      navigate("/");
    } catch (error) {
      const message = error instanceof Error ? error.message : "Login failed";
      toast.error(message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center p-4 relative">
      <BackgroundBlobs />
      <Card className="w-full max-w-md p-8 relative z-10">
        <Button
          variant="ghost"
          onClick={() => navigate("/courses")}
          className="absolute top-4 left-4 p-2 text-slate-400 hover:text-slate-600"
        >
          <Icons.ArrowLeft />
        </Button>

        <div className="text-center mb-8">
          <div className="w-16 h-16 bg-linear-to-tr from-violet-600 to-indigo-600 rounded-2xl mx-auto flex items-center justify-center mb-4 shadow-lg shadow-indigo-500/30">
            <span className="text-3xl">💪</span>
          </div>
          <h1 className="text-3xl font-bold text-slate-800 mb-2">FitFlow</h1>
          <p className="text-slate-500 mb-6">Let's get moving.</p>

          <div className="flex bg-slate-100 p-1 rounded-xl mb-4">
            <button
              onClick={() => setSelectedRole("student")}
              className={`flex-1 py-2 text-sm font-semibold rounded-lg transition-all ${
                selectedRole === "student"
                  ? "bg-white text-indigo-600 shadow-sm"
                  : "text-slate-500 hover:text-slate-700"
              }`}
            >
              Student
            </button>
            <button
              onClick={() => setSelectedRole("manager")}
              className={`flex-1 py-2 text-sm font-semibold rounded-lg transition-all ${
                selectedRole === "manager"
                  ? "bg-white text-indigo-600 shadow-sm"
                  : "text-slate-500 hover:text-slate-700"
              }`}
            >
              Manager
            </button>
          </div>
        </div>

        <div className="space-y-4">
          <div className="relative group">
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-400 group-focus-within:text-indigo-500 transition-colors">
              <Icons.Mail />
            </div>
            <Input
              type="email"
              placeholder="Email Address"
              value={formData.email}
              onChange={(e) => {
                const email = e.target.value;
                setFormData({ ...formData, email });
                setErrors((prev) => ({
                  ...prev,
                  email: !email.trim()
                    ? "Email is required."
                    : validateEmail(email)
                      ? undefined
                      : "Please enter a valid email address.",
                }));
              }}
            />
            {errors.email && (
              <p className="mt-1.5 text-xs text-rose-500">{errors.email}</p>
            )}
          </div>

          <div className="relative group">
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-400 group-focus-within:text-indigo-500 transition-colors">
              <Icons.Lock />
            </div>
            <Input
              type={showPassword ? "text" : "password"}
              placeholder="Password"
              value={formData.password}
              onChange={(e) => {
                const password = e.target.value;
                setFormData({ ...formData, password });
                setErrors((prev) => ({
                  ...prev,
                  password: !password.trim()
                    ? "Password is required."
                    : undefined,
                }));
              }}
            />
            <button
              type="button"
              onClick={() => setShowPassword((prev) => !prev)}
              className="absolute right-3 top-1/2 -translate-y-1/2 text-xs font-semibold text-slate-500 hover:text-indigo-600"
            >
              {showPassword ? "Hide" : "Show"}
            </button>
            {errors.password && (
              <p className="mt-1.5 text-xs text-rose-500">{errors.password}</p>
            )}
          </div>

          <Button className="w-full py-3.5 text-lg mt-2" onClick={handleLogin}>
            {loading ? "Signing In..." : "Sign In"}
          </Button>
        </div>

        <div className="mt-8 text-center text-sm text-slate-500">
          <button
            onClick={() => navigate("/register")}
            className="font-semibold text-indigo-600 hover:text-indigo-700 hover:underline"
          >
            Don't have an account? Sign up
          </button>
        </div>
      </Card>
    </div>
  );
};


export default Login;
