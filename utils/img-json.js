const fs = require("fs");
const path = require("path");

fs.readdir(process.argv[2], (err, files) => {
  if (err) {
    console.error("Error reading directory:", err);
    return;
  }
  const data = files.map((filename) => {
    const extensionSeparatorIndex = filename.lastIndexOf(".");
    const type =
      filename.substring(extensionSeparatorIndex) === ".gif" ? "gif" : "img";
    return {
      filename: filename,
      type,
      tags: filename.slice(0, extensionSeparatorIndex).split("-").join(" "),
    };
  });
  console.log();
  fs.writeFile("images.json", JSON.stringify(data, undefined, 2), (err) => {
    err
      ? console.error("Error writing file:", err)
      : console.log("File written successfully!");
  });
});
