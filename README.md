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
  - [x] backslash_in_bare_id
  - [ ] bare_arg
  - [x] block_comment_before_node_no_space
  - [ ] block_comment_before_node
  - [x] commented_arg
  - [ ] commented_child