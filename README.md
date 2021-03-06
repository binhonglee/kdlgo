# kdlgo

Go parser for the [KDL Document Language](https://github.com/kdl-org/kdl).

[![Go Reference](https://pkg.go.dev/badge/github.com/binhonglee/kdlgo.svg)](https://pkg.go.dev/github.com/binhonglee/kdlgo)

To reach 1.0 parity


- [x] Parse from file
- [x] Parse from string
- [x] String
- [x] Raw String
- [x] Boolean
- [x] Null
- [x] Array Value
- [x] Nested child object
- [x] Inline `;` node
- [x] "slashdash" comment `/-`
- [x] Multi-line node `\`
- [x] Line comment `//`
- [x] Block comment `/**/`
- [x] Quoted node
- [x] Inline `=` node
- [ ] Type Annotations (Currently this will probably either get parsed as part of the node or cause an error.)
  - [ ] Ignored
  - [ ] signed int
  - [ ] unsigned int
  - [ ] float
  - [ ] decimal

- [ ] Pass the tests (I'm going through them in alphabetical order. If its not listed and its before the listed ones, its passing. If its after the listed ones, I've not looked into it.)
  - [x] empty_child_whitespace
  - [ ] empty_quoted_node_id