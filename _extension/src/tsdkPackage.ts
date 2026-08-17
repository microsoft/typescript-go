import path from "node:path";

export function getNodeModulesCandidates(packagePath: string, packageName: string): string[] {
    const candidates: string[] = [];

    for (let current = path.dirname(packagePath); current !== path.dirname(current); current = path.dirname(current)) {
        if (path.basename(current) === "node_modules") {
            candidates.push(current);
            break;
        }
    }

    const metadataCandidate = packageName.startsWith("@")
        ? path.resolve(packagePath, "..", "..")
        : path.resolve(packagePath, "..");
    if (!candidates.includes(metadataCandidate)) {
        candidates.push(metadataCandidate);
    }

    return candidates;
}
