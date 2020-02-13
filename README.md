# structure editor for protobuf

this is an experiment to create a terminal-based structure editor for [proto3]. the goal is to overcome the limitations of computer languages and textual representation, to present a user interface that makes invalid input impossible, and try to create a modal keyboard control mechanism that is effective and convenient.

[proto3]: https://developers.google.com/protocol-buffers/docs/proto3

# updating requirements

nix-shell -A pip2nix --run "pip2nix generate -r requirements.txt --output=requirements.nix"
