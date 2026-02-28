// @strict: true
// @target: esnext
// @noEmit: true
// @jsx: preserve

// Repro: JSX parsing fails with spurious syntax errors in ternary expressions
// when the truthy branch has a parenthesized identifier on a separate line,
// and the falsy branch contains JSX with nested JSX attribute values containing
// function calls with multiple arguments including multi-property object literals.

/// <reference path="/.lib/react.d.ts" />

import * as React from 'react';

declare function t(key: string, params?: Record<string, any>): string;
declare function nf(v: any, opts: { precision: number; rounding: string }): string;

const HoverCardText = (p: { label: any; text: any; className?: string }) => null;
const DEFAULT_NULL_VALUE = '--';

export const Example = React.memo(function Example() {
  const isLogin = true;

  return (
    <div>
      <div>{t('label')}</div>
      {!isLogin ? (
        DEFAULT_NULL_VALUE
      ) : (
        <HoverCardText
          className="test"
          label={
            <div>
              {t('some.key', {
                s1: nf(1, { precision: 2, rounding: 'down' }),
                s2: nf(2, { precision: 2, rounding: 'down' }),
              })}
            </div>
          }
          text={
            <div>
              {nf(0, { precision: 2, rounding: 'down' })} USDT
            </div>
          }
        ></HoverCardText>
      )}
    </div>
  );
});
