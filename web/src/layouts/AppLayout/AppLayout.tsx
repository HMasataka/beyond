import { Outlet } from "react-router";

export function AppLayout() {
  return (
    <div className="min-h-screen">
      <header className="py-4 px-6 border-b">
        <h1 className="text-xl font-bold">beyond</h1>
      </header>
      <main className="p-6">
        <Outlet />
      </main>
    </div>
  );
}
