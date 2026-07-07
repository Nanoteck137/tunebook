{
  description = "Devshell for tunebook";

  inputs = {
    nixpkgs.url      = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url  = "github:numtide/flake-utils";

    gitignore.url = "github:hercules-ci/gitignore.nix";
    gitignore.inputs.nixpkgs.follows = "nixpkgs";

    versionctl.url = "github:nanoteck137/versionctl/0.3.0";
  };

  outputs = { self, nixpkgs, flake-utils, gitignore, ... }@inputs:
    flake-utils.lib.eachDefaultSystem (system:
      let
        overlays = [];
        pkgs = import nixpkgs {
          inherit system overlays;
        };

        version = pkgs.lib.strings.fileContents "${self}/version";
        fullVersion = ''${version}-${self.dirtyShortRev or self.shortRev or "dirty"}'';

        backend = pkgs.buildGoModule {
          pname = "tunebook";
          version = fullVersion;
          src = ./.;
          subPackages = ["cmd/tunebook"];

          ldflags = [
            "-X github.com/nanoteck137/tunebook.Version=${version}"
            "-X github.com/nanoteck137/tunebook.Commit=${self.dirtyRev or self.rev or "no-commit"}"
          ];

          vendorHash = "sha256-62DF7oLaS+0a2sOne8twsoBu/rSNP/IrMaNVcRt10Sw=";

          nativeBuildInputs = [ pkgs.makeWrapper ];

          postFixup = ''
            wrapProgram $out/bin/tunebook --prefix PATH : ${pkgs.lib.makeBinPath [ pkgs.ffmpeg pkgs.imagemagick ]}
          '';
        };

        web = pkgs.buildNpmPackage {
          name = "tunebook-web";
          version = fullVersion;

          src = gitignore.lib.gitignoreSource ./web;
          npmDepsHash = "sha256-xIrTtTfnssHyyOF8ruqfeEdcXrzKHENy5fL06B2vZD8=";

          PUBLIC_VERSION=version;
          PUBLIC_COMMIT=self.dirtyRev or self.rev or "no-commit";
          PUBLIC_API_ADDRESS="";

          installPhase = ''
            runHook preInstall
            cp -r build $out/
            runHook postInstall
          '';
        };

        docker = pkgs.dockerTools.buildLayeredImage {
          name = "tunebook";
          tag  = fullVersion;

          contents = [
            pkgs.dockerTools.caCertificates
            web
            backend
          ];

          extraCommands = ''
            mkdir -p data
            mkdir -p library
          '';

          config = {
            Entrypoint   = [ "/bin/tunebook" ];
            Cmd = [ "serve" ];
            WorkingDir = "/data";

            ExposedPorts = { 
              "3000/tcp" = {}; 
            };

            Env = [
              "TUNEBOOK_WEB=${web}"
              "TUNEBOOK_DATA_DIR=/data"
              "TUNEBOOK_LIBRARY_DIR=/library"
            ];
          };
        };
      in
      {
        packages = {
          default = backend;
          inherit backend web docker;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            air
            go
            gopls
            nodejs
            imagemagick
            ffmpeg
            just

            inputs.versionctl.packages.${system}.default
          ];
        };
      }
    ) // {
      nixosModules.backend = import ./nix/backend.nix { inherit self; };
      nixosModules.default = { ... }: {
        imports = [
          self.nixosModules.backend
        ];
      };
    };
}
