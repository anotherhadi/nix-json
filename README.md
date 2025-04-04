# Nix JSON

This repository generates and publishes several large `.json` files from the Nix ecosystem, updated daily using GitHub Actions.

## 🔧 What’s Inside?

The workflow builds and publishes the following files:

- `home-manager.json` — Options from [home-manager](https://github.com/nix-community/home-manager)
- `nixos.json` — Options from [nixpkgs](https://github.com/NixOS/nixpkgs) (NixOS module system)
- `nixpkgs.json` — Package metadata from nixpkgs (including all packages)
- `darwin.json` — Options from [nix-darwin](https://github.com/LnL7/nix-darwin)
- `nur.json` — Package metadata from [NUR](https://github.com/nix-community/NUR)

These files can be used for documentation, search engines, developer tools, or exploration of the Nix ecosystem.

## 📦 How Are Files Built?

The workflow is defined in `.github/workflows/update.yml` and runs every night at **00:00 UTC**. It:

1. Installs Nix and Go.
2. Runs custom Go scripts in `.github/go/` to extract and generate JSON.
3. Builds `home-manager/options.json` using a Nix flake.
4. Publishes the files as release assets to the [`latest`](https://github.com/anotherhadi/nix-json/releases/latest) release.

## 🔗 Download Files

The latest files are available at:

- `https://github.com/anotherhadi/nix-json/releases/latest/download/<filename>.json`

For example:

- [nixpkgs.json](https://github.com/anotherhadi/nix-json/releases/latest/download/nixpkgs.json)
- [home-manager.json](https://github.com/anotherhadi/nix-json/releases/latest/download/home-manager.json)
- [darwin.json](https://github.com/anotherhadi/nix-json/releases/latest/download/darwin.json)

---

Created by [@anotherhadi](https://github.com/anotherhadi) · Powered by [Nix](https://nixos.org) and [Github](https://github.com)
