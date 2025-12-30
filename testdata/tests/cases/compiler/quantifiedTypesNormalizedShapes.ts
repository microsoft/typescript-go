type NormalizedRecord<T extends { id: string }> =
  <Id extends string> { [K in Id]: Omit<T, "id"> & { id: K } }

interface Layer {
  id: string
  color: string
}

let layers: NormalizedRecord<Layer> = {
  a: {
    id: "a",
    color: "green"
  },
  b: {
    id: "a", // should have been "b"
    color: "blue"
  }
}
