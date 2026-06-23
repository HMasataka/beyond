import { createBrowserRouter } from "react-router";
import { AppLayout } from "./layouts/AppLayout";
import { HomePage } from "./pages/home";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <AppLayout />,
    children: [{ index: true, element: <HomePage /> }],
  },
]);
