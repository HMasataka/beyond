import type { ButtonHTMLAttributes } from "react";

type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement>;

export function Button({ type = "button", ...props }: ButtonProps) {
  return (
    <button
      type={type}
      className="py-2 px-4 font-medium text-white rounded bg-brand"
      {...props}
    />
  );
}
