import { useState } from "react";
import { useNavigate } from "react-router-dom";
import Button from "../components/ui/Button";
import Card from "../components/ui/Card";
import Input from "../components/ui/Input";
import Modal from "../components/ui/Modal";
import { Icons } from "../lib/icons";
import { useAuthStore } from "../store/authStore";
import BackgroundBlobs from "../components/ui/BackgroundBlobs";
import toast from "react-hot-toast";
import {
  extractUserIdFromToken,
  loginRequest,
  registerManagerRequest,
  registerStudentRequest,
  roleIdToFrontendRole,
} from "../lib/api";
import {
  evaluatePassword,
  getPasswordStrength,
  isPasswordValid,
  validateEmail,
} from "../lib/validation";

type RegisterErrors = {
  fullName?: string;
  email?: string;
  password?: string;
  invitationCode?: string;
};

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
  const [loading, setLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [showSuccessModal, setShowSuccessModal] = useState(false);
  const [errors, setErrors] = useState<RegisterErrors>({});

  const passwordChecks = evaluatePassword(formData.password);
  const passwordStrength = getPasswordStrength(formData.password);

  const validateForm = (): boolean => {
    const nextErrors: RegisterErrors = {};

    if (!formData.fullName.trim()) {
      nextErrors.fullName = "Full name is required.";
    }

    if (!formData.email.trim()) {
      nextErrors.email = "Email is required.";
    } else if (!validateEmail(formData.email)) {
      nextErrors.email = "Please enter a valid email address.";
    }

    if (!formData.password) {
      nextErrors.password = "Password is required.";
    } else if (!isPasswordValid(formData.password)) {
      nextErrors.password =
        "Password must be 8+ chars and include uppercase, lowercase, and a number.";
    }

    if (role === "manager" && !invitationCode.trim()) {
      nextErrors.invitationCode = "Invitation code is required for managers.";
    }

    setErrors(nextErrors);
    return Object.keys(nextErrors).length === 0;
  };

  const handleRegister = async () => {
    if (loading) {
      return;
    }

    if (!validateForm()) {
      toast.error("Please correct the highlighted fields.");
      return;
    }

    setLoading(true);
    try {
      if (role === "manager") {
        await registerManagerRequest(
          formData.fullName,
          formData.email,
          formData.password,
          invitationCode,
        );
      } else {
        await registerStudentRequest(
          formData.fullName,
          formData.email,
          formData.password,
        );
      }

      // Keep existing UX: auto-login immediately after successful registration.
      const loginData = await loginRequest(formData.email, formData.password);
      const frontendRole = roleIdToFrontendRole(loginData.role_id);
      const userId = extractUserIdFromToken(loginData.token);

      login(loginData.token, frontendRole, userId);
      setShowSuccessModal(true);
    } catch (error) {
      const message =
        error instanceof Error ? error.message : "Registration failed";
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
              onClick={() => {
                setRole("student");
                setErrors((prev) => ({ ...prev, invitationCode: undefined }));
              }}
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
              onChange={(e) => {
                const fullName = e.target.value;
                setFormData({ ...formData, fullName });
                setErrors((prev) => ({
                  ...prev,
                  fullName: fullName.trim() ? undefined : "Full name is required.",
                }));
              }}
            />
            {errors.fullName && (
              <p className="mt-1.5 text-xs text-rose-500">{errors.fullName}</p>
            )}
          </div>
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
                  password: !password
                    ? "Password is required."
                    : isPasswordValid(password)
                      ? undefined
                      : "Password must be 8+ chars and include uppercase, lowercase, and a number.",
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

            <div className="mt-2 rounded-lg border border-slate-200 bg-white p-3">
              <div className="mb-2 flex items-center justify-between text-xs">
                <span className="font-semibold text-slate-600">Password strength</span>
                <span
                  className={`font-semibold capitalize ${
                    passwordStrength === "strong"
                      ? "text-emerald-600"
                      : passwordStrength === "medium"
                        ? "text-amber-600"
                        : "text-rose-500"
                  }`}
                >
                  {passwordStrength}
                </span>
              </div>
              <div className="mb-2 h-2 overflow-hidden rounded-full bg-slate-100">
                <div
                  className={`h-full transition-all ${
                    passwordStrength === "strong"
                      ? "w-full bg-emerald-500"
                      : passwordStrength === "medium"
                        ? "w-2/3 bg-amber-500"
                        : "w-1/3 bg-rose-500"
                  }`}
                />
              </div>
              <ul className="space-y-1 text-xs text-slate-500">
                <li className={passwordChecks.hasMinLength ? "text-emerald-600" : ""}>
                  At least 8 characters
                </li>
                <li className={passwordChecks.hasUppercase ? "text-emerald-600" : ""}>
                  Contains an uppercase letter
                </li>
                <li className={passwordChecks.hasLowercase ? "text-emerald-600" : ""}>
                  Contains a lowercase letter
                </li>
                <li className={passwordChecks.hasNumber ? "text-emerald-600" : ""}>
                  Contains a number
                </li>
              </ul>
            </div>
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
                onChange={(e) => {
                  const code = e.target.value;
                  setInvitationCode(code);
                  setErrors((prev) => ({
                    ...prev,
                    invitationCode:
                      role === "manager" && !code.trim()
                        ? "Invitation code is required for managers."
                        : undefined,
                  }));
                }}
                className="border-indigo-100 focus:ring-indigo-500"
              />
              <p className="mt-1.5 text-xs text-slate-400 pl-1">
                Required for manager accounts.
              </p>
              {errors.invitationCode && (
                <p className="mt-1 text-xs text-rose-500">{errors.invitationCode}</p>
              )}
            </div>
          )}

          <Button
            className="w-full py-3.5 text-lg mt-2"
            onClick={handleRegister}
          >
            {loading ? "Creating Account..." : "Create Account"}
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

      <Modal
        isOpen={showSuccessModal}
        onClose={() => {
          setShowSuccessModal(false);
          toast.success("Registration successful! Welcome to FitFlow.");
          navigate("/courses");
        }}
      >
        <div className="p-6 text-center">
          <div className="w-12 h-12 rounded-full bg-emerald-100 text-emerald-600 mx-auto mb-4 flex items-center justify-center">
            <Icons.Check />
          </div>
          <h3 className="text-lg font-bold text-slate-900">Registration Complete</h3>
          <p className="text-sm text-slate-500 mt-2 mb-5">
            Your account is ready and you are now logged in.
          </p>
          <Button
            className="w-full"
            onClick={() => {
              setShowSuccessModal(false);
              toast.success("Registration successful! Welcome to FitFlow.");
              navigate("/courses");
            }}
          >
            Go to Courses
          </Button>
        </div>
      </Modal>
    </div>
  );
};


export default Register;
