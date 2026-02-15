import type { ReactNode } from "react";
import { GSService } from "../services/GSService";
import { GSServiceContext } from "./GSServiceContext";

export function GSServiceProvider({ children }: { children: ReactNode }) {
  const service = GSService.getInstance(import.meta.env.VITE_GS_API_URL);

  return (
    <GSServiceContext.Provider value={service}>
      {children}
    </GSServiceContext.Provider>
  );
}
