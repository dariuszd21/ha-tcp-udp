{ pkgs ? (
    let
      inherit (builtins) fetchTree fromJSON readFile;
      inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
    in
    import (fetchTree nixpkgs.locked) {
      overlays = [
        (import "${fetchTree gomod2nix.locked}/overlay.nix")
      ];
    }
  )
, mkGoEnv ? pkgs.mkGoEnv
, gomod2nix ? pkgs.gomod2nix
, golangci-lint ? pkgs.golangci-lint
, golangci-lint-langserver ? pkgs.golangci-lint-langserver
, gopls ? pkgs.gopls
}:

let
  goEnv = mkGoEnv { pwd = ./server; };
in
pkgs.mkShell {
  packages = [
    goEnv
    gomod2nix
    golangci-lint
    golangci-lint-langserver
    gopls
  ];
}
