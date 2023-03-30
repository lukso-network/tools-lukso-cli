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
      const proxyRequest = await fetch(`https://storage.googleapis.com/lks-lz-binaries-euw4${pr ? `/${pr}/install.sh` : '/install.sh'}`);
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
