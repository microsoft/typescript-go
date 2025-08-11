// Test file to reproduce Chinese character display issue
interface Point {
    上居中: string;
    下居中: string; 
    右居中: string;
    左居中: string;
}

class TSLine {
    setLengthTextPositionPreset(preset: "上居中" | "下居中" | "右居中" | "左居中", offsetX?: number, offsetY?: number, font_size?: number): void {
        
    }
}

let lines = new TSLine();
lines.setLengthTextPositionPreset(// cursor here