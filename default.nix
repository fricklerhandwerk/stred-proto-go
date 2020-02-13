let
  setup = import (builtins.fetchGit {
    url = "https://github.com/nix-community/setup.nix";
    ref = "v3.3.0";
  });
  nixpkgs = builtins.fetchGit {
    name = "nixpkgs";
    url = "https://github.com/NixOS/nixpkgs-channels";
    ref = "nixos-19.09";
    rev = "b9cb3b2fb2f45ac8f3a8f670c90739eb34207b0e";
  };
  pkgs = import (nixpkgs) {};
  lib = import "${nixpkgs}/lib";
  pythonPackages = pkgs.python3Packages;
in
setup {
  inherit pkgs pythonPackages;
  src = ./.;
}
