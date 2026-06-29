import { test, expect } from '@playwright/test';

test('language selector is visible after load', async ({ page }) => {
  await page.goto('/');
  await page.waitForLoadState('networkidle');
  await expect(page.locator('#lang-select')).toBeVisible({ timeout: 5000 });
});

test('language selector contains expected languages', async ({ page }) => {
  await page.goto('/');
  await page.waitForLoadState('networkidle');

  const select = page.locator('#lang-select');
  await expect(select).toBeVisible();

  const options = await select.locator('option').allTextContents();
  const langs = options.map(o => o.toLowerCase());

  expect(langs.some(l => l.includes('python'))).toBe(true);
  expect(langs.some(l => l.includes('go'))).toBe(true);
  expect(langs.some(l => l.includes('node'))).toBe(true);
  expect(langs.some(l => l.includes('java'))).toBe(true);
  expect(langs.some(l => l.includes('markdown'))).toBe(true);
});

test('selecting a language updates the selector value', async ({ page }) => {
  await page.goto('/');
  await page.waitForLoadState('networkidle');

  const select = page.locator('#lang-select');
  await expect(select).toBeVisible();

  // Switch to GoLang
  await select.selectOption('go');
  await expect(select).toHaveValue('go');
});

test('language persists across selector changes', async ({ page }) => {
  await page.goto('/');
  await page.waitForLoadState('networkidle');

  const select = page.locator('#lang-select');
  await select.selectOption('python3');
  await select.selectOption('nodejs');
  await expect(select).toHaveValue('nodejs');
});
