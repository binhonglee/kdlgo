# kdlgo

Go parser for the [KDL Document Language](https://github.com/kdl-org/kdl).

To reach 1.0 parity


- [x] String
- [x] Raw String
  - [ ] Recreating raw strings as regular string (according to spec)
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
- [ ] Inline `=` node
- [ ] Type Annotations
- [ ] Quoted node
- [ ] Parse from string