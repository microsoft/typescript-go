import * as fs from "node:fs";
import * as path from "node:path";

const compilerFolder = "./internal/compiler";

/** @type {Record<string, {package: string, members?: boolean, rename?: string}>} */
const refactorMapping = {
    "NodeId": {
        package: "ast",
        members: false,
    },
    "MergeId": {
        package: "ast",
        members: false,
    },
    "SymbolId": {
        package: "ast",
        members: false,
    },
    "Symbol": {
        package: "ast",
        members: false,
    },
    "SymbolTable": {
        package: "ast",
        members: false,
    },
    /*"ScriptTarget": {
        package: "core",
    },
    "LanguageVariant": {
        package: "core",
    },
    "ScriptKind": {
        package: "core",
    },
    "SyntaxKind": {
        package: "ast",
        rename: "Kind"
    },
    "NodeFlags": {
        package: "ast",
    },
    "SymbolFlags": {
        package: "ast",
    },
    "CheckFlags": {
        package: "ast",
    },*/
    // "Tristate": {
    //     package: "tristate",
    // } // wasn't ported following the same pattern
}

const astNames = []
for (const n of astNames) {
    refactorMapping[n] = { package: "ast", members: false };
}

const enums = Object.keys(refactorMapping);

const entries = fs.readdirSync(compilerFolder, {recursive: true});

outer: for (const entry of entries) {
    if (!entry.endsWith(".go")) continue;

    const localPath = path.join(compilerFolder, entry);
    console.log(localPath)
    let file = fs.readFileSync(localPath, {encoding: "utf-8"});

    for (const e of enums) {
        const newRootName = refactorMapping[e].rename || e;
        // replace bare references to the type with `package.type`
        file = file.replaceAll(new RegExp(`(\\W)${e}(\\W)`, "g"), (_, prefix, postfix) => `${prefix}${refactorMapping[e].package}.${newRootName}${postfix}`);
        // Replace all member references with `package.member`
        if (refactorMapping[e].members !== false) { 
            file = file.replaceAll(new RegExp(`(\\W)${e}(\\w+)`, "g"), (_, prefix, postfix) => `${prefix}${refactorMapping[e].package}.${newRootName}${postfix}`);
        }
    }
/*
    // Do tristate bespoke like

    // replace bare references to Tristate with `Tristate.Type`
    file = file.replaceAll(/(\W)Tristate(\W)/g, (_, prefix, postfix) => prefix+"core.Tristate"+postfix);
    // And member references to `tristate.Member`
    for (const member of ["TSUnknown", "TSTrue", "TSFalse"]) {
        file = file.replaceAll(new RegExp(`(\\W)${member}(\\W)`, "g"), (_, prefix, postfix) => `${prefix}core.${member}${postfix}`);
    }
*/
    fs.writeFileSync(localPath, file);
}