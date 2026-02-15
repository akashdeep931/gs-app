import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import { GSServiceProvider } from "./context/GSServiceProvider.tsx";
import "./index.css";
import App from "./App.tsx";

createRoot(document.getElementById("root")!).render(
  <BrowserRouter>
    <GSServiceProvider>
      <App />
    </GSServiceProvider>
  </BrowserRouter>,
);
