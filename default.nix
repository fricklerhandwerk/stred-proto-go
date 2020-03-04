let
  nixpkgs = builtins.fetchGit {
    name = "nixpkgs";
    url = "https://github.com/NixOS/nixpkgs-channels";
    ref = "nixos-19.09";
    rev = "b9cb3b2fb2f45ac8f3a8f670c90739eb34207b0e";
  };
  pkgs = import (nixpkgs) {};
in
pkgs.buildGo113Module {
  pname = "stred-proto";
  version = "0.0.0";
  # will change when `go.mod` changes - trust `nix` and update this with
  # whatever it comes up with
  modSha256 = "07h8l5p81695x10fzryjijm42av9nsxg5hzb862n2y0n0irslx8j";
  src = ./.;
  shellHook = ''
    # work around <https://github.com/NixOS/nixpkgs/issues/69401>
    unset GOPATH
  '';
}
