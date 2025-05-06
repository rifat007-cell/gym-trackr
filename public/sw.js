const urlsToCache = [
  "/",
  "app.js",
  "styles.css",
  "fonts",
  "components",
  "services",
  "bg.png",
  "https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:opsz,wght,FILL,GRAD@20..48,100..700,0..1,-50..200",
  "https://unpkg.com/pwacompat",
  "https://cdn.jsdelivr.net/npm/jwt-decode@4.0.0/build/cjs/index.min.js",
  "https://unpkg.com/@simplewebauthn/browser/dist/bundle/index.umd.min.js",
  "https://cdn.jsdelivr.net/npm/chart.js",
];
self.addEventListener("install", (event) => {
  let cacheUrls = async () => {
    const cache = await caches.open("gymtrack-v1");
    return cache.addAll(urlsToCache);
  };
  event.waitUntil(cacheUrls());
});

self.addEventListener("fetch", (event) => {
  event.respondWith(
    (async () => {
      const cache = await caches.open("gymtrack-v1");

      // from the cache;

      const cachedResponse = await cache.match(event.request);

      // Fetch the latest resource from the network
      const fetchPromise = fetch(event.request)
        .then((networkResponse) => {
          // Update the cache with the latest version
          cache.put(event.request, networkResponse.clone());
          return networkResponse;
        })
        .catch(() => cachedResponse); // In case of network failure, use cached response

      // return cached immediately and update cache in the background
      return cachedResponse || fetchPromise;
    })()
  );
});
