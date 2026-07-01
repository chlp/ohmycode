import { test, expect } from '@playwright/test';

// File locking: header-lock-btn -> set_locked -> read-only editor for everyone viewing the file.

test('lock button locks the file and makes the editor read-only', async ({ page }) => {
  await page.goto('/');
  await page.waitForLoadState('networkidle');

  await page.locator('#header-lock-btn').click();
  await expect(page.locator('#status-bar')).toHaveText('Locked', { timeout: 5000 });

  const isReadOnly = await page.evaluate(() => {
    return document.querySelector('#content-container .CodeMirror').CodeMirror.getOption('readOnly');
  });
  expect(isReadOnly).toBeTruthy();

  // Unlock again so the run stays hygienic if this file id gets reopened elsewhere.
  await page.locator('#header-lock-btn').click();
  await expect(page.locator('#status-bar')).not.toHaveText('Locked', { timeout: 5000 });
});

test('lock set by one client is reflected on another client viewing the same file', async ({ browser }) => {
  const ctxA = await browser.newContext();
  const ctxB = await browser.newContext();

  const pageA = await ctxA.newPage();
  await pageA.goto('/');
  await pageA.waitForLoadState('networkidle');
  const fileUrl = pageA.url();

  const pageB = await ctxB.newPage();
  await pageB.goto(fileUrl);
  await pageB.waitForLoadState('networkidle');

  await pageA.locator('#header-lock-btn').click();

  await expect(pageB.locator('#status-bar')).toHaveText('Locked', { timeout: 5000 });
  await expect(pageB.locator('#header-lock-btn')).toHaveAttribute('title', 'Unlock editing');
  const isReadOnlyOnB = await pageB.evaluate(() => {
    return document.querySelector('#content-container .CodeMirror').CodeMirror.getOption('readOnly');
  });
  expect(isReadOnlyOnB).toBeTruthy();

  await ctxA.close();
  await ctxB.close();
});
