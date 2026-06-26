// Conventional Commits を正本ルールとして commit メッセージと PR タイトルを検査する。
// 詳細な採用理由は docs/adr/0010-commit-convention-enforcement.md を参照。
export default {
  extends: ["@commitlint/config-conventional"],
  rules: {
    // 日本語の subject は語尾に「。」を打たない運用のため句点での誤検知を避ける。
    // config-conventional の subject-full-stop は "." のみ対象で「。」は素通りするが、
    // 意図を明示するためここで再宣言しておく。
    "subject-full-stop": [2, "never", "."],
  },
};
