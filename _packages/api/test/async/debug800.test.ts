import { API } from "@typescript/api/async";
import { createVirtualFileSystem } from "@typescript/api/fs";
import { formatSyntaxKind, SyntaxKind } from "@typescript/ast";
import { readFileSync } from "node:fs";
import { resolve } from "node:path";

const repoRoot = resolve(import.meta.dirname!, "..", "..", "..", "..");
const testFile = resolve(repoRoot, "_submodules/TypeScript/src/services/mapCode.ts");
const fileText = readFileSync(testFile, "utf-8");

(async () => {
  const api = new API({
    cwd: repoRoot,
    tsserverPath: resolve(repoRoot, `built/local/tsgo${process.platform === "win32" ? ".exe" : ""}`),
    fs: createVirtualFileSystem({
      "/tsconfig.json": JSON.stringify({ files: ["/src/testFile.ts"] }),
      "/src/testFile.ts": fileText,
    }),
  });
  const snapshot = await api.updateSnapshot({ openProject: "/tsconfig.json" });
  const project = snapshot.getProject("/tsconfig.json")!;
  const sf = await project.program.getSourceFile("/src/testFile.ts");
  if (!sf) { console.log("no sf"); await api.close(); return; }

  // Find the stmt containing position 800
  let stmt: any = undefined;
  sf.forEachChild((n: any) => {
    if (n.pos <= 800 && n.end > 800) { stmt = n; }
    return undefined;
  });
  
  console.log("stmt:", stmt ? formatSyntaxKind(stmt.kind) : "none", `[${stmt?.pos}, ${stmt?.end})`);
  
  if (stmt) {
    console.log("stmt.jsDoc:", stmt.jsDoc ? stmt.jsDoc.length : "none");
    if (stmt.jsDoc) {
      for (const jd of stmt.jsDoc) {
        console.log("  jsDoc:", formatSyntaxKind((jd as any).kind), `[${(jd as any).pos}, ${(jd as any).end})`);
      }
    }
    
    console.log("\nstmt.forEachChild (all children):");
    stmt.forEachChild((node: any) => {
      console.log(`  node: ${formatSyntaxKind(node.kind)} [${node.pos},${node.end})`);
      return undefined;
    }, (ns: any) => {
      for (const n of ns) {
        console.log(`  node-arr: ${formatSyntaxKind(n.kind)} [${n.pos},${n.end})`);
      }
      return undefined;
    });
    
    console.log("\nChecking jsDoc of ExportKeyword (first child):");
    let firstMod: any = undefined;
    stmt.forEachChild((n: any) => n, (ns: any) => { firstMod = ns[0]; return ns; });
    if (firstMod) {
      console.log("  firstMod:", formatSyntaxKind(firstMod.kind), `[${firstMod.pos}, ${firstMod.end})`);
      console.log("  firstMod.jsDoc:", (firstMod as any).jsDoc ? "present" : "undefined");
    }
  }

  await api.close();
})();
