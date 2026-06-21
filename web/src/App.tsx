import { useState } from "react";

export default function App() {
  const [count, setCount] = useState(0);

  return (
    <main>
      <h1>beyond web</h1>
      <p>React + TypeScript + Vite の最小サンプル。</p>
      <button type="button" onClick={() => setCount((c) => c + 1)}>
        count is {count}
      </button>
    </main>
  );
}
