{
  description = "Development environment for infra-mz toolkit";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };
  outputs = { self, nixpkgs, utils }:
    utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
                inherit system;
                config.allowUnfree = true;
          };
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go_1_25
            gopls
            go-tools
		gemini-cli
		claude-code
          ];
          shellHook = ''
export TERM=xterm-256color
            export COLORTERM=truecolor
            echo "🎶 Mozart Symphony: Infra-mz Toolkit Environment"
            go version
          '';
        };
      });
}
