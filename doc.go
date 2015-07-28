// Package fuse is a web framework with sane conventions.
// These conventions remove some of the boilerplate that naturally arises from
// using Go for web development.  Fuse automatically:
//
//     - serves files from the public directory
//     - loads templates from the template directory
//     - loads and saves sessions
//
// This means that each controller is more succinct and to the point.
package fuse
