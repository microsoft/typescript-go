export type GoodDocNode =
    | string
    | ["i", ...GoodDocNode[]]
    | ["li", index: number, ...GoodDocNode[]];

export type NamedRestDocNode1 =
    | string
    | ["i", ...tagged: NamedRestDocNode1[]]
    | ["li", index: number, ...NamedRestDocNode1[]];

export type NamedRestDocNode2 =
    | string
    | ["i", ...NamedRestDocNode2[]]
    | ["li", index: number, ...tagged: NamedRestDocNode2[]];
