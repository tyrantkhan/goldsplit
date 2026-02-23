# Changelog

## [0.0.11](https://github.com/tyrantkhan/goldsplit/compare/v0.0.10...v0.0.11) (2026-02-23)


### Bug Fixes

* correct UI labels and improve build setup ([525662f](https://github.com/tyrantkhan/goldsplit/commit/525662f894c4268baf91a116c6add9e739c3a2a0))
* correct UI labels and improve build setup ([08f0273](https://github.com/tyrantkhan/goldsplit/commit/08f027339bfd8a37d37c085b688b04d97c719488))

## [0.0.10](https://github.com/tyrantkhan/goldsplit/compare/v0.0.9...v0.0.10) (2026-02-22)


### Features

* suspend and resume in-progress runs ([1f03dbc](https://github.com/tyrantkhan/goldsplit/commit/1f03dbc544668269c5f2ac2f1b75e245f038eca0))
* suspend/resume runs, editor fixes, SoB improvements ([350ff65](https://github.com/tyrantkhan/goldsplit/commit/350ff65af2520a2d3713bf009cd1a35485a875dd))


### Bug Fixes

* edit incomplete runs and sticky header in attempt editor ([2d92ac6](https://github.com/tyrantkhan/goldsplit/commit/2d92ac6aa67c1ed60aa51da68b8cc391eae1dc7d))
* Sum of Best uses current run segments to fill gaps ([b730b41](https://github.com/tyrantkhan/goldsplit/commit/b730b416abb741d092312b75f5115714ddfaa994))


### Refactoring

* trash icon and confirm dialogs for delete actions ([415a2e2](https://github.com/tyrantkhan/goldsplit/commit/415a2e2028fc0dd5837b20f3f252cb4b83679d92))

## [0.0.9](https://github.com/tyrantkhan/goldsplit/compare/v0.0.8...v0.0.9) (2026-02-22)


### Bug Fixes

* PB fallback for comparison gaps, first-time segment marker ([d1e02df](https://github.com/tyrantkhan/goldsplit/commit/d1e02dfc0c4cd5bb21da23d127bee88708c263a2))

## [0.0.8](https://github.com/tyrantkhan/goldsplit/compare/v0.0.7...v0.0.8) (2026-02-22)


### Features

* comparison splits, timer coloring, and best-time rename ([352c415](https://github.com/tyrantkhan/goldsplit/commit/352c415523ac7bd56d103dfe5c426e1490ff8144))

## [0.0.7](https://github.com/tyrantkhan/goldsplit/compare/v0.0.6...v0.0.7) (2026-02-22)


### Bug Fixes

* handle effectively-skipped segments with same cumulative time ([24a72b4](https://github.com/tyrantkhan/goldsplit/commit/24a72b437e3729f9ad2ee2f7e8928d11614a1b0c))

## [0.0.6](https://github.com/tyrantkhan/goldsplit/compare/v0.0.5...v0.0.6) (2026-02-22)


### Refactoring

* compute comparison data on the fly from history ([3855c0c](https://github.com/tyrantkhan/goldsplit/commit/3855c0cef33aa6be5322644cefecf305761e216b))
* compute comparison data on the fly from history ([368e6ab](https://github.com/tyrantkhan/goldsplit/commit/368e6ab857c4edcb15936b7f6b6d4caf055f6fbd)), closes [#20](https://github.com/tyrantkhan/goldsplit/issues/20)

## [0.0.5](https://github.com/tyrantkhan/goldsplit/compare/v0.0.4...v0.0.5) (2026-02-22)


### Bug Fixes

* rename VERSION to version.txt so release-please updates it ([cedeaa3](https://github.com/tyrantkhan/goldsplit/commit/cedeaa3807aa6a21dc2f644bec147f41c7cf1d1d))

## [0.0.4](https://github.com/tyrantkhan/goldsplit/compare/v0.0.3...v0.0.4) (2026-02-22)


### Features

* add timer stats display, fix completion bug, and UI tweaks ([c9d351d](https://github.com/tyrantkhan/goldsplit/commit/c9d351d8b933605391f3f2521dcf8f358fd677b5)), closes [#12](https://github.com/tyrantkhan/goldsplit/issues/12)
* timer stats display, completion fix, and UI tweaks ([053ff32](https://github.com/tyrantkhan/goldsplit/commit/053ff32153327f556fd32d94f2c95ac5a58946b1))

## [0.0.3](https://github.com/tyrantkhan/goldsplit/compare/v0.0.2...v0.0.3) (2026-02-22)


### Bug Fixes

* append contributor list to versioned release notes ([90d0af0](https://github.com/tyrantkhan/goldsplit/commit/90d0af091432931087072cc663d3baac20363be6))

## [0.0.2](https://github.com/tyrantkhan/goldsplit/compare/v0.0.1...v0.0.2) (2026-02-22)


### Features

* unify CI/CD by merging release-please into release workflow ([27c3acf](https://github.com/tyrantkhan/goldsplit/commit/27c3acf6c2b1e41617b30f6da4e64bfe19041bef))
* unify CI/CD by merging release-please into release workflow ([dd6ffcf](https://github.com/tyrantkhan/goldsplit/commit/dd6ffcf64a5451f57855868a13d4f85acd7a7f0f)), closes [#7](https://github.com/tyrantkhan/goldsplit/issues/7)

## 0.0.1 (2026-02-22)


### Features

* add About page with version info and release-please automation ([a253322](https://github.com/tyrantkhan/goldsplit/commit/a253322205be8dec308ff4173fcf5b73161ec72d))
* add About page with version info and release-please automation ([2dda0e1](https://github.com/tyrantkhan/goldsplit/commit/2dda0e19c14600d50fdae1f7a0f06ed2b42c31eb)), closes [#3](https://github.com/tyrantkhan/goldsplit/issues/3)
* add native macOS About dialog and fix menu bar name ([19d06c3](https://github.com/tyrantkhan/goldsplit/commit/19d06c31e25e9d7fcaa9dba57b89a681dc690645))


### Bug Fixes

* force initial release to 0.0.1 via release-as ([f30fbd4](https://github.com/tyrantkhan/goldsplit/commit/f30fbd44922dcdb3b0715d1b248c06f20cb49f86))
* gofmt formatting in main.go ([62b5889](https://github.com/tyrantkhan/goldsplit/commit/62b5889c78c2290415f5b465a9970dfea0ab120e))
* remove explicit token from release-please action ([b718d3f](https://github.com/tyrantkhan/goldsplit/commit/b718d3fe0079af97d19b43fd3aea09cd738e265f))
* rename Linux binary to lowercase before packaging ([c3c51e0](https://github.com/tyrantkhan/goldsplit/commit/c3c51e03d567109baa7a996d757abd878a8c3a17))
