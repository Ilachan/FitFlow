import { Routes, Route, Navigate } from "react-router-dom";
import { Toaster } from "react-hot-toast";
import Login from "./pages/Login";
import Register from "./pages/Register";
import Browse from "./pages/Browse";
import MySchedule from "./pages/MySchedule";
import Profile from "./pages/Profile";
import MainLayout from "./layouts/MainLayout";
import AuthGuard from "./components/AuthGuard";

function App() {
  return (
    <main>
      <Toaster position="top-center" reverseOrder={false} />
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />

        {/* Main Layout for public & protected routes */}
        <Route element={<MainLayout />}>
          {/* Public Routes */}
          <Route path="/courses" element={<Browse />} />
          <Route path="/" element={<Navigate to="/courses" replace />} />

          {/* Protected Routes */}
          <Route element={<AuthGuard />}>
            <Route path="/my-schedule" element={<MySchedule />} />
            <Route path="/profile" element={<Profile />} />
          </Route>
        </Route>

        {/* Fallback for unknown routes */}
        <Route path="*" element={<Navigate to="/login" replace />} />
      </Routes>
    </main>
  );
}

export default App;
