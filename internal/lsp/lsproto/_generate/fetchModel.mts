import fs from "node:fs";
import path from "node:path";
import url from "node:url";

const __filename = url.fileURLToPath(new URL(import.meta.url));
const __dirname = path.dirname(__filename);

const metaModelPath = path.join(__dirname, "metaModel.json");
const metaModelSchemaPath = path.join(__dirname, "metaModelSchema.mts");

const hash = "66a087310eea0d60495ba3578d78f70409c403d9";

const metaModelURL = `https://raw.githubusercontent.com/microsoft/vscode-languageserver-node/${hash}/protocol/metaModel.json`;
const metaModelSchemaURL = `https://raw.githubusercontent.com/microsoft/vscode-languageserver-node/${hash}/tools/src/metaModel.ts`;

const metaModelResponse = await fetch(metaModelURL);
let metaModel = await metaModelResponse.text();
metaModel = metaModel.replaceAll('"_InitializeParams"', '"InitializeParamsBase"');
fs.writeFileSync(metaModelPath, metaModel);

const metaModelSchemaResponse = await fetch(metaModelSchemaURL);
const metaModelSchema = await metaModelSchemaResponse.text();
fs.writeFileSync(metaModelSchemaPath, metaModelSchema);
