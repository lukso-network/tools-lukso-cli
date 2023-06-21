const cors = {
  "Access-Control-Allow-Origin": "*",
  "Access-Control-Allow-Methods": "GET,HEAD,OPTIONS",
  "Access-Control-Max-Age": "60",
  "Access-Control-Allow-Headers": "content-type",
  "Cache-Control": "max-age=60",
};

export default {
  async fetch(request) {
    const url = new URL(request.url);
    const match = /^\/?(\d*|preview|l16)?\/?$/.exec(url.pathname);
    if (match) {
      const pr = match[1];
      const url = `https://storage.googleapis.com/lks-lz-binaries-euw4${
        pr ? `/${pr}/install.sh` : "/install.sh"
      }`;
      const proxyRequest = await fetch(url);
      if (proxyRequest.status !== 200) {
        console.error("error", proxyRequest, url);
        return new Response(proxyRequest.statusTest, {
          status: proxyRequest.status,
          headers: cors,
        });
      }
      const output = await proxyRequest.text();
      return new Response(output, {
        headers: {
          ...proxyRequest.headers,
          "content-type": "text/plain",
          ...cors,
        },
      });
    }
    return new Response("File not found", {
      status: 404,
      headers: cors,
    });
  },
};
