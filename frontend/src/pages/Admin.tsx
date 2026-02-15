import { useState, useEffect, type ChangeEvent } from "react";
import { Link } from "react-router-dom";
import { motion, AnimatePresence } from "framer-motion";
import ArrowBackRoundedIcon from "@mui/icons-material/ArrowBackRounded";
import DeleteOutlineRoundedIcon from "@mui/icons-material/DeleteOutlineRounded";
import AddRoundedIcon from "@mui/icons-material/AddRounded";
import InventoryOutlinedIcon from "@mui/icons-material/InventoryOutlined";
import CircularProgress from "@mui/material/CircularProgress";
import { useGSService } from "../hooks/useGSService";
import { ValidationError, NetworkError } from "../types/errors";

export default function Admin() {
  const gsService = useGSService();

  const [packs, setPacks] = useState<number[]>([]);
  const [loading, setLoading] = useState(true);
  const [adding, setAdding] = useState(false);
  const [deletingSize, setDeletingSize] = useState<number | null>(null);
  const [newSize, setNewSize] = useState("");
  const [validationError, setValidationError] = useState<string | null>(null);
  const [toast, setToast] = useState<string | null>(null);

  function showToast(message: string) {
    setToast(message);
    setTimeout(() => setToast(null), 4000);
  }

  useEffect(() => {
    async function fetchPacks() {
      try {
        const data = await gsService.getPacks();
        setPacks(data.packs);
      } catch (err: unknown) {
        showToast(
          err instanceof Error ? err.message : "Failed to load packs",
        );
      } finally {
        setLoading(false);
      }
    }

    fetchPacks();
  }, [gsService]);

  function handleNewSizeChange(e: ChangeEvent<HTMLInputElement>) {
    const raw = e.target.value;
    if (raw === "" || /^\d+$/.test(raw)) {
      setNewSize(raw);
      setValidationError(null);
    }
  }

  async function handleAdd() {
    const size = parseInt(newSize, 10);

    if (!newSize || isNaN(size) || size <= 0) {
      setValidationError("Enter a positive number");
      return;
    }

    if (packs.includes(size)) {
      setValidationError(`Pack size ${size} already exists`);
      return;
    }

    setValidationError(null);

    try {
      setAdding(true);
      const data = await gsService.addPack(size);
      setPacks(data.packs);
      setNewSize("");
    } catch (err) {
      if (err instanceof ValidationError) {
        setValidationError(err.message);
      } else {
        showToast(
          err instanceof NetworkError
            ? err.message
            : "An unexpected error occurred",
        );
      }
    } finally {
      setAdding(false);
    }
  }

  async function handleDelete(size: number) {
    try {
      setDeletingSize(size);
      const data = await gsService.deletePack(size);
      setPacks(data.packs);
    } catch (err) {
      showToast(
        err instanceof Error ? err.message : "Failed to delete pack",
      );
    } finally {
      setDeletingSize(null);
    }
  }

  return (
    <div className="flex min-h-screen flex-col bg-slate-950 px-4 py-8">
      <div className="mx-auto w-full max-w-lg">
        <Link
          to="/"
          className="inline-flex items-center gap-1 text-sm text-slate-400 transition-colors hover:text-indigo-400"
        >
          <ArrowBackRoundedIcon sx={{ fontSize: 18 }} />
          Back to home
        </Link>

        <h1 className="mt-6 text-3xl font-bold text-white">
          Pack Configuration
        </h1>
        <p className="mt-2 text-slate-400">
          Manage the available pack sizes for shipments.
        </p>
      </div>

      <div className="mx-auto mt-10 w-full max-w-lg">
        <label className="mb-3 block text-sm font-medium text-slate-300">
          Add new pack size
        </label>

        <div className="flex items-center gap-3">
          <input
            type="text"
            inputMode="numeric"
            placeholder="e.g. 750"
            value={newSize}
            onChange={handleNewSizeChange}
            onKeyDown={(e) => {
              if (e.key === "Enter" && !adding) handleAdd();
            }}
            className="flex-1 rounded-xl border border-slate-700 bg-slate-800/50 px-4 py-3 text-lg font-semibold text-white outline-none transition-colors placeholder:font-normal placeholder:text-slate-600 focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500/50"
          />

          <motion.button
            whileTap={{ scale: 0.95 }}
            transition={{ duration: 0.1, ease: "easeOut" }}
            disabled={adding || !newSize}
            onClick={handleAdd}
            className="flex shrink-0 cursor-pointer items-center gap-2 rounded-xl bg-indigo-600 px-5 py-3 font-semibold text-white transition-colors hover:bg-indigo-500 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {adding ? (
              <CircularProgress size={20} color="inherit" />
            ) : (
              <AddRoundedIcon sx={{ fontSize: 20 }} />
            )}
            Add
          </motion.button>
        </div>

        <AnimatePresence>
          {validationError && (
            <motion.p
              initial={{ opacity: 0, y: -4 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -4 }}
              className="mt-2 text-sm text-red-400"
            >
              {validationError}
            </motion.p>
          )}
        </AnimatePresence>
      </div>

      <div className="mx-auto mt-10 w-full max-w-lg">
        <h2 className="mb-4 flex items-center gap-2 text-lg font-semibold text-white">
          <InventoryOutlinedIcon sx={{ fontSize: 22 }} />
          Current Packs
        </h2>

        {loading ? (
          <div className="space-y-3">
            {[1, 2, 3].map((i) => (
              <div
                key={i}
                className="h-14 animate-pulse rounded-xl bg-slate-800/50"
              />
            ))}
          </div>
        ) : packs.length === 0 ? (
          <p className="text-sm text-slate-500">No packs configured.</p>
        ) : (
          <ul className="space-y-2">
            <AnimatePresence initial={false}>
              {packs.map((size) => (
                <motion.li
                  key={size}
                  layout
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  exit={{ opacity: 0, x: 20, transition: { duration: 0.2 } }}
                  transition={{ duration: 0.3 }}
                  className="flex items-center justify-between rounded-xl border border-slate-700/50 bg-slate-800/50 px-5 py-3.5"
                >
                  <span className="text-lg font-semibold text-white">
                    {size.toLocaleString()}
                    <span className="ml-2 text-sm font-normal text-slate-400">
                      items
                    </span>
                  </span>

                  <motion.button
                    whileTap={{ scale: 0.95 }}
                    transition={{ duration: 0.1, ease: "easeOut" }}
                    disabled={packs.length <= 1 || deletingSize !== null}
                    onClick={() => handleDelete(size)}
                    className="flex cursor-pointer items-center justify-center rounded-lg p-2 text-slate-400 transition-colors hover:bg-red-500/10 hover:text-red-400 disabled:cursor-not-allowed disabled:opacity-30"
                    title={
                      packs.length <= 1
                        ? "Cannot remove the last pack"
                        : `Delete pack ${size}`
                    }
                  >
                    {deletingSize === size ? (
                      <CircularProgress size={20} color="inherit" />
                    ) : (
                      <DeleteOutlineRoundedIcon sx={{ fontSize: 20 }} />
                    )}
                  </motion.button>
                </motion.li>
              ))}
            </AnimatePresence>
          </ul>
        )}
      </div>

      <AnimatePresence>
        {toast && (
          <motion.div
            initial={{ opacity: 0, y: 40 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: 40 }}
            transition={{ duration: 0.25 }}
            className="fixed bottom-6 left-1/2 z-50 -translate-x-1/2 rounded-xl border border-red-500/30 bg-red-950/90 px-5 py-3 text-sm text-red-200 shadow-lg backdrop-blur-sm"
          >
            {toast}
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
