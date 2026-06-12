import http from "node:http"
import fs from "node:fs"
import path from "node:path"
import { fileURLToPath } from "node:url"

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const clientDir = path.join(__dirname, "dist", "client")

const MIME = {
  ".js": "text/javascript",
  ".css": "text/css",
  ".html": "text/html",
  ".json": "application/json",
  ".png": "image/png",
  ".jpg": "image/jpeg",
  ".jpeg": "image/jpeg",
  ".gif": "image/gif",
  ".svg": "image/svg+xml",
  ".ico": "image/x-icon",
  ".webp": "image/webp",
  ".woff": "font/woff",
  ".woff2": "font/woff2",
  ".map": "application/json",
}

const { default: serverEntry } = await import("./dist/server/server.js")
const { fetch: handler } = serverEntry

const server = http.createServer(async (req, res) => {
  try {
    const url = new URL(req.url ?? "/", `http://${req.headers.host ?? "localhost"}`)
    const filePath = path.join(clientDir, url.pathname === "/" ? "index.html" : url.pathname)

    if (url.pathname.startsWith("/assets/") && fs.existsSync(filePath)) {
      const ext = path.extname(filePath)
      res.writeHead(200, { "Content-Type": MIME[ext] ?? "application/octet-stream" })
      fs.createReadStream(filePath).pipe(res)
      return
    }

    const request = new Request(url, {
      method: req.method,
      headers: req.headers,
      body: req.method !== "GET" && req.method !== "HEAD" ? req : undefined,
    })
    const response = await handler(request)
    const body = await response.text()
    res.writeHead(response.status, Object.fromEntries(response.headers))
    res.end(body)
  } catch (err) {
    res.writeHead(500)
    res.end("Internal Server Error")
  }
})

const port = process.env.PORT ? Number(process.env.PORT) : 3000
server.listen(port, () => {
  console.log(`Front SSR server running on http://localhost:${port}`)
})
