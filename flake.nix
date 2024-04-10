{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-23.11";
    serve.url = "github:a-h/serve";
    xc.url = "github:joerdav/xc";
  };

  outputs = { nixpkgs, serve, xc, ... }:
    let
      pkgsForSystem = system: import nixpkgs {
        inherit system;
      };
      forAllSystems = f: {
        x86_64-linux = f "x86_64-linux";
        aarch64-linux = f "aarch64-linux";
        x86_64-darwin = f "x86_64-darwin";
        aarch64-darwin = f "aarch64-darwin";
      };
    in
    {
      devShell = forAllSystems (system:
        let
          pkgs = pkgsForSystem system;
        in
        pkgs.mkShell {
          buildInputs = [
            pkgs.openssl
            serve.outputs.packages.${system}.default
            pkgs.testssl # Use to test the server.
            xc.outputs.packages.${system}.xc
          ];
        });
    };
}
