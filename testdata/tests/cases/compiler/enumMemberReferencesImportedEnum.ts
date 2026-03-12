// @Filename: imported.ts
export enum Imported {
	Option1 = 1,
	Option2 = "hello",
}

// @Filename: usage.ts
import { Imported } from "./imported.js";
export enum Usage {
	Option1 = Imported.Option1,
	Option2 = 2,
	Option3 = Imported.Option2,
}
