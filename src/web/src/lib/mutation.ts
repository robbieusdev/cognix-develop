import axios from "axios";
import useSWRMutation from "swr/mutation";

export const useMutation = <
  const RequestData extends Record<string, any> | any[] | undefined,
  ResponseData extends Record<string, any> = Record<string, any>
>(
  path: string | null,
  method: "POST" | "PUT" | "DELETE" | "PATCH",
  headers?: Record<string, string>,
  options?: Parameters<
    typeof useSWRMutation<ResponseData, any, string, RequestData>
  >[2]
) => {
  return useSWRMutation<ResponseData, any, string | null, RequestData>(
    path,
    (key, { arg }) =>
      axios(key, { method: method, data: arg, headers: headers }),
    options
  );
};
