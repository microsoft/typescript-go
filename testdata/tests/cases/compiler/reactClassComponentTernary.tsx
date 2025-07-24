// @strict: true
// @jsx: react-jsx
// @noEmit: true

/// <reference path="/.lib/react16.d.ts" />

import { Component } from "react";

type PropsTest = Readonly<{
  test: string,
}>;

class X extends Component<PropsTest, { yo: string }> {
  static defaultProps = {
    test: "x",
  };

  render() {
    return "test";
  }
}

class Y extends Component<PropsTest, { hey: string }> {
  render() {
    return "test";
  }
}

const XorY = Math.random() > 0.5 ? X : Y;

function z() {
  return <XorY test="test" />;
}
