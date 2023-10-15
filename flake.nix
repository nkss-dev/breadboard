{
  description = "breadboard's development environment";

  inputs = {
    go-env.url = "https://flakehub.com/f/GetPsyched/go-env/0.x.x.tar.gz";
    go-env.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = inputs@{ nixpkgs, go-env, ... }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
      go-env-pkgs = go-env.outputs.packages.${system};
    in
    {
      devShells.${system}.default = pkgs.mkShell {
        buildInputs = [
          pkgs.sqlc
          go-env-pkgs.default
          go-env-pkgs.vscode
        ];
      };
    };
}
