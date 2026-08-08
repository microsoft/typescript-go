// @noEmit: true
// @noTypesAndSymbols: true
// @strict: true
// @target: esnext

type BodyStyle =
    | "sedan" | "coupe" | "hatchback" | "suv" | "wagon" | "pickup"
    | "minivan" | "crossover" | "limousine" | "roadster" | "convertible";
type SportBodyStyle = "hatchback" | "roadster" | "minivan" | "crossover" | "wagon" | "pickup" | "convertible";

type ModelId = string;
type Injector = "a" | "b";
type Intercooler = { type: string; attribute: string };
type Transmission = { savesOn: boolean; filter: unknown };
type Grille = { range?: unknown } | { empty: true } | { unfiltered: true };
type Gasket<T extends SportBodyStyle = SportBodyStyle> = Record<Injector, { style?: T }>;
type PaintScheme<T extends BodyStyle = BodyStyle> = [T] extends ["pickup"]
    ? { attributes: Intercooler[] }
    : Intercooler;
type Hubcap<T extends BodyStyle = BodyStyle> = [T] extends [SportBodyStyle] ? Gasket<T> : { type: ModelId };
type Trim<T extends SportBodyStyle = SportBodyStyle> = T extends unknown ? `${T}-trim` : never;
type TrimSet<T extends BodyStyle = BodyStyle> = [T] extends [SportBodyStyle]
    ? { [K in Trim<T>]?: string }
    : { [K in Trim]?: string };
type TrimEnabled<T extends BodyStyle = BodyStyle> = [T] extends [SportBodyStyle]
    ? { [K in Trim<T>]?: { enabled: boolean } }
    : { [K in Trim]?: { enabled: boolean } };
type Hood<T extends BodyStyle = BodyStyle> = T extends "sedan" | "suv" | "wagon" | "pickup" | "hatchback" | "minivan" | "crossover"
    ? { style: T }
    : never;

interface CarSpec<T extends BodyStyle = BodyStyle> {
    paintScheme?: PaintScheme<T>;
    sortOrder?: Hubcap<T>;
    modelYear?: { groupingOn?: boolean; filter?: Grille };
    saves?: { savesOn: boolean } | Transmission;
    types?: T extends SportBodyStyle ? ModelId[] : never;
    bodyStyle?: T;
    sharedConfig?: { modelYear?: unknown };
    trim?: TrimSet<T>;
    trimEnabled?: TrimEnabled<T>;
    modelOptions?: Hood<T>;
}

type DealerTrim = any;
type DealerTrims = Record<string, DealerTrim>;
type WithRequired<O, K extends keyof O> = Required<Pick<O, K>> & O;

interface Car<T extends BodyStyle = BodyStyle> extends CarSpec<T> {}
declare class Car {
    isRoadster(): this is Car<"roadster">;
    isConvertible(): this is Car<"convertible">;
    hasModelYear(): this is WithRequired<Car, "modelYear"> & {
        modelYear: { groupingOn: true };
    };
    hasPickupPaint(): this is WithRequired<Car<"pickup">, "paintScheme"> & {
        paintScheme: { attributes: Intercooler[] };
    };
    hasPickupSetting(): this is (WithRequired<Car<"pickup">, "modelYear"> & {
        modelYear: { groupingOn: true };
    }) | (WithRequired<Car<"pickup">, "saves"> & {
        saves: { savesOn: true };
    });
    hasCustomPaint(): this is WithRequired<Car<Exclude<BodyStyle, "pickup">>, "paintScheme">;
    getTrims(
        trims: DealerTrims,
        trim?: TrimSet,
        trimEnabled?: TrimEnabled,
        filter?: (key: Trim, definition: DealerTrim) => boolean,
    ): Record<Trim, DealerTrim>;
    hasTransmission(): this is Car & { saves: Transmission };
    hasFilter(): this is Car & { modelYear: { filter: Grille } };
    hasSort(type: Injector): this is Car & { sortOrder: Gasket };
    hasSportSort(): this is Car<SportBodyStyle> & { sortOrder: Gasket };
}

declare class CarInspector extends Car {}

export function repro(car: CarInspector): boolean {
    if (car.isRoadster() || car.isConvertible()) {}
    return car.hasCustomPaint();
}
