import { Button } from "@/components/ui/button";
import CognixLow from "@/assets/svgs/cognix-sm.svg?react";
import { useState } from "react";
import axios from "axios";

export function LoginComponent() {
  const [error] = useState<string | null>(null);

  function login(): void {
    axios
      .get(
        `${import.meta.env.VITE_PLATFORM_API_LOGIN_URL}?redirect_url=${
          window.location.origin
        }`
      )
      .then((response) => {
        if (response.status === 200) {
          const authUrl = response.data;
          window.location.href = authUrl.data;
        }
      });
  }

  return (
    <>
      <div className="flex flex-col items-center justify-center h-screen space-y-5">
        <CognixLow className="w-20 h-20" />
        <span className="text-2xl font-bold">Log In to CogniX</span>
        <div className="flex items-center justify-center">
          <Button
            variant="outline"
            size="xl"
            className="shadow-none bg-secondary"
            type="button"
            onClick={() => login()}
          >
            Continue with Google
          </Button>
        </div>
        {error && <p className="text-red-500">{error}</p>}
      </div>
    </>
  );
}

export { LoginComponent as Component };
