import axios from "axios";
import { toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import { capitalize } from "./utils";

axios.interceptors.request.use(
  (config) => {
    const accessTokenString = localStorage.getItem("access_token");
    const token = accessTokenString
      ? JSON.parse(String(localStorage.getItem("access_token")))
      : "";
    config.headers["Authorization"] = "Bearer " + token;
    config.headers["Content-Type"] = "application/json";
    config.headers["Accept"] = "application/json";
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

axios.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    const { response } = error;
    if (response && response.status && response.status >= 400) {
      toast.error(capitalize(response.data.error));
    }
    return Promise.reject(error);
  }
);

export default axios;
