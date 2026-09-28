{
  description = "Miku Researcher - High-speed, token-efficient autonomous research MCP server in pure Go";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "miku-researcher";
          version = "1.0.0";
          src = ./.;

          vendorHash = null; # Pure Go standard library, no external go.mod dependencies!

          env.CGO_ENABLED = 0;
          ldflags = [ "-s" "-w" ];

          meta = with pkgs.lib; {
            description = "High-speed, token-efficient autonomous research MCP server";
            homepage = "https://github.com/surtr85/miku-researcher";
            license = licenses.mit;
            maintainers = [ ];
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools
          ];
        };
      }
    );
}
