# structure editor for protobuf

this is an experiment to create a terminal-based structure editor for [proto3]. the goal is to overcome the limitations of computer languages and textual representation, to present a user interface that makes invalid input impossible, and try to create a modal keyboard control mechanism that is effective and convenient.

the core library is designed to be impossible to misuse, and will only create valid document models. the API guides the client implementation at compile time via the type system, and by only allowing to create objects which are tentative initially. their fields are set individually with strings and integers from user input, and validated each. only when all fields are valid, such an object can become part of the document. there is no creation of structs or passing of pointers involved at any point.

Go was chosen for this prototype because it is the one statically typed language I am fluent in. it is obviously not an optimal tool for a safe implementation, and hence much swearing is involved in the process. Python was used for a first draft, but due to its dynamic typing and lack of module isolation it is obviously not up to the task of actually making errors in API client implementation impossible. while Haskell or Idris may be much more suitable for a safe, type-driven design, I simply lack the proficiency. it looks like it is too much of a hassle to work with graphs and user interfaces. Rust seems to have a most suitable combination of strong static type system and reference handling, but no practical experience with the language is too much of a risk to build something explorative.

[proto3]: https://developers.google.com/protocol-buffers/docs/proto3

