{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = [
    pkgs.go
    pkgs.sqlite
    pkgs.sqlitestudio
  ];
}