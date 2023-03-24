const showdown = require("showdown");
const fs = require("fs/promises");
const path = require("path");
const walkdir = require("walkdir");
const hljs = require("highlight.js");
const { argv } = require("yargs");
const fm = require("frontmatter");

showdown.extension("highlight", function () {
  function htmlunencode(text) {
    return text
      .replace(/&amp;/g, "&")
      .replace(/&lt;/g, "<")
      .replace(/&gt;/g, ">");
  }
  return [
    {
      type: "output",
      filter: function (text, converter, options) {
        var left = "<pre><code\\b[^>]*>",
          right = "</code></pre>",
          flags = "g";
        var replacement = function (wholeMatch, match, left, right) {
          match = htmlunencode(match);
          var lang = (left.match(/class=\"([^ \"]+)/) || [])[1];
          left = left.slice(0, 18) + "hljs " + left.slice(18);
          if (lang && hljs.getLanguage(lang)) {
            return (
              left + hljs.highlight(match, { language: lang }).value + right
            );
          } else {
            return left + hljs.highlightAuto(match).value + right;
          }
        };
        return showdown.helper.replaceRecursiveRegExp(
          text,
          replacement,
          left,
          right,
          flags
        );
      },
    },
  ];
});

async function run() {
  return new Promise(async (resolve, reject) => {
    const styleData = await fs.readFile(__dirname + "/style.css", "utf-8");
    const highlightingStyles = await fs.readFile(
      __dirname + "/node_modules/highlight.js/styles/atom-one-dark.css",
      "utf-8"
    );
    if (!argv.in || !argv.out) {
      throw new Error("Missing in and/out out arguments");
    }
    const dir = argv.in;
    const exists = await fs
      .stat(argv.out)
      .then((stat) => stat.isDirectory())
      .catch(() => false);
    if (!exists) {
      await fs.mkdir(argv.out, { recursive: true });
    }
    const events = walkdir(dir);
    events.on("file", async (infile) => {
      events.pause();
      try {
        const rel = path.relative(dir, infile);
        if (!/\.(md|markdown)$/.test(infile)) {
          const outfile = path.join(argv.out, rel);
          await fs.copyFile(infile, outfile);
          return;
        }
        const outfile = path
          .join(argv.out, rel)
          .replace(/\.(md|markdown)$/, ".html");

        const text = await fs.readFile(infile, "utf-8");
        const { data, content } = fm(text);
        converter = new showdown.Converter({
          ghCompatibleHeaderId: true,
          simpleLineBreaks: true,
          ghMentions: true,
          extensions: ["highlight"],
          tables: true,
        });

        var preContent =
          `
        <html>
          <head>
            <title>` +
          (data?.title || rel || "") +
          `</title>
            <meta name="viewport" content="width=device-width, initial-scale=1">
            <meta charset="UTF-8">`;

        preContent += `
          </head>
          <body>
            <div id='content'>
        `;

        let postContent =
          `
  
            </div>
            <style type='text/css'>` +
          styleData +
          `</style>
            <style type='text/css'>` +
          highlightingStyles +
          `</style>
          </body>
        </html>`;

        html = preContent + converter.makeHtml(content) + postContent;

        converter.setFlavor("github");
        // console.log(html);
        await fs.writeFile(outfile, html);
        console.log("Done, saved to " + outfile);
      } catch (err) {
        reject(err);
        throw err;
      } finally {
        events.resume();
      }
    });
    events.on("error", (err) => {
      reject(err);
    });
    events.on("end", () => {
      resolve();
    });
  });
}

if (require.main === module) {
  run()
    .then(() => {
      console.log("Done");
      process.exit();
    })
    .catch((err) => console.error(err));
}
