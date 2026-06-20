#!/usr/bin/env bash
#
# テンプレート初期化スクリプト
#
# このテンプレートから生成されたリポジトリで（template-cleanup.yml ワークフロー
# 経由で）一度だけ実行され、テンプレートのプレースホルダ名を新しいリポジトリに
# 合わせて書き換える:
#
#   github.com/HMasataka/beyond  ->  github.com/<owner>/<repo>   (Go モジュールパス)
#   beyond-db                    ->  <repo>-db                   (データベース名)
#   beyond / Beyond              ->  <repo> / <Repo>             (プロジェクト名)
#
# owner 名 "HMasataka" 単体は意図的に置換しない。これにより oapi-generator の
# git submodule (github.com:HMasataka/oapi-generator.git) など、実在する上流への
# 参照を壊さないようにしている。
#
set -euo pipefail

: "${GITHUB_REPOSITORY:?GITHUB_REPOSITORY は必須}"        # owner/repo
: "${GITHUB_REPOSITORY_OWNER:?GITHUB_REPOSITORY_OWNER は必須}"

FULL="${GITHUB_REPOSITORY}"
REPO="${FULL#*/}"

# 小文字（DB 名 / プロジェクト ID 用）と Title-case（README 見出し用）の派生を作る。
REPO_LOWER="$(printf '%s' "$REPO" | tr '[:upper:]' '[:lower:]')"
REPO_TITLE="$(printf '%s' "${REPO_LOWER:0:1}" | tr '[:lower:]' '[:upper:]')${REPO_LOWER:1}"

OLD_MODULE="github.com/HMasataka/beyond"
NEW_MODULE="github.com/${FULL}"

echo "テンプレート初期化:"
echo "  module: ${OLD_MODULE} -> ${NEW_MODULE}"
echo "  name:   beyond -> ${REPO_LOWER} (Beyond -> ${REPO_TITLE})"

# プレースホルダを含む、追跡済みかつ非バイナリのテキストファイルを収集する。
# 初期化スクリプト自身は最後に削除するため除外する。
mapfile -t FILES < <(
  git ls-files \
    | grep -vE '^\.github/(template-cleanup/|workflows/template-cleanup\.yml$)' \
    | while IFS= read -r f; do
        if grep -EIiq 'beyond' "$f"; then printf '%s\n' "$f"; fi
      done
)

export OLD_MODULE NEW_MODULE REPO_LOWER REPO_TITLE

for f in "${FILES[@]}"; do
  # 置換順が重要: 先に完全なモジュールパスを書き換え、末尾の "beyond" を消してから
  # 単体名の置換を行う。"beyond-db" は "beyond" より先に処理する。
  # \b（単語境界）と \Q..\E（リテラル化）を環境差なく扱うため sed ではなく perl を使う。
  perl -i -pe '
    s|\Q$ENV{OLD_MODULE}\E|$ENV{NEW_MODULE}|g;
    s|beyond-db|$ENV{REPO_LOWER}-db|g;
    s|\bBeyond\b|$ENV{REPO_TITLE}|g;
    s|\bbeyond\b|$ENV{REPO_LOWER}|g;
  ' "$f"
  echo "  updated: $f"
done

# 再実行されないように初期化用の仕組みを削除する。
rm -rf .github/template-cleanup
rm -f .github/workflows/template-cleanup.yml

echo "テンプレート初期化が完了しました。"
