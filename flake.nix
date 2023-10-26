{
  description = "Chronicle project Go 1.21.0, TailwindCSS, and Flyctl";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
    flake-utils.url = "github:numtide/flake-utils";
	  templ.url = "github:a-h/templ";
  };

  outputs = { self, nixpkgs, flake-utils, templ, ... }:
    flake-utils.lib.eachDefaultSystem (system: 
      let
        pkgs = nixpkgs.legacyPackages.${system};
        zsh = pkgs.zsh;
        go = pkgs.go_1_21;
        tailwindcss = pkgs.tailwindcss;
        flyctl = pkgs.flyctl;
		    templBinary = templ.packages.${system}.default;
      in

      {
        devShell = pkgs.mkShell {
          buildInputs = [ go tailwindcss flyctl templBinary zsh];
          SHELL = "${pkgs.zsh}/bin/zsh";
          shellHook = ''
            if [ -z "$IN_NIX_SHELL_ZSH_STARTED" ]; then
              export IN_NIX_SHELL_ZSH_STARTED=1
              exec $SHELL
            fi
          '';
        };
      }
    );
}

