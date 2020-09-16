const ASSETS = [
  "index.html",
  "app.js",
  "app.css",
  "func.js",
  "Chart.bundle.min.js",
  "logo.png",
  "offline.html"
];

let cache_name = "Eos_3.0.0";

self.addEventListener("install", event => {
  console.log("installing...");
  event.waitUntil(
    caches
      .open(cache_name)
      .then(cache => {
        return cache.addAll(assets);
      })
      .catch(err => console.log(err))
  );
});

self.addEventListener("fetch", event => {
  if (event.request.url === "https://eos.lyrenhex.com/") {
      event.respondWith(
          fetch(event.request).catch(err =>
              self.cache.open(cache_name).then(cache => cache.match("/offline.html"))
          )
      );
  } else {
      event.respondWith(
          fetch(event.request).catch(err =>
              caches.match(event.request).then(response => response)
          )
      );
  }
});
