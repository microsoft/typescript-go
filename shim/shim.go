// Package shim and its subpackages re-export the subset of this repository's
// internal API needed by external consumers of this fork (Parable's psgen
// TypeScript reader). Each subpackage is a pure veneer over the corresponding
// internal package: type aliases for types and var/const bindings for
// functions and constants. No wrapper logic lives here.
//
// The surface is intentionally minimal and is expanded only as consumer code
// demands. See the fork notice in the repository README for the policy
// governing these packages and the upstream upgrade procedure.
package shim
