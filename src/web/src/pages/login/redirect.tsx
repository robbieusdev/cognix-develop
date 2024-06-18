import { useLocalStorage } from "@/lib/local-store";
import { router } from "@/main";
import axios from "axios";
import { useEffect } from "react";

function RedirectComponent() {
  const { set } = useLocalStorage();

  useEffect(() => {
    axios
      .get(
        `${import.meta.env.VITE_PLATFORM_API_LOGIN_CALLBACK_URL}?${
          window.location.href.split("?")[1]
        }`
      )
      .then((response) => {
        if (response.status === 200) {
          set("access_token", response.data.data);
        }
        router.navigate("/");
        return "";
      })
  }, []);

  return <></>;
}

export { RedirectComponent as Component };
