// @strict: true
// @noEmit: true

// Test case for recursion guard in destructuring pattern type checking
// This pattern has potential for infinite recursion in getSyntheticElementAccess/getParentElementAccess
const { c, f }: string | number | symbol = { c: 0, f };
