import { useState, type ChangeEvent } from "react";
import { Link } from "react-router-dom";
import { motion, AnimatePresence } from "framer-motion";
import ArrowBackRoundedIcon from "@mui/icons-material/ArrowBackRounded";
import RemoveRoundedIcon from "@mui/icons-material/RemoveRounded";
import AddRoundedIcon from "@mui/icons-material/AddRounded";
import LocalShippingOutlinedIcon from "@mui/icons-material/LocalShippingOutlined";
import InventoryOutlinedIcon from "@mui/icons-material/InventoryOutlined";
import CircularProgress from "@mui/material/CircularProgress";
import { useGSService } from "../hooks/useGSService";
import { ValidationError } from "../types/errors";
import type { ShipmentResponse } from "../types/shipment";

const MIN = 0;
const MAX = 1000000;

export default function User() {
  const gsService = useGSService();

  const [quantity, setQuantity] = useState(0);
  const [inputValue, setInputValue] = useState("0");
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<ShipmentResponse | null>(null);
  const [validationError, setValidationError] = useState<string | null>(null);
  const [toast, setToast] = useState<string | null>(null);

  function showToast(message: string) {
    setToast(message);
    setTimeout(() => setToast(null), 4000);
  }

  function updateQuantity(value: number) {
    const clamped = Math.max(MIN, Math.min(MAX, value));
    setQuantity(clamped);
    setInputValue(String(clamped));
    setValidationError(null);
  }

  function handleInputChange(e: ChangeEvent<HTMLInputElement>) {
    const raw = e.target.value;

    if (raw === "") {
      setInputValue("");
      setQuantity(0);
      setValidationError(null);
      return;
    }

    if (!/^\d+$/.test(raw)) return;

    const num = parseInt(raw, 10);
    if (num > MAX) return;

    setInputValue(raw);
    setQuantity(num);
    setValidationError(null);
  }

  function handleInputBlur() {
    setInputValue(String(quantity));
  }

  async function handleCalculate() {
    setResult(null);
    setValidationError(null);

    try {
      setLoading(true);
      const data = await gsService.calculateShipment(quantity);
      setResult(data);
    } catch (err) {
      if (err instanceof ValidationError) {
        setValidationError(err.message);
      } else {
        showToast(
          err instanceof Error ? err.message : "An unexpected error occurred",
        );
      }
    } finally {
      setLoading(false);
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
          Shipment Calculator
        </h1>
        <p className="mt-2 text-slate-400">
          Enter the number of items to calculate optimal pack sizes.
        </p>
      </div>

      <div className="mx-auto mt-10 w-full max-w-lg">
        <label className="mb-3 block text-sm font-medium text-slate-300">
          Quantity
        </label>

        <div className="flex items-center gap-3">
          <motion.button
            whileTap={{ scale: 0.95 }}
            transition={{ duration: 0.1, ease: "easeOut" }}
            disabled={quantity <= MIN}
            onClick={() => updateQuantity(quantity - 1)}
            className="flex size-12 shrink-0 cursor-pointer items-center justify-center rounded-xl border border-slate-700 bg-slate-800 text-slate-300 transition-colors hover:border-slate-600 hover:bg-slate-700 disabled:cursor-not-allowed disabled:opacity-40"
          >
            <RemoveRoundedIcon />
          </motion.button>

          <div className="relative flex-1">
            <input
              type="text"
              inputMode="numeric"
              value={inputValue}
              onChange={handleInputChange}
              onBlur={handleInputBlur}
              className="w-full rounded-xl border border-slate-700 bg-slate-800/50 px-4 py-3 text-center text-lg font-semibold text-white outline-none transition-colors focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500/50"
            />
            <AnimatePresence mode="wait">
              <motion.span
                key={quantity}
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                exit={{ opacity: 0 }}
                transition={{ duration: 0.15 }}
                className="pointer-events-none absolute inset-0"
              />
            </AnimatePresence>
          </div>

          <motion.button
            whileTap={{ scale: 0.95 }}
            transition={{ duration: 0.1, ease: "easeOut" }}
            disabled={quantity >= MAX}
            onClick={() => updateQuantity(quantity + 1)}
            className="flex size-12 shrink-0 cursor-pointer items-center justify-center rounded-xl border border-slate-700 bg-slate-800 text-slate-300 transition-colors hover:border-slate-600 hover:bg-slate-700 disabled:cursor-not-allowed disabled:opacity-40"
          >
            <AddRoundedIcon />
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

        <motion.button
          whileTap={{ scale: 0.95 }}
          transition={{ duration: 0.1, ease: "easeOut" }}
          disabled={loading || quantity <= 0}
          onClick={handleCalculate}
          className="mt-6 flex w-full cursor-pointer items-center justify-center gap-2 rounded-xl bg-indigo-600 px-6 py-3.5 font-semibold text-white transition-colors hover:bg-indigo-500 disabled:cursor-not-allowed disabled:opacity-50"
        >
          {loading ? (
            <CircularProgress size={20} color="inherit" />
          ) : (
            <LocalShippingOutlinedIcon sx={{ fontSize: 20 }} />
          )}
          {loading ? "Calculating…" : "Calculate"}
        </motion.button>
      </div>

      <AnimatePresence>
        {result && (
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: 20 }}
            transition={{ duration: 0.3 }}
            className="mx-auto mt-10 w-full max-w-lg"
          >
            <div className="rounded-2xl border border-slate-700/50 bg-slate-800/50 p-6">
              <h2 className="flex items-center gap-2 text-lg font-semibold text-white">
                <InventoryOutlinedIcon sx={{ fontSize: 22 }} />
                Packs Breakdown
              </h2>

              <div className="mt-4 space-y-2">
                {Object.entries(result.packs)
                  .sort(([a], [b]) => Number(b) - Number(a))
                  .map(([size, count]) => (
                    <div
                      key={size}
                      className="flex items-center justify-between rounded-lg bg-slate-900/50 px-4 py-3"
                    >
                      <span className="text-slate-300">
                        Pack of{" "}
                        <span className="font-semibold text-white">
                          {Number(size).toLocaleString()}
                        </span>
                      </span>
                      <span className="rounded-lg bg-indigo-500/10 px-3 py-1 text-sm font-semibold text-indigo-400">
                        ×{count}
                      </span>
                    </div>
                  ))}
              </div>

              <div className="mt-4 flex items-center justify-between border-t border-slate-700/50 pt-4">
                <span className="text-slate-400">Total items shipped</span>
                <span className="text-xl font-bold text-white">
                  {result.items_shipped.toLocaleString()}
                </span>
              </div>
            </div>
          </motion.div>
        )}
      </AnimatePresence>

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
