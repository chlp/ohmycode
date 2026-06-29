import { test, expect } from '@playwright/test';

test('file name element is visible and editable', async ({ page }) => {
  await page.goto('/');
  await page.waitForLoadState('networkidle');

  const fileNameEl = page.locator('#file-name');
  await expect(fileNameEl).toBeVisible({ timeout: 5000 });

  const ce = await fileNameEl.getAttribute('contenteditable');
  expect(ce).toBe('true');
});

test('clicking file name focuses it', async ({ page }) => {
  await page.goto('/');
  await page.waitForLoadState('networkidle');

  const fileNameEl = page.locator('#file-name');
  await fileNameEl.click();

  // Should be focused — verify by checking activeElement
  const isFocused = await page.evaluate(() => {
    return document.activeElement === document.getElementById('file-name');
  });
  expect(isFocused).toBe(true);
});

test('file name has non-empty text after WS sync', async ({ page }) => {
  await page.goto('/');
  await page.waitForLoadState('networkidle');

  // Wait for WS init round-trip to populate the file name
  const fileNameEl = page.locator('#file-name');
  await expect(fileNameEl).not.toBeEmpty({ timeout: 5000 });
});

test('editing file name changes its text', async ({ page }) => {
  await page.goto('/');
  await page.waitForLoadState('networkidle');

  const fileNameEl = page.locator('#file-name');
  await fileNameEl.click();
  await page.keyboard.selectAll();
  await page.keyboard.type('My Renamed File');

  await expect(fileNameEl).toContainText('My Renamed File');
});
