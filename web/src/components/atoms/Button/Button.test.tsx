import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { Button } from "./Button";

describe("Button", () => {
  it("子要素をラベルとして表示する", () => {
    render(<Button>送信</Button>);

    expect(screen.getByRole("button", { name: "送信" })).toBeInTheDocument();
  });

  it("クリックすると onClick を呼ぶ", () => {
    const onClick = vi.fn();
    render(<Button onClick={onClick}>送信</Button>);

    fireEvent.click(screen.getByRole("button", { name: "送信" }));

    expect(onClick).toHaveBeenCalledTimes(1);
  });

  it("デフォルトの type は button", () => {
    render(<Button>送信</Button>);

    expect(screen.getByRole("button", { name: "送信" })).toHaveAttribute(
      "type",
      "button",
    );
  });
});
