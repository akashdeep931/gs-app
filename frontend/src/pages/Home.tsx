import { useNavigate } from "react-router-dom";
import { motion } from "framer-motion";
import PersonOutlineRoundedIcon from "@mui/icons-material/PersonOutlineRounded";
import AdminPanelSettingsOutlinedIcon from "@mui/icons-material/AdminPanelSettingsOutlined";

const roles = [
  {
    label: "User",
    path: "/user",
    description: "Browse and order packs",
    icon: <PersonOutlineRoundedIcon sx={{ fontSize: 32 }} />,
  },
  {
    label: "Admin",
    path: "/admin",
    description: "Manage pack configurations",
    icon: <AdminPanelSettingsOutlinedIcon sx={{ fontSize: 32 }} />,
  },
];

export default function Home() {
  const navigate = useNavigate();

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-linear-to-br from-slate-950 via-slate-900 to-slate-950 px-4">
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6 }}
        className="mb-16 text-center"
      >
        <h1 className="text-5xl font-bold tracking-tight text-white sm:text-6xl">
          GS App
        </h1>
        <p className="mt-4 text-lg text-slate-400">
          Select your role to continue
        </p>
      </motion.div>

      <div className="flex flex-col gap-6 sm:flex-row">
        {roles.map((role, i) => (
          <motion.button
            key={role.path}
            onClick={() => navigate(role.path)}
            initial={{ opacity: 0, y: 30 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.2 + i * 0.15 }}
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.97 }}
            className="group flex w-64 cursor-pointer flex-col items-center gap-4 rounded-2xl border border-slate-700/50 bg-slate-800/50 px-8 py-10 backdrop-blur-sm transition-colors hover:border-indigo-500/60 hover:bg-slate-800"
          >
            <div className="rounded-xl bg-indigo-500/10 p-4 text-indigo-400 transition-colors group-hover:bg-indigo-500/20 group-hover:text-indigo-300">
              {role.icon}
            </div>
            <span className="text-xl font-semibold text-white">
              {role.label}
            </span>
            <span className="text-sm text-slate-400">{role.description}</span>
          </motion.button>
        ))}
      </div>
    </div>
  );
}
