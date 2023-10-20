{
  description = "Chronicle project Go 1.21.0, TailwindCSS, and Flyctl";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system: 
      let
        pkgs = nixpkgs.legacyPackages.${system};
        go = pkgs.go_1_21;
        tailwindcss = pkgs.tailwindcss;
        flyctl = pkgs.flyctl;
      in

      {
        devShell = pkgs.mkShell {
          buildInputs = [ go tailwindcss flyctl ];
        };
      }
    );
}

