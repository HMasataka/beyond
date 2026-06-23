import { useQuery } from "@tanstack/react-query";
import { apiClient } from "../../api/client";

export function useHealthz() {
  return useQuery({
    queryKey: ["healthz"],
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/healthz");
      if (error) {
        throw error;
      }
      return data;
    },
  });
}
