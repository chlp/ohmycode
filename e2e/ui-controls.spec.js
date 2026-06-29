import { test, expect } from '@playwright/test';

// UI controls: buttons, status bar, header icons

test('history button is visible in file header', async ({ page }) => {
  await page.goto('/');
  await page.waitForLoadState('networkidle');
  await expect(page.locator('#header-history-btn')).toBeVisible({ timeout: 5000 });
});

test('lock button is visible in file header', async ({ page }) => {
  await page.goto('/');
  await page.waitForLoadState('networkidle');
  await expect(page.locator('#header-lock-btn')).toBeVisible({ timeout: 5000 });
});

test('encrypt button is visible in file header', async ({ page }) => {
  await page.goto('/');
  await page.waitForLoadState('networkidle');
  await expect(page.locator('#header-encrypt-btn')).toBeVisible({ timeout: 5000 });
});

test('download button is visible in file header', async ({ page }) => {
  await page.goto('/');
  await page.waitForLoadState('networkidle');
  await expect(page.locator('#header-download-btn')).toBeVisible({ timeout: 5000 });
});

test('controls container becomes visible after WS init', async ({ page }) => {
  await page.goto('/');
  await page.waitForLoadState('networkidle');
  // controls-container starts as display:none but JS shows it after file is received
  await expect(page.locator('#controls-container')).toBeVisible({ timeout: 5000 });
});

test('status bar becomes visible after WS init', async ({ page }) => {
  await page.goto('/');
  await page.waitForLoadState('networkidle');
  await expect(page.locator('#status-bar')).toBeVisible({ timeout: 5000 });
});

test('new file button creates a new file URL', async ({ page }) => {
  await page.goto('/');
  await page.waitForLoadState('networkidle');
  const firstUrl = page.url();

  await page.locator('#header-new-file-btn').click();
  await page.waitForLoadState('networkidle');
  const newUrl = page.url();

  expect(newUrl).not.toBe(firstUrl);
  const path = new URL(newUrl).pathname.slice(1);
  expect(path).toMatch(/^[0-9A-Za-z]{22}$/);
});

test('versions button is visible', async ({ page }) => {
  await page.goto('/');
  await page.waitForLoadState('networkidle');
  await expect(page.locator('#header-versions-btn')).toBeVisible({ timeout: 5000 });
});
