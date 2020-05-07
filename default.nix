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
  pname = "stred-proto-go";
  version = "0.1.0";
  # will change when `go.mod` changes - trust `nix` and update this with
  # whatever it comes up with
  modSha256 = "04xgy6gd4zgld74vdys9a787c4chcil4zw8qmnfijxjy3r0m9lba";
  src = ./.;
  shellHook = ''
    # work around <https://github.com/NixOS/nixpkgs/issues/69401>
    unset GOPATH
  '';
}
