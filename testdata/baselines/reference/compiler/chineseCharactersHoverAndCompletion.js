//// [tests/cases/compiler/chineseCharactersHoverAndCompletion.ts] ////

//// [chineseCharactersHoverAndCompletion.ts]
// Test Chinese characters in hover and completion display
interface 中文界面 {
    上居中: string;
    下居中: string; 
    右居中: string;
    左居中: string;
}

class 中文类 {
    setLengthTextPositionPreset(preset: "上居中" | "下居中" | "右居中" | "左居中"): void {
        console.log("设置位置: " + preset);
    }
    
    获取中文属性(): 中文界面 {
        return {
            上居中: "上居中",
            下居中: "下居中",
            右居中: "右居中",
            左居中: "左居中"
        };
    }
}

let 实例 = new 中文类();
let 属性对象 = 实例.获取中文属性();
实例.setLengthTextPositionPreset("上居中");

//// [chineseCharactersHoverAndCompletion.js]
class 中文类 {
    setLengthTextPositionPreset(preset) {
        console.log("设置位置: " + preset);
    }
    获取中文属性() {
        return {
            上居中: "上居中",
            下居中: "下居中",
            右居中: "右居中",
            左居中: "左居中"
        };
    }
}
let 实例 = new 中文类();
let 属性对象 = 实例.获取中文属性();
实例.setLengthTextPositionPreset("上居中");
