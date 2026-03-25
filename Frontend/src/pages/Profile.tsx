import React, { useCallback, useEffect, useMemo, useState } from "react";
import Input from "../components/ui/Input";
import Button from "../components/ui/Button";
import Card from "../components/ui/Card";
import Modal from "../components/ui/Modal";
import { Icons } from "../lib/icons";
import { useAuthStore } from "../store/authStore";
import toast from "react-hot-toast";
import {
  createManagerInviteCodeRequest,
  getProfileRequest,
  updateProfileRequest,
  type UpdateUserProfilePayload,
  type UserProfile,
} from "../lib/api";
import {
  formatPhoneNumberUS,
  isAgeAtLeast,
  isValidAvatarUrl,
  isValidPhoneNumberUS,
  validateEmail,
} from "../lib/validation";

const GENDER_OPTIONS = ["Male", "Female", "Other"] as const;
type GenderOption = (typeof GENDER_OPTIONS)[number];

const isGenderOption = (value: string): value is GenderOption =>
  GENDER_OPTIONS.includes(value as GenderOption);

const normalizeFromApi = (value: string | undefined) => (value || "").trim();
const normalizeGenderFromApi = (value: string | undefined) => {
  const normalized = normalizeFromApi(value);
  return isGenderOption(normalized) ? normalized : "";
};

type ProfileErrors = {
  name?: string;
  email?: string;
  avatar_url?: string;
  date_of_birth?: string;
  phone_number?: string;
};

const initialProfileState: UserProfile = {
  name: "",
  email: "",
  avatar_url: "",
  date_of_birth: "",
  gender: "",
  phone_number: "",
  address: "",
};

const Profile = () => {
  const { role, token } = useAuthStore();
  const [loading, setLoading] = useState(false);
  const [generatingInviteCode, setGeneratingInviteCode] = useState(false);
  const [initializing, setInitializing] = useState(true);
  const [showConfirmSave, setShowConfirmSave] = useState(false);
  const [generatedInviteCode, setGeneratedInviteCode] = useState("");
  const [inviteeEmail, setInviteeEmail] = useState("");
  const [expireHours, setExpireHours] = useState("");
  const [formData, setFormData] = useState<UserProfile>(initialProfileState);
  const [originalData, setOriginalData] =
    useState<UserProfile>(initialProfileState);
  const [errors, setErrors] = useState<ProfileErrors>({});

  const avatarFallback = useMemo(() => {
    const name = formData.name.trim() || "User";
    return `https://ui-avatars.com/api/?name=${encodeURIComponent(name)}&background=6366f1&color=ffffff`;
  }, [formData.name]);

  const avatarPreview = formData.avatar_url.trim() || avatarFallback;

  const isDirty = useMemo(() => {
    return JSON.stringify(formData) !== JSON.stringify(originalData);
  }, [formData, originalData]);

  const validateSingleField = useCallback(
    (field: keyof UserProfile, value: string): string | undefined => {
      if (field === "name" && !value.trim()) {
        return "Full name is required.";
      }

      if (field === "email") {
        if (!value.trim()) {
          return "Email is required.";
        }
        if (!validateEmail(value)) {
          return "Please enter a valid email address.";
        }
      }

      if (field === "avatar_url" && !isValidAvatarUrl(value)) {
        return "Avatar URL must start with http:// or https://.";
      }

      if (field === "phone_number" && !isValidPhoneNumberUS(value)) {
        return "Phone number must match (XXX) XXX-XXXX.";
      }

      if (field === "date_of_birth") {
        if (value && !isAgeAtLeast(value, 13)) {
          return "You must be at least 13 years old.";
        }
      }

      return undefined;
    },
    [],
  );

  const validateAllFields = useCallback((): boolean => {
    const nextErrors: ProfileErrors = {
      name: validateSingleField("name", formData.name),
      email: validateSingleField("email", formData.email),
      avatar_url: validateSingleField("avatar_url", formData.avatar_url),
      phone_number: validateSingleField("phone_number", formData.phone_number),
      date_of_birth: validateSingleField(
        "date_of_birth",
        formData.date_of_birth,
      ),
    };

    const compactErrors = Object.fromEntries(
      Object.entries(nextErrors).filter(([, error]) => Boolean(error)),
    ) as ProfileErrors;

    setErrors(compactErrors);
    return Object.keys(compactErrors).length === 0;
  }, [formData, validateSingleField]);

  const loadProfile = useCallback(async () => {
    if (!token) {
      setInitializing(false);
      return;
    }

    try {
      const data = await getProfileRequest(token);
      const normalizedProfile: UserProfile = {
        name: normalizeFromApi(data.name),
        email: normalizeFromApi(data.email),
        avatar_url: normalizeFromApi(data.avatar_url),
        date_of_birth: normalizeFromApi(data.date_of_birth),
        gender: normalizeGenderFromApi(data.gender),
        phone_number: normalizeFromApi(data.phone_number),
        address: normalizeFromApi(data.address),
      };

      setFormData(normalizedProfile);
      setOriginalData(normalizedProfile);
      setErrors({});
    } catch (error) {
      const message =
        error instanceof Error ? error.message : "Failed to load profile";
      toast.error(message);
    } finally {
      setInitializing(false);
    }
  }, [token]);

  useEffect(() => {
    let active = true;

    setInitializing(true);

    void (async () => {
      await loadProfile();
      if (!active) {
        return;
      }
    })();

    return () => {
      active = false;
    };
  }, [loadProfile]);

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>,
  ) => {
    const { name } = e.target;
    const fieldName = name as keyof UserProfile;
    const rawValue = e.target.value;
    const value =
      fieldName === "phone_number" ? formatPhoneNumberUS(rawValue) : rawValue;

    setFormData((prev) => ({
      ...prev,
      [fieldName]: value,
    }));
    setErrors((prev) => ({
      ...prev,
      [fieldName]: validateSingleField(fieldName, value),
    }));
  };

  const handleSave = async () => {
    if (loading || !token) {
      return;
    }

    setLoading(true);
    try {
      const payload: UpdateUserProfilePayload = {
        name: formData.name.trim(),
        email: formData.email.trim(),
      };

      const avatarUrl = formData.avatar_url.trim();
      if (avatarUrl !== "") {
        payload.avatar_url = avatarUrl;
      }

      const dateOfBirth = formData.date_of_birth.trim();
      if (dateOfBirth !== "") {
        payload.date_of_birth = dateOfBirth;
      }

      const phoneNumber = formData.phone_number.trim();
      if (phoneNumber !== "") {
        payload.phone_number = phoneNumber;
      }

      const address = formData.address.trim();
      if (address !== "") {
        payload.address = address;
      }

      const gender = formData.gender.trim();
      if (gender !== "") {
        payload.gender = gender;
      }

      await updateProfileRequest(token, payload);
      await loadProfile();
      toast.success("Profile saved successfully!");
    } catch (error) {
      const message =
        error instanceof Error ? error.message : "Failed to save profile";
      toast.error(message);
    } finally {
      setLoading(false);
    }
  };

  const openSaveConfirmation = () => {
    if (!validateAllFields()) {
      toast.error("Please fix the validation errors before saving.");
      return;
    }

    if (!isDirty) {
      toast("No changes to save.", { icon: "ℹ️" });
      return;
    }

    setShowConfirmSave(true);
  };

  const handleGenerateInviteCode = async () => {
    if (generatingInviteCode || !token) {
      return;
    }

    if (!inviteeEmail.trim()) {
      toast.error("Invitee email is required.");
      return;
    }

    if (!validateEmail(inviteeEmail)) {
      toast.error("Please enter a valid invitee email.");
      return;
    }

    const expireHoursValue = Number(expireHours);
    if (
      !Number.isInteger(expireHoursValue) ||
      expireHoursValue < 1 ||
      expireHoursValue > 720
    ) {
      toast.error("Expire hours must be an integer between 1 and 720.");
      return;
    }

    setGeneratingInviteCode(true);
    try {
      const data = await createManagerInviteCodeRequest(token, {
        invitee_email: inviteeEmail.trim(),
        expire_hours: expireHoursValue,
      });
      setGeneratedInviteCode(data.code || "");
      toast.success("Invite code generated successfully.");
    } catch (error) {
      const message =
        error instanceof Error
          ? error.message
          : "Failed to generate invite code";
      toast.error(message);
    } finally {
      setGeneratingInviteCode(false);
    }
  };

  const handleCopyInviteCode = async () => {
    if (!generatedInviteCode) {
      toast.error("Please generate an invite code first.");
      return;
    }

    try {
      await navigator.clipboard.writeText(generatedInviteCode);
      toast.success("Invite code copied.");
    } catch {
      toast.error("Failed to copy invite code.");
    }
  };

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-3xl font-bold bg-clip-text text-transparent bg-linear-to-r from-violet-600 to-indigo-600">
            User Profile
          </h1>
          <p className="text-slate-500 mt-2">
            Manage your personal information and preferences
          </p>
          {initializing && (
            <p className="text-xs text-slate-400 mt-2">Loading profile...</p>
          )}
        </div>
        <Button
          onClick={openSaveConfirmation}
          className={loading ? "opacity-50 cursor-not-allowed" : ""}
        >
          {loading ? (
            <span className="flex items-center gap-2">
              <span className="w-4 h-4 border-2 border-white/20 border-t-white rounded-full animate-spin" />
              Saving...
            </span>
          ) : (
            <>
              <Icons.Check />
              Save Changes
            </>
          )}
        </Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="space-y-6">
          <Card className="p-6 text-center space-y-4">
            <div className="relative w-32 h-32 mx-auto">
              <img
                src={avatarPreview}
                alt="Profile"
                onError={(event) => {
                  event.currentTarget.src = avatarFallback;
                }}
                className="w-full h-full rounded-full object-cover border-4 border-indigo-50"
              />
              <div className="absolute -bottom-2 left-1/2 -translate-x-1/2 px-3 py-1 bg-indigo-600 text-white text-xs font-bold rounded-full shadow-lg capitalize">
                {role}
              </div>
            </div>
            <div className="pt-2">
              <h2 className="text-xl font-bold text-slate-800">
                {formData.name}
              </h2>
              <p className="text-indigo-600 font-medium">{formData.email}</p>
            </div>
          </Card>

          {role === "supermanager" && (
            <Card className="p-6 space-y-4 border-indigo-100 bg-indigo-50/30">
              <h3 className="text-lg font-bold text-slate-800 flex items-center gap-2">
                <span className="w-1 h-6 bg-indigo-500 rounded-full" />
                Manager Invite Code
              </h3>
              <p className="text-sm text-slate-500">
                Generate a manager invitation code and share it securely.
              </p>
              <div className="space-y-3">
                <div className="relative group">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-400 group-focus-within:text-indigo-500 transition-colors">
                    <Icons.Mail />
                  </div>
                  <Input
                    placeholder="Invitee email"
                    value={inviteeEmail}
                    onChange={(e) => setInviteeEmail(e.target.value)}
                    className="px-4 py-3 text-sm"
                  />
                </div>
                <div className="relative group">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-400 group-focus-within:text-indigo-500 transition-colors text-xs font-bold">
                    H
                  </div>
                  <Input
                    type="number"
                    min={1}
                    max={720}
                    placeholder="Expire hours"
                    value={expireHours}
                    onChange={(e) => setExpireHours(e.target.value)}
                    className="pl-4 py-3 text-sm"
                  />
                </div>
                <div className="relative group">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-400 group-focus-within:text-indigo-500 transition-colors">
                    <Icons.Key className="w-4 h-4" />
                  </div>
                  <Input
                    placeholder="Generated invitation code"
                    value={generatedInviteCode}
                    readOnly
                    className="px-4 py-3 text-sm bg-slate-100 cursor-not-allowed"
                  />
                </div>
                <div className="grid grid-cols-1 gap-2">
                  <Button
                    variant="outline"
                    onClick={() => {
                      void handleGenerateInviteCode();
                    }}
                    className={`w-full py-2.5 text-sm ${generatingInviteCode ? "opacity-50 cursor-not-allowed" : ""}`}
                  >
                    {generatingInviteCode ? "Generating..." : "Generate Code"}
                  </Button>
                  <Button
                    variant="secondary"
                    onClick={() => {
                      void handleCopyInviteCode();
                    }}
                    className="w-full py-2.5 text-sm"
                  >
                    Copy Code
                  </Button>
                </div>
                <p className="text-xs text-slate-500">
                  Expire hours must be between 1 and 720.
                </p>
              </div>
            </Card>
          )}
        </div>

        <div className="md:col-span-2 space-y-6">
          <Card className="p-6 space-y-4">
            <h3 className="text-lg font-bold text-slate-800 mb-4 flex items-center gap-2">
              <span className="w-1 h-6 bg-indigo-500 rounded-full" />
              Account Information
            </h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <label className="text-sm font-medium text-slate-700">
                  Full Name
                </label>
                <Input
                  name="name"
                  value={formData.name}
                  onChange={handleChange}
                  placeholder="Enter full name"
                />
                {errors.name && (
                  <p className="text-xs text-rose-500">{errors.name}</p>
                )}
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium text-slate-700">
                  Email
                </label>
                <Input
                  name="email"
                  type="email"
                  value={formData.email}
                  onChange={handleChange}
                  placeholder="Enter email"
                />
                {errors.email && (
                  <p className="text-xs text-rose-500">{errors.email}</p>
                )}
              </div>
              <div className="space-y-2 md:col-span-2">
                <label className="text-sm font-medium text-slate-700">
                  Avatar URL
                </label>
                <Input
                  name="avatar_url"
                  value={formData.avatar_url}
                  onChange={handleChange}
                  placeholder="https://example.com/avatar.png"
                />
                {errors.avatar_url && (
                  <p className="text-xs text-rose-500">{errors.avatar_url}</p>
                )}
              </div>
            </div>
          </Card>

          <Card className="p-6 space-y-4">
            <h3 className="text-lg font-bold text-slate-800 mb-4 flex items-center gap-2">
              <span className="w-1 h-6 bg-purple-500 rounded-full" />
              Personal Details
            </h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <label className="text-sm font-medium text-slate-700">
                  Date of Birth
                </label>
                <Input
                  name="date_of_birth"
                  type="date"
                  value={formData.date_of_birth}
                  onChange={handleChange}
                />
                {errors.date_of_birth && (
                  <p className="text-xs text-rose-500">
                    {errors.date_of_birth}
                  </p>
                )}
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium text-slate-700">
                  Gender
                </label>
                <select
                  name="gender"
                  value={formData.gender}
                  onChange={handleChange}
                  className="w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none transition-all"
                >
                  {formData.gender.trim() === "" && (
                    <option value="">Select Gender</option>
                  )}
                  <option value="Male">Male</option>
                  <option value="Female">Female</option>
                  <option value="Other">Other</option>
                </select>
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium text-slate-700">
                  Phone Number
                </label>
                <Input
                  name="phone_number"
                  value={formData.phone_number}
                  onChange={handleChange}
                  placeholder="(555) 000-0000"
                />
                {errors.phone_number && (
                  <p className="text-xs text-rose-500">{errors.phone_number}</p>
                )}
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium text-slate-700">
                  Address
                </label>
                <Input
                  name="address"
                  value={formData.address}
                  onChange={handleChange}
                  placeholder="Street Address"
                />
              </div>
            </div>
          </Card>
        </div>
      </div>

      <Modal isOpen={showConfirmSave} onClose={() => setShowConfirmSave(false)}>
        <div className="p-6">
          <h3 className="text-lg font-bold text-slate-900">
            Confirm Profile Update
          </h3>
          <p className="text-sm text-slate-500 mt-2">
            Are you sure you want to save these changes to your profile?
          </p>
          <div className="mt-6 grid grid-cols-2 gap-3">
            <Button
              variant="secondary"
              onClick={() => setShowConfirmSave(false)}
            >
              Cancel
            </Button>
            <Button
              onClick={async () => {
                setShowConfirmSave(false);
                await handleSave();
              }}
            >
              Confirm Save
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
};

export default Profile;
