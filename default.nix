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
  goPackages = pkgs.go113Packages;
in
{
  python = setup {
    inherit pkgs pythonPackages;
    src = ./.;
  };
  go = pkgs.buildGo113Module {
    pname = "stred-proto";
    version = "0.0.0";
    # will change when `go.mod` changes - trust `nix` and update this with
    # whatever it comes up with
    modSha256 = "0sjjj9z1dhilhpc8pq4154czrb79z9cm044jvn75kxcjv6v5l2m5";
    src = ./.;
  };
}
