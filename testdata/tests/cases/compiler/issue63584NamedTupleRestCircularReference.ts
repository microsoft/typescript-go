export type GoodDocNode =
    | string
    | ["i", ...GoodDocNode[]]
    | ["li", index: number, ...GoodDocNode[]];

export type BadDocNode1 =
    | string
    | ["i", ...tagged: BadDocNode1[]]
    | ["li", index: number, ...BadDocNode1[]];

export type BadDocNode2 =
    | string
    | ["i", ...BadDocNode2[]]
    | ["li", index: number, ...tagged: BadDocNode2[]];
