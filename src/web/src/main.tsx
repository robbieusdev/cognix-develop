import React from "react";
import ReactDOM from "react-dom/client";
import { RouterProvider, createBrowserRouter } from "react-router-dom";
import "@/global.css";
import "@/lib/axios";
import AuthProvider from "./context/AuthContext";
import { ToastContainer } from "react-toastify";
import { MessagesProvider } from "./context/ChatContext";

export const router = createBrowserRouter([
  {
    path: "/login",
    lazy: () => import("@/pages/login"),
  },
  {
    path: "/google/callback",
    lazy: () => import("@/pages/login/redirect"),
  },
  {
    path: "/",
    lazy: () => import("@/pages/platform"),
    children: [
      {
        lazy: () => import("@/pages/chat"),
        index: true,
      },
      {
        path: "/chat/:chatId",
        lazy: () => import("@/pages/chat"),
      },
      {
        path: "settings",
        children: [
          {
            path: "connectors",
            children: [
              {
                path: "existing-connectors",
                lazy: () =>
                  import("@/pages/settings/connectors/existing-connectors"),
              },
            ],
          },
          {
            path: "feedback",
            lazy: () => import("@/pages/settings/feedback"),
          },
          {
            path: "embeddings",
            lazy: () => import("@/pages/settings/embeddings"),
          },
          {
            path: "llms",
            lazy: () => import("@/pages/settings/llms"),
          },
          {
            path: "users",
            lazy: () => import("@/pages/settings/users"),
          },
          {
            path: "config",
            lazy: () => import("@/pages/settings/config"),
          },
        ],
      },
    ],
  },
]);

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <AuthProvider>
      <MessagesProvider>
        <ToastContainer position="bottom-center" />
        <RouterProvider router={router} />
      </MessagesProvider>
    </AuthProvider>
  </React.StrictMode>
);
