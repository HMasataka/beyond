// ADR 0005 の web 依存ルールを機械強制する。詳細は docs/adr/0008 を参照。
/** @type {import('dependency-cruiser').IConfiguration} */
module.exports = {
  forbidden: [
    {
      name: "no-circular",
      severity: "error",
      comment: "循環参照は層構造を壊す（ADR 0005）。",
      from: {},
      to: { circular: true },
    },
    {
      name: "components-to-page-or-layout",
      severity: "error",
      comment:
        "共通コンポーネントはページ固有・レイアウトを import しない（ADR 0005）。",
      from: { path: "^src/components/" },
      to: { path: "^src/(pages|layouts)/" },
    },
    {
      name: "layout-to-page",
      severity: "error",
      comment:
        "レイアウトは router の結節点経由でページを包む。直接 import しない（ADR 0005）。",
      from: { path: "^src/layouts/" },
      to: { path: "^src/pages/" },
    },
    {
      name: "api-no-upward",
      severity: "error",
      comment:
        "api は最下層の共有インフラ。上位層を import しない（ADR 0005）。",
      from: { path: "^src/api/" },
      to: { path: "^src/(components|layouts|pages)/|^src/router\\." },
    },
    {
      name: "no-cross-page",
      severity: "error",
      comment:
        "ページ間の直接 import は禁止。共有は共通へ昇格する（ADR 0005）。",
      // $1 は from.path の capture group（ページ名）に展開される。
      // pathNot で同一ページ配下を除外し、別ページへの import だけを禁止する。
      from: { path: "^src/pages/([^/]+)/" },
      to: { path: "^src/pages/[^/]+/", pathNot: "^src/pages/$1/" },
    },
    // ^src/pages/[^/]+ はページ 1 階層前提（ADR 0005）。共通とページ固有、両方の粒度サブ層に同じ向きを適用する。
    {
      name: "atoms-no-upward",
      severity: "error",
      comment:
        "粒度の依存は下位から上位への一方向。atoms は molecules/organisms を import しない（ADR 0005）。",
      from: { path: "(^src/components|^src/pages/[^/]+)/atoms/" },
      to: { path: "(^src/components|^src/pages/[^/]+)/(molecules|organisms)/" },
    },
    {
      name: "molecules-no-upward",
      severity: "error",
      comment:
        "粒度の依存は下位から上位への一方向。molecules は organisms を import しない（ADR 0005）。",
      from: { path: "(^src/components|^src/pages/[^/]+)/molecules/" },
      to: { path: "(^src/components|^src/pages/[^/]+)/organisms/" },
    },
  ],
  options: {
    tsConfig: { fileName: "tsconfig.json" },
    tsPreCompilationDeps: true,
    moduleSystems: ["es6", "cjs"],
    includeOnly: "^src/",
    doNotFollow: { path: "node_modules" },
    exclude: {
      path: [
        "node_modules",
        "dist",
        "^src/api/schema\\.ts$",
        "\\.test\\.(ts|tsx)$",
        "^src/test/",
      ],
    },
  },
};
