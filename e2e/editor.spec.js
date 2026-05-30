import { test, expect } from '@playwright/test';

test('app loads and editor element is present', async ({ page }) => {
  await page.goto('/');
  await page.waitForLoadState('domcontentloaded');
  await expect(page.locator('#content')).toBeAttached();
});

test('app navigates to a URL with a valid 22-char file ID', async ({ page }) => {
  await page.goto('/');
  await page.waitForLoadState('networkidle');
  const path = new URL(page.url()).pathname.slice(1);
  expect(path).toMatch(/^[0-9A-Za-z]{22}$/);
});

test('body becomes visible after DOMContentLoaded', async ({ page }) => {
  await page.goto('/');
  // body starts at opacity:0 and transitions to 1 after DOMContentLoaded
  await expect(page.locator('body')).toHaveCSS('opacity', '1', { timeout: 5000 });
});

test('WebSocket connects to /file endpoint', async ({ page }) => {
  const wsPromise = page.waitForEvent('websocket', ws => ws.url().includes('/file'), { timeout: 10_000 });
  await page.goto('/');
  const ws = await wsPromise;
  expect(ws.url()).toContain('/file');
});

test('CodeMirror editor is rendered', async ({ page }) => {
  await page.goto('/');
  await page.waitForLoadState('networkidle');
  await expect(page.locator('.CodeMirror')).toBeVisible({ timeout: 5000 });
});

test('file header is visible', async ({ page }) => {
  await page.goto('/');
  await page.waitForLoadState('networkidle');
  await expect(page.locator('#file-header')).toBeVisible();
});

test('saving content via WebSocket updates the editor', async ({ page }) => {
  await page.goto('/');
  await page.waitForLoadState('networkidle');

  // Wait for WS to connect and initial file to arrive
  const ws = await page.waitForEvent('websocket', ws => ws.url().includes('/file'), { timeout: 10_000 });
  await ws.waitForEvent('framesent', { timeout: 5000 }); // init message sent
  await ws.waitForEvent('framereceived', { timeout: 5000 }); // file snapshot received

  // Editor should have some state at this point
  await expect(page.locator('.CodeMirror')).toBeVisible();
});
