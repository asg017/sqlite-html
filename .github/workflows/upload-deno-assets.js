const fs = require("fs").promises;

const compiled_extensions = [
  {
    path: "sqlite-html-macos/html0.dylib",
    name: "deno-darwin-x86_64.html0.dylib",
  },
  {
    path: "sqlite-html-linux_x86/html0.so",
    name: "deno-linux-x86_64.html0.so",
  },
  {
    path: "sqlite-html-windows/html0.dll",
    name: "deno-windows-x86_64.html0.dll",
  },
];

module.exports = async ({ github, context }) => {
  const { owner, repo } = context.repo;
  const release = await github.rest.repos.getReleaseByTag({
    owner,
    repo,
    tag: process.env.GITHUB_REF.replace("refs/tags/", ""),
  });
  const release_id = release.data.id;

  await Promise.all(
    compiled_extensions.map(async ({ name, path }) => {
      return github.rest.repos.uploadReleaseAsset({
        owner,
        repo,
        release_id,
        name,
        data: await fs.readFile(path),
      });
    })
  );
};
