import { test, expect } from '@playwright/test';

// Verifies the hash-based cache busting mechanism end-to-end.
//
// Hash versioning only applies when use_dynamic_files=false (embedded/prod mode).
// In dev mode (use_dynamic_files=true, volume-mounted files), ?v= stays as the
// hardcoded numeric value from index.html — those tests are skipped automatically.

let embeddedMode = null;
const isEmbeddedMode = async (request) => {
  if (embeddedMode !== null) return embeddedMode;
  const body = await (await request.get('/')).text();
  // Embedded mode rewrites ?v=<number> to ?v=<hexhash>; dynamic mode leaves numerics.
  embeddedMode = !/\?v=\d+"/.test(body);
  return embeddedMode;
};

test('index.html is served with no-cache', async ({ request }) => {
  const resp = await request.get('/');
  expect(resp.status()).toBe(200);
  expect(resp.headers()['cache-control']).toContain('no-cache');
});

test('index.html contains versioned asset references', async ({ request }) => {
  if (!await isEmbeddedMode(request)) {
    test.skip(true, 'hash versioning only applies in embedded mode (use_dynamic_files=false)');
  }
  const resp = await request.get('/');
  const body = await resp.text();
  // All ?v= references should be a short hex hash, not the old hardcoded numbers
  expect(body).toContain('?v=');
  expect(body).not.toMatch(/\?v=\d+"/); // no bare numeric versions like ?v=18"
});

test('main.js is served with versioned module imports', async ({ request }) => {
  if (!await isEmbeddedMode(request)) {
    test.skip(true, 'hash versioning only applies in embedded mode (use_dynamic_files=false)');
  }
  const resp = await request.get('/js/main.js');
  expect(resp.status()).toBe(200);
  const body = await resp.text();
  // Every relative import should be versioned
  const lines = body.split('\n').filter(l => l.includes('"./') && l.includes('.js'));
  for (const line of lines) {
    expect(line).toContain('?v=');
  }
});

test('versioned JS asset has immutable cache headers', async ({ request }) => {
  if (!await isEmbeddedMode(request)) {
    test.skip(true, 'hash versioning only applies in embedded mode (use_dynamic_files=false)');
  }
  // Extract the build hash from index.html
  const indexResp = await request.get('/');
  const indexBody = await indexResp.text();
  const match = indexBody.match(/main\.js\?v=([a-f0-9]+)/);
  expect(match).not.toBeNull();
  const hash = match[1];

  const jsResp = await request.get(`/js/main.js?v=${hash}`);
  expect(jsResp.status()).toBe(200);
  const cc = jsResp.headers()['cache-control'];
  expect(cc).toContain('immutable');
  expect(cc).toContain('max-age=31536000');
});

test('all JS modules loaded by the browser have versioned URLs', async ({ page, request }) => {
  if (!await isEmbeddedMode(request)) {
    test.skip(true, 'hash versioning only applies in embedded mode (use_dynamic_files=false)');
  }
  const jsUrls = new Set();
  page.on('request', req => {
    const url = req.url();
    if (url.includes('/js/') && url.includes('.js')) {
      jsUrls.add(url);
    }
  });

  await page.goto('/');
  await page.waitForLoadState('networkidle');

  expect(jsUrls.size).toBeGreaterThan(0);
  for (const url of jsUrls) {
    expect(url, `unversioned JS module: ${url}`).toContain('?v=');
  }
});

test('build hash is consistent across requests', async ({ request }) => {
  if (!await isEmbeddedMode(request)) {
    test.skip(true, 'hash versioning only applies in embedded mode (use_dynamic_files=false)');
  }
  const extractHash = async () => {
    const body = await (await request.get('/')).text();
    return body.match(/main\.js\?v=([a-f0-9]+)/)?.[1];
  };
  const h1 = await extractHash();
  const h2 = await extractHash();
  expect(h1).toBe(h2);
});
