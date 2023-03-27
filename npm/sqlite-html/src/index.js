import { join } from "node:path";
import { fileURLToPath } from "node:url";
import { arch, platform } from "node:process";
import { statSync } from "node:fs";

const supportedPlatforms = [
  ["darwin", "x64"],
  ["darwin", "arm64"],
  ["win32", "x64"],
  ["linux", "x64"],
];

function validPlatform(platform, arch) {
  return (
    supportedPlatforms.find(([p, a]) => platform == p && arch === a) !== null
  );
}
function extensionSuffix(platform) {
  if (platform === "win32") return "dll";
  if (platform === "darwin") return "dylib";
  return "so";
}
function platformPackageName(platform, arch) {
  const os = platform === "win32" ? "windows" : platform;
  return `sqlite-html-${os}-${arch}`;
}

export function getLoadablePath() {
  if (!validPlatform(platform, arch)) {
    throw new Error(
      `Unsupported platform for sqlite-html, on a ${platform}-${arch} machine, but not in supported platforms (${supportedPlatforms
        .map(([p, a]) => `${p}-${a}`)
        .join(",")}). Consult the sqlite-html NPM package README for details. `
    );
  }
  const packageName = platformPackageName(platform, arch);
  const loadablePath = join(
    fileURLToPath(new URL(".", import.meta.url)),
    "..",
    "..",
    packageName,
    "lib",
    `html0.${extensionSuffix(platform)}`
  );
  if (!statSync(loadablePath, { throwIfNoEntry: false })) {
    throw new Error(
      `Loadble extension for sqlite-html not found. Was the ${packageName} package installed? Avoid using the --no-optional flag, as the optional dependencies for sqlite-html are required.`
    );
  }

  return loadablePath;
}
