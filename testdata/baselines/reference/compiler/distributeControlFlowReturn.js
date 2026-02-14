//// [tests/cases/compiler/distributeControlFlowReturn.ts] ////

//// [distributeControlFlowReturn.ts]
type Layer = 
  | { labels: { selectedFieldName: "A" | "B" } }
  | { labels: { selectedFieldName: "C" } }

declare let layers: Layer[]

let newLayers: Layer[] = layers.map(layer => {
  return { labels: layer.labels }
})

let newLayersManualFix: Layer[] = layers.map(layer => {
  switch (layer.labels.selectedFieldName) {
    case "A": return { labels: layer.labels }
    case "B": return { labels: layer.labels }
    case "C": return { labels: layer.labels }
  }
})

let newLayersDistributeFix: Layer[] = layers.map(layer => {
  distribute (layer) {
    return { labels: layer.labels }
  }
})


//// [distributeControlFlowReturn.js]
let newLayers = layers.map(layer => {
    return { labels: layer.labels };
});
let newLayersManualFix = layers.map(layer => {
    switch (layer.labels.selectedFieldName) {
        case "A": return { labels: layer.labels };
        case "B": return { labels: layer.labels };
        case "C": return { labels: layer.labels };
    }
});
let newLayersDistributeFix = layers.map(layer => {
    {
        return { labels: layer.labels };
    }
});
