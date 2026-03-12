import React, { useState } from "react";
import Input from "../components/ui/Input";
import Button from "../components/ui/Button";
import Card from "../components/ui/Card";
import { Icons } from "../lib/icons";
import { useAuthStore } from "../store/authStore";
import toast from "react-hot-toast";

const Profile = () => {
  const { role, login, token } = useAuthStore();
  const [loading, setLoading] = useState(false);
  const [upgrading, setUpgrading] = useState(false);
  const [invitationCode, setInvitationCode] = useState("");
  const [formData, setFormData] = useState({
    name: "John Doe",
    email: "john@example.com",
    date_of_birth: "2000-01-01",
    gender: "Male",
    phone_number: "+1 (555) 123-4567",
    address: "123 Fitness St, Wellness City",
  });

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSave = () => {
    setLoading(true);
    // Mock API call
    setTimeout(() => {
      console.log("Saved profile:", formData);
      setLoading(false);
      toast.success("Profile saved successfully!");
    }, 1000);
  };

  const handleUpgrade = () => {
    if (!invitationCode) {
      toast.error("Please enter an invitation code.");
      return;
    }

    if (invitationCode !== "FITFLOW2026") {
      toast.error("Invalid invitation code.");
      return;
    }

    setUpgrading(true);
    // Mock API call to upgrade role
    setTimeout(() => {
      if (token) {
        login(token, "manager");
        toast.success("Account upgraded to Manager successfully!");
        setInvitationCode("");
      }
      setUpgrading(false);
    }, 1000);
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
        </div>
        <Button
          onClick={handleSave}
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
        {/* Left Column - Avatar & Basic Info */}
        <div className="space-y-6">
          <Card className="p-6 text-center space-y-4">
            <div className="relative w-32 h-32 mx-auto">
              <img
                src={`https://picsum.photos/id/${1}/200/200`}
                alt="Profile"
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

          {/* Upgrade Section for Students */}
          {role === "student" && (
            <Card className="p-6 space-y-4 border-indigo-100 bg-indigo-50/30">
              <h3 className="text-lg font-bold text-slate-800 flex items-center gap-2">
                <span className="w-1 h-6 bg-indigo-500 rounded-full" />
                Upgrade Account
              </h3>
              <p className="text-sm text-slate-500">
                Become a manager to access advanced features.
              </p>
              <div className="space-y-3">
                <div className="relative group">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-400 group-focus-within:text-indigo-500 transition-colors">
                    <Icons.Key className="w-4 h-4" />
                  </div>
                  <Input
                    placeholder="Invitation Code"
                    value={invitationCode}
                    onChange={(e) => setInvitationCode(e.target.value)}
                    className="pl-9 py-2.5 text-sm"
                  />
                </div>
                <Button
                  variant="outline"
                  onClick={handleUpgrade}
                  className="w-full py-2.5 text-sm"
                  disabled={upgrading}
                >
                  {upgrading ? (
                    <span className="flex items-center gap-2">
                      <span className="w-3 h-3 border-2 border-indigo-600/20 border-t-indigo-600 rounded-full animate-spin" />
                      Upgrading...
                    </span>
                  ) : (
                    "Upgrade to Manager"
                  )}
                </Button>
              </div>
            </Card>
          )}
        </div>

        {/* Right Column - Form Fields */}
        <div className="md:col-span-2 space-y-6">
          {/* Account Information */}
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
              </div>
            </div>
          </Card>

          {/* Personal Details */}
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
                  <option value="">Select Gender</option>
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
                  placeholder="+1 (555) 000-0000"
                />
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
    </div>
  );
};

export default Profile;
