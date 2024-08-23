# Changelog

## 2024-08-23
### v0.9.0
It will actually set the environment-variables properly through os.Setenv now.

It will always fall back to the last .env-file in the list. So if there are two files in env.EnvironmentFiles, then it will start by setting and using the first one, then go about and repeat this with the second file. Never unsetting, but always overwriting with new data.

## 2023-02-02
### v0.7.0
Added force-env tag.

## 2023-01-01
### v0.6.5
Fixed error where we would provide faultu value when, for example, a URL would contain a "=".

### v0.6.4
Moved the code out to root.

### v0.6.3
Updated readme.

### v0.6.2
Moved code outside of the pkg-folder and changed name of the Config-struct to Attributes.

## 2022-12-12
### v0.6.1
Updated comments and readme.

### v0.6
Added functionality for multi-level structs.

## 2022-11-11
### v0.5.4
Added a "default"-tag for environment variables and refactored the code.

### v0.5.3
Made sure pre-added values aren't overwritten.

### v0.5.2
Patched up the readme.

### v0.5.1
Updated go.mod

### v0.5
The project is reaching a point of release.

