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

import BackgroundBlobs from "../components/ui/BackgroundBlobs";

const Login = () => {
  const navigate = useNavigate();
  const { login } = useAuthStore();
  const [role, setRole] = useState<"student" | "manager">("student");
  const [formData, setFormData] = useState({
    email: "",
    password: "",
  });
  const [loading, setLoading] = useState(false);

  const handleLogin = async () => {
    if (loading) {
      return;
    }

    if (!formData.email || !formData.password) {
      toast.error("Please enter email and password.");
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
          <p className="text-slate-500">Let's get moving.</p>
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
              <Icons.Mail />
            </div>
            <Input
              type="email"
              placeholder="Email Address"
              value={formData.email}
              onChange={(e) => setFormData({ ...formData, email: e.target.value })}
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
              onChange={(e) => setFormData({ ...formData, password: e.target.value })}
            />
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
