{
  description = "Development environment for infra-mz toolkit";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };
  outputs = { self, nixpkgs, utils }:
    utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go_1_23
            gopls
            go-tools
          ];
          shellHook = ''
            echo "🎶 Mozart Symphony: Infra-mz Toolkit Environment"
            go version
          '';
        };
      });
}
