import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen, waitFor } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import { HomePage } from "./index";

const get = vi.hoisted(() => vi.fn());
vi.mock("../../api/client", () => ({ apiClient: { GET: get } }));

function renderHome() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  render(
    <QueryClientProvider client={queryClient}>
      <HomePage />
    </QueryClientProvider>,
  );
}

describe("HomePage", () => {
  afterEach(() => {
    get.mockReset();
  });

  it("取得中はローディングを表示する", () => {
    get.mockReturnValue(new Promise(() => {}));
    renderHome();

    expect(screen.getByText("API status: loading")).toBeInTheDocument();
  });

  it("成功時はステータスを表示する", async () => {
    get.mockResolvedValue({ data: { status: "ok" }, error: undefined });
    renderHome();

    await waitFor(() =>
      expect(screen.getByText("API status: ok")).toBeInTheDocument(),
    );
  });

  it("失敗時はエラーを表示する", async () => {
    get.mockResolvedValue({ data: undefined, error: { message: "boom" } });
    renderHome();

    await waitFor(() =>
      expect(screen.getByText("API status: error")).toBeInTheDocument(),
    );
  });
});
