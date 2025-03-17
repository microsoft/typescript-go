//// [tests/cases/conformance/salsa/thisPropertyAssignmentInherited.ts] ////

//// [thisPropertyAssignmentInherited.js]
export class Element {
  /**
   * @returns {String}
   */
  get textContent() {
    return  ''
  }
  set textContent(x) {}
  cloneNode() { return this}
}
export class HTMLElement extends Element {}
export class TextElement extends HTMLElement {
  get innerHTML() { return this.textContent; }
  set innerHTML(html) { this.textContent = html; }
  toString() {
  }
}



//// [thisPropertyAssignmentInherited.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.TextElement = exports.HTMLElement = exports.Element = void 0;
class Element {
    get textContent() {
        return '';
    }
    set textContent(x) { }
    cloneNode() { return this; }
}
exports.Element = Element;
class HTMLElement extends Element {
}
exports.HTMLElement = HTMLElement;
class TextElement extends HTMLElement {
    get innerHTML() { return this.textContent; }
    set innerHTML(html) { this.textContent = html; }
    toString() {
    }
}
exports.TextElement = TextElement;
