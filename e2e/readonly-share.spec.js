import { test, expect } from '@playwright/test';

// Read-only share links: client-side AES-GCM encryption + a read-only key/token,
// verifying the round trip decrypts correctly and the RO viewer can't edit.

test('encrypted read-only share link lets a second client view but not edit', async ({ browser }) => {
  const ctx = await browser.newContext();
  const page = await ctx.newPage();
  await page.goto('/');
  await page.waitForLoadState('networkidle');

  const editor = page.locator('#content-container .CodeMirror');
  await editor.click();
  await page.keyboard.press('Control+A');
  await page.keyboard.insertText('encrypted-marker-42');

  await page.locator('#header-encrypt-btn').click();
  await page.locator('.encrypt-action-btn', { hasText: 'Enable Encryption' }).click();

  const roLinkInput = page
    .locator('.encrypt-link-row', { hasText: 'Read-only link' })
    .locator('input');
  await expect(roLinkInput).toHaveValue(/ro=/, { timeout: 10000 });
  const roLink = await roLinkInput.inputValue();

  await ctx.close();

  const ctx2 = await browser.newContext();
  const page2 = await ctx2.newPage();
  await page2.goto(roLink);
  await page2.waitForLoadState('networkidle');

  const code2 = await page2.locator('#content-container .CodeMirror-code').textContent();
  expect(code2).toContain('encrypted-marker-42');

  await expect(page2.locator('#header-lock-btn')).toBeDisabled();
  const isReadOnly = await page2.evaluate(() => {
    return document.querySelector('#content-container .CodeMirror').CodeMirror.getOption('readOnly');
  });
  expect(isReadOnly).toBeTruthy();

  await ctx2.close();
});
