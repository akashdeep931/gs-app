import { useContext } from "react";
import { GSService } from "../services/GSService";
import { GSServiceContext } from "../context/GSServiceContext";

export function useGSService(): GSService {
  const ctx = useContext(GSServiceContext);

  if (!ctx) {
    throw new Error("useGSService must be used within a GSServiceProvider");
  }

  return ctx;
}
