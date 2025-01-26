{
  description = "1Password Vault Plugin";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.11";
    gomod2nix = {
      url = "github:tweag/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    flake-utils.url = "github:numtide/flake-utils";
    # Used by shell.nix as a compat shim.
    flake-compat = {
      url = "github:edolstra/flake-compat";
      flake = false;
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      flake-compat,
      gomod2nix,
      ...
    }:
    let
      vault-plugin-secrets-onepassword =
        pkgs:
        pkgs.buildGoModule rec {
          pname = "vault-plugin-secrets-onepassword";
          version = "1.1.0";
          vendorHash = null;

          src = ./.;
          modules = ./gomod2nix.toml;
          nativeBuildInputs = pkgs.lib.optionals pkgs.stdenv.isLinux [ pkgs.makeWrapper ];
          doCheck = false;

          buildPhase = with pkgs; ''
            go build .
          '';
          installPhase = with pkgs; ''
            mkdir -p $out/bin
            go build  -o $out/bin/op-connect .
          '';

          postInstall = pkgs.lib.optionalString pkgs.stdenv.isLinux ''
          '';
        };

      flakeForSystem =
        nixpkgs: system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
          vp = vault-plugin-secrets-onepassword pkgs;
        in
        {
          packages = {
            default = vp;
            vault-plugin-secrets-onepassword = vp;
          };

          devShell = pkgs.mkShell {
            packages = with pkgs; [
              # system tools
              automake
              curl
              which
              act
              gcc
              ruby
              git
              sqlite-interactive

              pre-commit

              # go tools
              go
              gopls
              gotools
              go-tools
              gomod2nix.packages.${system}.default
              delve
              golangci-lint
              goreleaser
            ];
          };
        };
    in
    flake-utils.lib.eachDefaultSystem (system: flakeForSystem nixpkgs system);
}
