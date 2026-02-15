import { createContext } from "react";
import { GSService } from "../services/GSService";

export const GSServiceContext = createContext<GSService | null>(null);
