# kdlgo

Go parser for the [KDL Document Language](https://github.com/kdl-org/kdl).

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
- [ ] "slashdash" comment `/-`
  - [ ] support for node
  - [x] support for value
- [x] Multi-line node `\`
- [x] Line comment `//`
- [x] Block comment `/**/`
- [x] Quoted node
- [ ] Inline `=` node
- [ ] Type Annotations (Currently this will probably either get parsed as part of the node or cause an error.)
  - [ ] Ignored
  - [ ] signed int
  - [ ] unsigned int
  - [ ] float
  - [ ] decimal