const cors = {
  "Access-Control-Allow-Origin": "*",
  "Access-Control-Allow-Methods": "GET,HEAD,OPTIONS",
  "Access-Control-Max-Age": "86400",
  "Access-Control-Allow-Headers": "content-type",
};

export default {
  async fetch(request) {
    const url = new URL(request.url);
    const match = /^\/?(l?\d*)?\/?$/.exec(url.pathname);
    if (match) {
      const pr = match[1];
      url.pathname = pr
        ? `/tools-lukso-cli/pr-preview/pr-${pr}/index.html`
        : `/tools-lukso-cli/index.html`;
      url.hostname = "lukso-network.github.io";
      const fetchUrl = url.toString();
      const proxyRequest = await fetch(fetchUrl);
      if (proxyRequest.status !== 200) {
        return new Response(proxyRequest.statusTest, {
          status: proxyRequest.status,
          headers: cors,
        });
      }
      const output = await proxyRequest.text();
      return new Response(output, {
        headers: { "content-type": "text/plain", ...cors },
      });
    }
    return new Response("File not found", {
      status: 404,
      headers: cors,
    });
  },
};
