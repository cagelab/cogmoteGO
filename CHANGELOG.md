# Changelog

cogmoteGO change log.

## [0.1.32](https://github.com/cagelab/cogmoteGO/compare/v0.1.31...v0.1.32) (2026-02-14)


### Features

* :alien: add /data/mock endpoint for test ([b00709e](https://github.com/cagelab/cogmoteGO/commit/b00709e8d486c23bd1b3ea88708f063a7b928cf2))
* :art: rename projects to cogmoteGO ([67b2922](https://github.com/cagelab/cogmoteGO/commit/67b2922003396be95f960f5b449716d9077b92da))
* :boom: add /api prefix for all endpoint ([#103](https://github.com/cagelab/cogmoteGO/issues/103)) ([b799648](https://github.com/cagelab/cogmoteGO/commit/b7996487cc159a15b8b6dab6980e23a80427b0a3))
* :egg: add more monkey sayings ([3439784](https://github.com/cagelab/cogmoteGO/commit/34397842942a0219406eb30317a3997d4bada187))
* :rocket: Add install.sh script for easy deployment and version management ([#85](https://github.com/cagelab/cogmoteGO/issues/85)) ([1cd8d15](https://github.com/cagelab/cogmoteGO/commit/1cd8d1577b4e6306b8911c17ec9f7daff5fa7319))
* :rocket: add port to server ([af86396](https://github.com/cagelab/cogmoteGO/commit/af86396acc68f2478f1f111e2510df89a55f34eb))
* :sparkles: Add API endpoints for retrieving latest broadcast data ([#106](https://github.com/cagelab/cogmoteGO/issues/106)) ([bc260c8](https://github.com/cagelab/cogmoteGO/commit/bc260c82ce761f7496ba90a847fc33899fc4b9cc))
* :sparkles: add build, test targets to justfile ([c81f081](https://github.com/cagelab/cogmoteGO/commit/c81f081b2e0f50ff9470ec5f5ccb7d59ef5ea3d1))
* :sparkles: add CLI config command for configuration management ([#158](https://github.com/cagelab/cogmoteGO/issues/158)) ([c88f36a](https://github.com/cagelab/cogmoteGO/commit/c88f36ae42fa111224b7c7f1a057be29c8fe9b48))
* :sparkles: Add correct rate calculation to trial data ([#114](https://github.com/cagelab/cogmoteGO/issues/114)) ([67077d4](https://github.com/cagelab/cogmoteGO/commit/67077d40af86a4e38c1b8eb6b04105b3af7583b5))
* :sparkles: Add CORS middleware to allow requests from localhost:1420 ([#98](https://github.com/cagelab/cogmoteGO/issues/98)) ([fd8014b](https://github.com/cagelab/cogmoteGO/commit/fd8014b0d64b581fbefaf377372fb3c00f45a7e2))
* :sparkles: add default data endpoint: data ([68077b6](https://github.com/cagelab/cogmoteGO/commit/68077b6d18823acbf4e30d067f233e863983d1c3))
* :sparkles: add DELETE endpoints for command proxies and improve error handling ([d785dc7](https://github.com/cagelab/cogmoteGO/commit/d785dc745dd5a9b8a5b9d4e24a8cc4edb65b6f73))
* :sparkles: Add device routes to display system and hardware information ([#102](https://github.com/cagelab/cogmoteGO/issues/102)) ([3bddb38](https://github.com/cagelab/cogmoteGO/commit/3bddb381ec9967a46a00982ab63e7d7a36da9465))
* :sparkles: add email config and recipients management APIs ([#152](https://github.com/cagelab/cogmoteGO/issues/152)) ([fb8a38b](https://github.com/cagelab/cogmoteGO/commit/fb8a38b501e713880abaec93bb0724084b307c8f))
* :sparkles: Add email credentials storage with keyring ([#140](https://github.com/cagelab/cogmoteGO/issues/140)) ([30c3659](https://github.com/cagelab/cogmoteGO/commit/30c365916071874471ebbfc62b083a06f33d2bd9))
* :sparkles: Add experiment status management endpoints ([#125](https://github.com/cagelab/cogmoteGO/issues/125)) ([7e8dcd2](https://github.com/cagelab/cogmoteGO/commit/7e8dcd28badd70a322fe08ba53a75d0b555b4b9e))
* :sparkles: add experiments module with CRUD and start/stop endpoints ([5843102](https://github.com/cagelab/cogmoteGO/commit/584310280c60443a74cdbc0ba8a80b3a9e34d7bc))
* :sparkles: add GET /alive endpoint ([d912fe6](https://github.com/cagelab/cogmoteGO/commit/d912fe6ab0521fbd0f9592f50b91fa34023120c6))
* :sparkles: add GIN_MODE env var support ([971aa1e](https://github.com/cagelab/cogmoteGO/commit/971aa1ee9e95dd8205c4f3b76c9d9aada5cd9f7f))
* :sparkles: Add GIN_MODE env var support and H2C support ([7308b10](https://github.com/cagelab/cogmoteGO/commit/7308b10e38a37e822543306aff863200d3bf3a03))
* :sparkles: add git-cliff for generate changelog ([9af8c05](https://github.com/cagelab/cogmoteGO/commit/9af8c05eac2f4250b04464c1098f3a284164ef45))
* :sparkles: add instance unique identifier ([#153](https://github.com/cagelab/cogmoteGO/issues/153)) ([03ccab6](https://github.com/cagelab/cogmoteGO/commit/03ccab699793cd79f5fe18297bd97f64c3c3aed6))
* :sparkles: Add macOS AMD64 build target to GitHub workflow ([#93](https://github.com/cagelab/cogmoteGO/issues/93)) ([05e8c2e](https://github.com/cagelab/cogmoteGO/commit/05e8c2e68dfc90c2eea20f169e6f76f8381cf768))
* :sparkles: add request body to create stream endpoint ([3037186](https://github.com/cagelab/cogmoteGO/commit/3037186cf20f0f142dd23029a088c1f280984fd6))
* :sparkles: add slog as new log system ([#47](https://github.com/cagelab/cogmoteGO/issues/47)) ([7e68cf1](https://github.com/cagelab/cogmoteGO/commit/7e68cf19ddce8f4fb6672d7846ae339e54acbd48))
* :sparkles: add timeout for handshake (5 secs) ([#45](https://github.com/cagelab/cogmoteGO/issues/45)) ([bc3e962](https://github.com/cagelab/cogmoteGO/commit/bc3e96268e1b36bdd17be8cb9e02209720e6d027))
* :sparkles: Add tools directory with Python mock data generator ([#112](https://github.com/cagelab/cogmoteGO/issues/112)) ([e2e2ba4](https://github.com/cagelab/cogmoteGO/commit/e2e2ba4bed86c4549729f3664347f650c2cce89c))
* :sparkles: Add verbose flag ([#88](https://github.com/cagelab/cogmoteGO/issues/88)) ([94d86e2](https://github.com/cagelab/cogmoteGO/commit/94d86e22718aad56c4c75b165c2953b3f9e6157e))
* :sparkles: Add version info to device report ([#123](https://github.com/cagelab/cogmoteGO/issues/123)) ([d2b5838](https://github.com/cagelab/cogmoteGO/commit/d2b583802cd48dd69fcef58df0bcf2a7abb7cf4c))
* :sparkles: Add Windows installer and update README ([#117](https://github.com/cagelab/cogmoteGO/issues/117)) ([eab0c90](https://github.com/cagelab/cogmoteGO/commit/eab0c90a490ea75fa7e7532016dcda3d950bad51))
* :sparkles: Added /health endpoint for quickly checking cogmoteGO running status ([#11](https://github.com/cagelab/cogmoteGO/issues/11)) ([01970d6](https://github.com/cagelab/cogmoteGO/commit/01970d6466364a958768bbcfcc16256da9bafc00))
* :sparkles: added native services for Mac, Linux, and Windows ([#83](https://github.com/cagelab/cogmoteGO/issues/83)) ([7e9e074](https://github.com/cagelab/cogmoteGO/commit/7e9e07460d6de378ec74dc601598bedc7ccabf73))
* :sparkles: complete the email system prototype ([#141](https://github.com/cagelab/cogmoteGO/issues/141)) ([08b6050](https://github.com/cagelab/cogmoteGO/commit/08b60504668cf99698539dd7cbda817eb8929542))
* :sparkles: Enhance email with HTML and attachments support ([#143](https://github.com/cagelab/cogmoteGO/issues/143)) ([d48f975](https://github.com/cagelab/cogmoteGO/commit/d48f9759b89d01ad3495acee3cb292050ed78ed4))
* :sparkles: Enhance trial data with statistics tracking ([#128](https://github.com/cagelab/cogmoteGO/issues/128)) ([fc77134](https://github.com/cagelab/cogmoteGO/commit/fc771341412c8f459d2a777e4856b54aa542290b))
* :sparkles: Improve experimental management capabilities ([#84](https://github.com/cagelab/cogmoteGO/issues/84)) ([f6eeabb](https://github.com/cagelab/cogmoteGO/commit/f6eeabb72f5c4167f606904d6d944348576b89c1))
* :sparkles: make server port configurable ([#157](https://github.com/cagelab/cogmoteGO/issues/157)) ([ab230bc](https://github.com/cagelab/cogmoteGO/commit/ab230bc8c22b69176d227bdeecc1b5a925304151))
* :sparkles: preliminarily implement the lazy pirate pattern and configurable timeout ([#150](https://github.com/cagelab/cogmoteGO/issues/150)) ([6f8111c](https://github.com/cagelab/cogmoteGO/commit/6f8111c6dcfb6faf78a6cb12ac6be5f6e3ff371b))
* :sparkles: upload data to obs-studio ([#135](https://github.com/cagelab/cogmoteGO/issues/135)) ([64f487b](https://github.com/cagelab/cogmoteGO/commit/64f487b19bb4a7ade5e621737745384b41755735))
* :tada: init project ([6d212f3](https://github.com/cagelab/cogmoteGO/commit/6d212f34c37a95944b6490f0b36fe6cd10b1d16a))
* :watch: add time benchmarking for command proxy ([#62](https://github.com/cagelab/cogmoteGO/issues/62)) ([2918f72](https://github.com/cagelab/cogmoteGO/commit/2918f7231966e990ad8a14363b16b26e779633ad))
* :wrench: Replace kardianos/service with Ccccraz fork and update dependencies ([#108](https://github.com/cagelab/cogmoteGO/issues/108)) ([9082f71](https://github.com/cagelab/cogmoteGO/commit/9082f717b2e9614df8c5e23cff4584f8053898b8))
* :wrench: Update Windows installer exit behavior ([#118](https://github.com/cagelab/cogmoteGO/issues/118)) ([85143b5](https://github.com/cagelab/cogmoteGO/commit/85143b5971381abf8a4dbee21efb10d862600e2b))
* ✨ add availability check and improve timeout  ([#63](https://github.com/cagelab/cogmoteGO/issues/63)) ([1eb5208](https://github.com/cagelab/cogmoteGO/commit/1eb52086b7afe8cd464890284730eaeebf817c3f))
* ✨ add cobra to support command-line arguments such as: --version ([#28](https://github.com/cagelab/cogmoteGO/issues/28)) ([563c0bf](https://github.com/cagelab/cogmoteGO/commit/563c0bf9dc2c11aa2d79fe11f178bdd2e607c5dc))
* 🎉 add command proxy endpoint ([#3](https://github.com/cagelab/cogmoteGO/issues/3)) ([2230af5](https://github.com/cagelab/cogmoteGO/commit/2230af55812f52cfcfbdf4c3e9f0df00fc8c1e83))
* 🐳 add Docker support with multi-stage build ([#60](https://github.com/cagelab/cogmoteGO/issues/60)) ([708e87e](https://github.com/cagelab/cogmoteGO/commit/708e87ebf373a8b4f31ea2891b67979772056a67))
* 🛠️ add systemd service file for running cogmoteGO as a daemon ([#75](https://github.com/cagelab/cogmoteGO/issues/75)) ([c1c16b4](https://github.com/cagelab/cogmoteGO/commit/c1c16b442d4d368135b6cb5647749074b501a934))
* add build, clean, deps targets to Makefile ([a52e944](https://github.com/cagelab/cogmoteGO/commit/a52e944c6efb343bf21e9b9d79c0422f40c71259))


### Bug Fixes

* :ambulance: Add compatibility with sh ([#86](https://github.com/cagelab/cogmoteGO/issues/86)) ([1f26605](https://github.com/cagelab/cogmoteGO/commit/1f2660590d509ab4a1a14b70d11d315fbaaeb166))
* :ambulance: Auto-detect service username and require password parameter ([#89](https://github.com/cagelab/cogmoteGO/issues/89)) ([b090ed9](https://github.com/cagelab/cogmoteGO/commit/b090ed9da2ab5a9ad15ae6d8ff3483f18a8ca376))
* :ambulance: fix the broadcast module ([1908a24](https://github.com/cagelab/cogmoteGO/commit/1908a24a3c0c1ae819c86ad90f3153163962cf4b))
* :ambulance: fixed crash issue on proxy delete ([#44](https://github.com/cagelab/cogmoteGO/issues/44)) ([c1c50ed](https://github.com/cagelab/cogmoteGO/commit/c1c50ed64d4bd8e91b327d6fd16f7c5c70bfbc0f))
* :ambulance: hotfix timeout problem ([#79](https://github.com/cagelab/cogmoteGO/issues/79)) ([d51d083](https://github.com/cagelab/cogmoteGO/commit/d51d0839984be47d759b68383ba6779942cb3d52))
* :bug: add newline character to version output ([#50](https://github.com/cagelab/cogmoteGO/issues/50)) ([210930c](https://github.com/cagelab/cogmoteGO/commit/210930c5448004502619c2daf2c37a185aea145b))
* :bug: add the missing Get /cmds/proxies endpoint ([726fe22](https://github.com/cagelab/cogmoteGO/commit/726fe2275c22a8fb46d6507f0a928f3c404844b7))
* :bug: Change serviceCmd password flags from Persistent to local ([#91](https://github.com/cagelab/cogmoteGO/issues/91)) ([a334a6e](https://github.com/cagelab/cogmoteGO/commit/a334a6e0221b6521d9faa869ba1a7f43e332a9fb))
* :bug: fix changelog path problem ([f1e5e0e](https://github.com/cagelab/cogmoteGO/commit/f1e5e0e52a581a15e9ae8391926573eba2cb6bfe))
* :bug: fix handshake logic bug ([e702639](https://github.com/cagelab/cogmoteGO/commit/e702639357a89e5d4104b2d4973d6b6962715371))
* :bug: fix permission problem for changelog ([01474bd](https://github.com/cagelab/cogmoteGO/commit/01474bddfd45cd5d6c5f7d7817249605220592c5))
* :bug: fix permission problem for changelog ([b1218d9](https://github.com/cagelab/cogmoteGO/commit/b1218d9de465135bb96fa699d3f479e70973a23c))
* :bug: fix permission problem for changlog realse ([4fdc57a](https://github.com/cagelab/cogmoteGO/commit/4fdc57ad82ff275b839fb6165fa4af8d43b9f25e))
* :bug: Fixed an issue where matlab would block when the 47th message was sent ([0b2ec0d](https://github.com/cagelab/cogmoteGO/commit/0b2ec0d41bbe6ea82fc5d86b2d1f95d4b4500ca2))
* :bug: Fixed the error that accessing Delete /cmds/proxies for the second time would time out ([#77](https://github.com/cagelab/cogmoteGO/issues/77)) ([0bbcb6e](https://github.com/cagelab/cogmoteGO/commit/0bbcb6ee658995ec26a339c3723754c94fa69d01))
* :bug: Fixed the issue where the program does not exit automatically when the port is not available ([#145](https://github.com/cagelab/cogmoteGO/issues/145)) ([1f6ecf1](https://github.com/cagelab/cogmoteGO/commit/1f6ecf1992adeaea6d52ade973c3ec8b089131a8))
* :bug: remove validation_level in release.yml ([5c2e0a6](https://github.com/cagelab/cogmoteGO/commit/5c2e0a695393242bd0d972f0195a8e18b314d58d))
* :bug: select correct version for changelog ([75fd759](https://github.com/cagelab/cogmoteGO/commit/75fd759f01d94fe8ad4d2160a3493bca9b21c654))
* :bug: update CORS config to allow electron origin ([#124](https://github.com/cagelab/cogmoteGO/issues/124)) ([d1d1fc3](https://github.com/cagelab/cogmoteGO/commit/d1d1fc392c74bf6bf8d7527bffd0d0ab386e0158))
* :bug: update CORS config to allow tauri://localhost origin ([#120](https://github.com/cagelab/cogmoteGO/issues/120)) ([1ea1ba7](https://github.com/cagelab/cogmoteGO/commit/1ea1ba7f9797cfdd929e135099a14465029338c0))
* :bug: update CORS config to allow tauri.localhost origin for windows ([#121](https://github.com/cagelab/cogmoteGO/issues/121)) ([78f78c2](https://github.com/cagelab/cogmoteGO/commit/78f78c247371fa46dedf9c3c2888d9bda0f28311))
* :pencil2: bordercast_endpoints -&gt; broadcast_endpoints ([f03e2cb](https://github.com/cagelab/cogmoteGO/commit/f03e2cb94a52f489c35c1c9bd350e756096be13c))
* ⏱️ replace context timeout with socket timeout in handshake process ([#66](https://github.com/cagelab/cogmoteGO/issues/66)) ([5dfc34d](https://github.com/cagelab/cogmoteGO/commit/5dfc34de4ed81801bdaeb20bf6f5d419e1446b77))
* ⏱️ replace context timeout with socket timeout in handshake process ([#67](https://github.com/cagelab/cogmoteGO/issues/67)) ([82477a2](https://github.com/cagelab/cogmoteGO/commit/82477a20b4667b9f7445434b3017b4c02da96c6b))
* **ci:** :green_heart: add release-please to upload-release job dependencies ([#166](https://github.com/cagelab/cogmoteGO/issues/166)) ([9c22dbc](https://github.com/cagelab/cogmoteGO/commit/9c22dbc44dadc5902a2362b87623b181609abf41))

## [0.1.31](https://github.com/cagelab/cogmoteGO/compare/v0.1.30...v0.1.31) (2026-02-14)


### Features

* :sparkles: add CLI config command for configuration management ([#158](https://github.com/cagelab/cogmoteGO/issues/158)) ([c88f36a](https://github.com/cagelab/cogmoteGO/commit/c88f36ae42fa111224b7c7f1a057be29c8fe9b48))
* :sparkles: add email config and recipients management APIs ([#152](https://github.com/cagelab/cogmoteGO/issues/152)) ([fb8a38b](https://github.com/cagelab/cogmoteGO/commit/fb8a38b501e713880abaec93bb0724084b307c8f))
* :sparkles: add instance unique identifier ([#153](https://github.com/cagelab/cogmoteGO/issues/153)) ([03ccab6](https://github.com/cagelab/cogmoteGO/commit/03ccab699793cd79f5fe18297bd97f64c3c3aed6))
* :sparkles: make server port configurable ([#157](https://github.com/cagelab/cogmoteGO/issues/157)) ([ab230bc](https://github.com/cagelab/cogmoteGO/commit/ab230bc8c22b69176d227bdeecc1b5a925304151))

## [unreleased]

### 🚜 Refactor

- :art: fix output format with old powershell (#134) by @Ccccraz in [#134](https://github.com/Ccccraz/cogmoteGO/pull/134)



## [0.1.24] - 2025-09-11

### 🚀 Features

- :sparkles: Enhance trial data with statistics tracking (#128) by @Ccccraz in [#128](https://github.com/Ccccraz/cogmoteGO/pull/128)


### 🚜 Refactor

- :recycle: Move status endpoints to separate module (#127) by @Ccccraz in [#127](https://github.com/Ccccraz/cogmoteGO/pull/127)



**Full Changelog**: https://github.com/Ccccraz/cogmoteGO/compare/v0.1.23...v0.1.24

## [0.1.23] - 2025-09-06

### 🚀 Features

- :sparkles: Add experiment status management endpoints (#125) by @Ccccraz in [#125](https://github.com/Ccccraz/cogmoteGO/pull/125)


### 🐛 Bug Fixes

- :bug: update CORS config to allow electron origin (#124) by @Ccccraz in [#124](https://github.com/Ccccraz/cogmoteGO/pull/124)



**Full Changelog**: https://github.com/Ccccraz/cogmoteGO/compare/v0.1.22...v0.1.23

## [0.1.22] - 2025-08-17

### 🚀 Features

- :sparkles: Add version info to device report (#123) by @Ccccraz in [#123](https://github.com/Ccccraz/cogmoteGO/pull/123)



**Full Changelog**: https://github.com/Ccccraz/cogmoteGO/compare/v0.1.21...v0.1.22

## [0.1.21] - 2025-08-14

### 🐛 Bug Fixes

- :bug: update CORS config to allow tauri.localhost origin for windows (#121) by @Ccccraz in [#121](https://github.com/Ccccraz/cogmoteGO/pull/121)


### 📚 Documentation

- :pencil2: lost '| sh' part by @Ccccraz



**Full Changelog**: https://github.com/Ccccraz/cogmoteGO/compare/v0.1.20...v0.1.21

## [0.1.20] - 2025-08-14

### 🚀 Features

- :sparkles: Add Windows installer and update README (#117) by @Ccccraz in [#117](https://github.com/Ccccraz/cogmoteGO/pull/117)

- :wrench: Update Windows installer exit behavior (#118) by @Ccccraz in [#118](https://github.com/Ccccraz/cogmoteGO/pull/118)


### 🐛 Bug Fixes

- :bug: update CORS config to allow tauri://localhost origin (#120) by @Ccccraz in [#120](https://github.com/Ccccraz/cogmoteGO/pull/120)



**Full Changelog**: https://github.com/Ccccraz/cogmoteGO/compare/v0.1.19...v0.1.20

## [0.1.19] - 2025-08-10

### 🚀 Features

- :sparkles: Add tools directory with Python mock data generator (#112) by @Ccccraz in [#112](https://github.com/Ccccraz/cogmoteGO/pull/112)

- :sparkles: Add correct rate calculation to trial data (#114) by @Ccccraz in [#114](https://github.com/Ccccraz/cogmoteGO/pull/114)



**Full Changelog**: https://github.com/Ccccraz/cogmoteGO/compare/v0.1.18...v0.1.19

## [0.1.17] - 2025-07-23

### 🚀 Features

- :wrench: Replace kardianos/service with Ccccraz fork and update dependencies (#108) by @Ccccraz in [#108](https://github.com/Ccccraz/cogmoteGO/pull/108)



## New Contributors
* @iandol made their first contribution in [#107](https://github.com/Ccccraz/cogmoteGO/pull/107)

**Full Changelog**: https://github.com/Ccccraz/cogmoteGO/compare/v0.1.16...v0.1.17

## [0.1.16] - 2025-07-08

### 🚀 Features

- :sparkles: Add API endpoints for retrieving latest broadcast data (#106) by @Ccccraz in [#106](https://github.com/Ccccraz/cogmoteGO/pull/106)



**Full Changelog**: https://github.com/Ccccraz/cogmoteGO/compare/v0.1.15...v0.1.16

## [0.1.15] - 2025-06-21

### 🚀 Features

- :sparkles: Add CORS middleware to allow requests from localhost:1420 (#98) by @Ccccraz in [#98](https://github.com/Ccccraz/cogmoteGO/pull/98)

- :sparkles: Add device routes to display system and hardware information (#102) by @Ccccraz in [#102](https://github.com/Ccccraz/cogmoteGO/pull/102)

- :boom: add /api prefix for all endpoint (#103) by @Ccccraz in [#103](https://github.com/Ccccraz/cogmoteGO/pull/103)


### 📚 Documentation

- :memo: add winget (#95) by @Ccccraz in [#95](https://github.com/Ccccraz/cogmoteGO/pull/95)



**Full Changelog**: https://github.com/Ccccraz/cogmoteGO/compare/v0.1.14...v0.1.15

## [0.1.14] - 2025-06-15

### 🚀 Features

- :rocket: Add install.sh script for easy deployment and version management (#85) by @Ccccraz in [#85](https://github.com/Ccccraz/cogmoteGO/pull/85)

- :sparkles: Add verbose flag (#88) by @Ccccraz in [#88](https://github.com/Ccccraz/cogmoteGO/pull/88)

- :sparkles: Add macOS AMD64 build target to GitHub workflow (#93) by @Ccccraz in [#93](https://github.com/Ccccraz/cogmoteGO/pull/93)


### 🐛 Bug Fixes

- :ambulance: Add compatibility with sh (#86) by @Ccccraz in [#86](https://github.com/Ccccraz/cogmoteGO/pull/86)

- :ambulance: Auto-detect service username and require password parameter (#89) by @Ccccraz in [#89](https://github.com/Ccccraz/cogmoteGO/pull/89)

- :bug: Change serviceCmd password flags from Persistent to local (#91) by @Ccccraz in [#91](https://github.com/Ccccraz/cogmoteGO/pull/91)


### 📚 Documentation

- :memo: Add installation and service instructions to README (#92) by @Ccccraz in [#92](https://github.com/Ccccraz/cogmoteGO/pull/92)


### ⚙️ Miscellaneous Tasks

- :arrow_up: Upgrade Go version from 1.24.3 to 1.24.4 (#94) (#94) by @Ccccraz in [#94](https://github.com/Ccccraz/cogmoteGO/pull/94)



**Full Changelog**: https://github.com/Ccccraz/cogmoteGO/compare/v0.1.13...v0.1.14

## [0.1.13] - 2025-06-14

### 🚀 Features

- :sparkles: added native services for Mac, Linux, and Windows (#83) by @Ccccraz in [#83](https://github.com/Ccccraz/cogmoteGO/pull/83)

- :sparkles: Improve experimental management capabilities (#84) by @Ccccraz in [#84](https://github.com/Ccccraz/cogmoteGO/pull/84)



**Full Changelog**: https://github.com/Ccccraz/cogmoteGO/compare/v0.1.12...v0.1.13

## [0.1.12] - 2025-06-06

### 🐛 Bug Fixes

- :ambulance: hotfix timeout problem (#79) by @Ccccraz in [#79](https://github.com/Ccccraz/cogmoteGO/pull/79)



**Full Changelog**: https://github.com/Ccccraz/cogmoteGO/compare/v0.1.11...v0.1.12

## [0.1.11] - 2025-06-05

### 🚀 Features

- 🛠️ add systemd service file for running cogmoteGO as a daemon (#75) by @Ccccraz in [#75](https://github.com/Ccccraz/cogmoteGO/pull/75)


### 🐛 Bug Fixes

- :bug: Fixed the error that accessing Delete /cmds/proxies for the second time would time out (#77) by @Ccccraz in [#77](https://github.com/Ccccraz/cogmoteGO/pull/77)



**Full Changelog**: https://github.com/Ccccraz/cogmoteGO/compare/v0.1.10...v0.1.11

## [0.1.10] - 2025-05-25

### 🚀 Features

- 🐳 add Docker support with multi-stage build (#60) by @Ccccraz in [#60](https://github.com/Ccccraz/cogmoteGO/pull/60)

- :watch: add time benchmarking for command proxy (#62) by @Ccccraz in [#62](https://github.com/Ccccraz/cogmoteGO/pull/62)

- ✨ add availability check and improve timeout  (#63) by @Ccccraz in [#63](https://github.com/Ccccraz/cogmoteGO/pull/63)


### 🐛 Bug Fixes

- ⏱️ replace context timeout with socket timeout in handshake process (#66) by @Ccccraz in [#66](https://github.com/Ccccraz/cogmoteGO/pull/66)

- ⏱️ replace context timeout with socket timeout in handshake process (#67) by @Ccccraz in [#67](https://github.com/Ccccraz/cogmoteGO/pull/67)


### 🚜 Refactor

- 🔄 simplify release workflow by removing tag message extraction (#68) by @Ccccraz in [#68](https://github.com/Ccccraz/cogmoteGO/pull/68)



**Full Changelog**: https://github.com/Ccccraz/cogmoteGO/compare/v0.1.9...v0.1.10

## [0.1.9] - 2025-05-24

### 🐛 Bug Fixes

- :bug: add newline character to version output (#50) by @Ccccraz in [#50](https://github.com/Ccccraz/cogmoteGO/pull/50)



**Full Changelog**: https://github.com/Ccccraz/cogmoteGO/compare/v0.1.8...v0.1.9

## [0.1.8] - 2025-05-24

### 🚜 Refactor

- :recycle: standardize API error responses using commonTypes.APIError (#49) by @Ccccraz in [#49](https://github.com/Ccccraz/cogmoteGO/pull/49)



**Full Changelog**: https://github.com/Ccccraz/cogmoteGO/compare/v0.1.7...v0.1.8

## [0.1.7] - 2025-05-23

### 🚀 Features

- :sparkles: add slog as new log system (#47) by @Ccccraz in [#47](https://github.com/Ccccraz/cogmoteGO/pull/47)



**Full Changelog**: https://github.com/Ccccraz/cogmoteGO/compare/v0.1.6...v0.1.7

## [0.1.6] - 2025-05-23

### 🚀 Features

- :sparkles: add timeout for handshake (5 secs) (#45) by @Ccccraz in [#45](https://github.com/Ccccraz/cogmoteGO/pull/45)


### 🐛 Bug Fixes

- :ambulance: fixed crash issue on proxy delete (#44) by @Ccccraz in [#44](https://github.com/Ccccraz/cogmoteGO/pull/44)


### 📚 Documentation

- :memo: update readme (#36) by @Ccccraz in [#36](https://github.com/Ccccraz/cogmoteGO/pull/36)



**Full Changelog**: https://github.com/Ccccraz/cogmoteGO/compare/v0.1.5...v0.1.6

## [0.1.5] - 2025-05-17

### 🚀 Features

- :egg: add more monkey sayings by @Ccccraz in [#20](https://github.com/Ccccraz/cogmoteGO/pull/20)

- ✨ add cobra to support command-line arguments such as: --version (#28) by @Ccccraz in [#28](https://github.com/Ccccraz/cogmoteGO/pull/28)


### 🐛 Bug Fixes

- :pencil2: bordercast_endpoints -> broadcast_endpoints by @Ccccraz in [#26](https://github.com/Ccccraz/cogmoteGO/pull/26)

- :bug: add the missing Get /cmds/proxies endpoint by @Ccccraz in [#27](https://github.com/Ccccraz/cogmoteGO/pull/27)


### 💼 Other

- :wastebasket: give up justfile by @Ccccraz in [#29](https://github.com/Ccccraz/cogmoteGO/pull/29)

- :building_construction: add Windows build script with CGO and version embedding by @Ccccraz

- :building_construction: add Linux/macOS build script with CGO and version embedding by @Ccccraz in [#31](https://github.com/Ccccraz/cogmoteGO/pull/31)

- :ambulance: fix the problem of pre-release version number formatting error by @Ccccraz in [#32](https://github.com/Ccccraz/cogmoteGO/pull/32)


### ⚙️ Miscellaneous Tasks

- 📝 add GNU LESSER GENERAL PUBLIC LICENSE by @Ccccraz in [#30](https://github.com/Ccccraz/cogmoteGO/pull/30)

- 👷 optimized the release and build (#33) by @Ccccraz in [#33](https://github.com/Ccccraz/cogmoteGO/pull/33)



**Full Changelog**: https://github.com/Ccccraz/cogmoteGO/compare/v0.1.4...v0.1.5

## [0.1.4] - 2025-05-11

### 🐛 Bug Fixes

- :ambulance: fix the broadcast module by @Ccccraz in [#19](https://github.com/Ccccraz/cogmoteGO/pull/19)



**Full Changelog**: https://github.com/Ccccraz/cogmoteGO/compare/v0.1.3...v0.1.4

## [0.1.3] - 2025-05-11

### 🚀 Features

- :sparkles: add experiments module with CRUD and start/stop endpoints by @Ccccraz in [#18](https://github.com/Ccccraz/cogmoteGO/pull/18)



**Full Changelog**: https://github.com/Ccccraz/cogmoteGO/compare/v0.1.3-alpha.8...v0.1.3

## [0.1.3-alpha.8] - 2025-05-10

### 🚀 Features

- :sparkles: add DELETE endpoints for command proxies and improve error handling by @Ccccraz in [#17](https://github.com/Ccccraz/cogmoteGO/pull/17)



**Full Changelog**: https://github.com/Ccccraz/cogmoteGO/compare/v0.1.3-alpha.7...v0.1.3-alpha.8

## [0.1.3-alpha.7] - 2025-05-10

### 🚀 Features

- :sparkles: add GET /alive endpoint by @Ccccraz in [#16](https://github.com/Ccccraz/cogmoteGO/pull/16)



**Full Changelog**: https://github.com/Ccccraz/cogmoteGO/compare/v0.1.3-alpha.6...v0.1.3-alpha.7

## [0.1.3-alpha.6] - 2025-05-09

### 📚 Documentation

- :memo: add new changelog rules by @Ccccraz

- 📝 named release (#15) by @Ccccraz in [#15](https://github.com/Ccccraz/cogmoteGO/pull/15)


### ⚙️ Miscellaneous Tasks

- :bug: fix release workflow tag handling by @Ccccraz



**Full Changelog**: https://github.com/Ccccraz/cogmoteGO/compare/v0.1.3-alpha.5...v0.1.3-alpha.6

## [0.1.3-alpha.5] - 2025-05-09

### 🚀 Features

- :sparkles: add git-cliff for generate changelog by @Ccccraz in [#14](https://github.com/Ccccraz/cogmoteGO/pull/14)


### 🐛 Bug Fixes

- :bug: fix permission problem for changlog realse by @Ccccraz

- :bug: fix permission problem for changelog by @Ccccraz

- :bug: fix changelog path problem by @Ccccraz

- :bug: fix permission problem for changelog by @Ccccraz

- :bug: select correct version for changelog by @Ccccraz

- :bug: remove validation_level in release.yml by @Ccccraz


### 📚 Documentation

- :loud_sound: by @Ccccraz

- :loud_sound: by @Ccccraz

- :loud_sound: update changelog by @Ccccraz

- :loud_sound: test log by @Ccccraz


### ⚙️ Miscellaneous Tasks

- Test by @Ccccraz

- :bug: get correct tag name by @Ccccraz

- Update changelog and fix release issues by @Ccccraz



**Full Changelog**: https://github.com/Ccccraz/cogmoteGO/compare/v0.1.3-alpha.4...v0.1.3-alpha.5

## [0.1.3-alpha.4] - 2025-05-04

### 🐛 Bug Fixes

- :bug: fix handshake logic bug by @Ccccraz


### 🚜 Refactor

- :heavy_minus_sign: remove uuid deps by @Ccccraz



**Full Changelog**: https://github.com/Ccccraz/cogmoteGO/compare/v0.1.3-alpha.3...v0.1.3-alpha.4

## [0.1.3-alpha.3] - 2025-04-22

### 🚀 Features

- :alien: add /data/mock endpoint for test by @Ccccraz in [#13](https://github.com/Ccccraz/cogmoteGO/pull/13)


### ⚙️ Miscellaneous Tasks

- :recycle: direct require gopsutil (#12) by @Ccccraz in [#12](https://github.com/Ccccraz/cogmoteGO/pull/12)

- :heavy_minus_sign: remove redundancy mod cmp by @Ccccraz



**Full Changelog**: https://github.com/Ccccraz/cogmoteGO/compare/v0.1.3-alpha.2...v0.1.3-alpha.3

## [0.1.3-alpha.2] - 2025-04-21

### 🚀 Features

- :sparkles: Added /health endpoint for quickly checking cogmoteGO running status (#11) by @Ccccraz in [#11](https://github.com/Ccccraz/cogmoteGO/pull/11)



**Full Changelog**: https://github.com/Ccccraz/cogmoteGO/compare/v0.1.3-alpha.1...v0.1.3-alpha.2

## [0.1.3-alpha.1] - 2025-04-21

### 🚀 Features

- 🎉 add command proxy endpoint (#3) by @Ccccraz in [#3](https://github.com/Ccccraz/cogmoteGO/pull/3)


### 💼 Other

- :green_heart: Fix building and releasing CI/CD to fit CGO enabled projects (#4) by @Ccccraz in [#4](https://github.com/Ccccraz/cogmoteGO/pull/4)

- :green_heart: add libzmq-mt-4_3_5.dll in windows release build (#5) by @Ccccraz in [#5](https://github.com/Ccccraz/cogmoteGO/pull/5)

- 💚 Get the correct tagname (#7) by @Ccccraz in [#7](https://github.com/Ccccraz/cogmoteGO/pull/7)

- :construction_worker: add cache for vcpkg (#8) by @Ccccraz in [#8](https://github.com/Ccccraz/cogmoteGO/pull/8)

- :green_heart: correct get artifact by @Ccccraz

- :green_heart: Abandon multi-workflow solution (#9) by @Ccccraz in [#9](https://github.com/Ccccraz/cogmoteGO/pull/9)

- 👷 Add build workflow for macOS (#10) by @Ccccraz in [#10](https://github.com/Ccccraz/cogmoteGO/pull/10)


### 🚜 Refactor

- Rename cogmoteGO.go to main.go by @Ccccraz

- :art: Modify the project structure to prepare for future by @Ccccraz in [#2](https://github.com/Ccccraz/cogmoteGO/pull/2)

- :art: Separate build and release into two workflows (#6) by @Ccccraz in [#6](https://github.com/Ccccraz/cogmoteGO/pull/6)

- 🎨 Separate build and release into two workflows by @Ccccraz



**Full Changelog**: https://github.com/Ccccraz/cogmoteGO/compare/v0.1.2...v0.1.3-alpha.1

## [0.1.2] - 2025-04-15

### 🚀 Features

- :tada: init project by @Ccccraz

- :rocket: add port to server by @Ccccraz

- Add build, clean, deps targets to Makefile by @Ccccraz

- :sparkles: add default data endpoint: data by @Ccccraz

- :art: rename projects to cogmoteGO by @Ccccraz

- :sparkles: add build, test targets to justfile by @Ccccraz

- :sparkles: add GIN_MODE env var support by @Ccccraz

- :sparkles: add request body to create stream endpoint by @Ccccraz

- :sparkles: Add GIN_MODE env var support and H2C support by @Ccccraz


### 🐛 Bug Fixes

- :bug: Fixed an issue where matlab would block when the 47th message was sent by @Ccccraz


### 💼 Other

- :fire: give makefile sloution by @Ccccraz

- :construction_worker: init auto build and release action (#1) by @Ccccraz in [#1](https://github.com/Ccccraz/cogmoteGO/pull/1)


### 🚜 Refactor

- :truck: move main.go to cmd/main.go by @Ccccraz

- :art: Simplify the project structure by @Ccccraz

- :zap: Concurrent processing of all subscribers by @Ccccraz

- :zap: add buffered channels to subscribers by @Ccccraz

- :art: Simplify the project structure by @Ccccraz


### 📚 Documentation

- :memo: Init readme by @Ccccraz


### 🧪 Testing

- :test_tube: add e2e tests for puremote-server by @Ccccraz

- :test_tube: update new test for default data endpoint by @Ccccraz

- :recycle: rename go file to cogmoteGO by @Ccccraz


### ⚙️ Miscellaneous Tasks

- :construction: by @Ccccraz

- Add bin/ to .gitignore by @Ccccraz

- :arrow_up: gin by @Ccccraz

- :art: remove unused signal import by @Ccccraz



## New Contributors
* @Ccccraz made their first contribution

**Full Changelog**: https://github.com/Ccccraz/cogmoteGO/compare/v0.1.1...v0.1.2

<!-- generated by git-cliff -->
