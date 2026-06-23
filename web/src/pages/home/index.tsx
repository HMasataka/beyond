import { Button } from "../../components/atoms/Button";
import { useHealthz } from "./api";

// 成功状態では data が型レベルで確定するため、guard 後はフォールバック不要。
function statusLabel(query: ReturnType<typeof useHealthz>): string {
  if (query.isPending) {
    return "loading";
  }
  if (query.isError) {
    return "error";
  }
  return query.data.status;
}

export function HomePage() {
  const query = useHealthz();

  return (
    <section className="flex flex-col gap-4 items-start">
      <p>API status: {statusLabel(query)}</p>
      <Button onClick={() => query.refetch()}>再取得</Button>
    </section>
  );
}
