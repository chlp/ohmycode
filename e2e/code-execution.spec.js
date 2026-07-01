import { test, expect } from '@playwright/test';

// Code execution: run_task_with_content -> runner -> result pane.
// Requires a live app with the python3 runner container up (see e2e/CLAUDE.md).

test('running python code shows its output in the result pane', async ({ page }) => {
  await page.goto('/');
  await page.waitForLoadState('networkidle');

  await page.locator('#lang-select').selectOption('python3');

  const editor = page.locator('#content-container .CodeMirror');
  await editor.click();
  await page.keyboard.press('Control+A');
  await page.keyboard.insertText('print(21 + 21)');

  // Wait for a runner to come online (background sync ticks every ~1s).
  await expect(page.locator('#run-button')).toBeEnabled({ timeout: 15000 });

  await page.locator('#run-button').click();

  const resultCode = page.locator('#result-container .CodeMirror-code');
  await expect(resultCode).toContainText('42', { timeout: 15000 });
});

test('clean result button clears a previous run result', async ({ page }) => {
  await page.goto('/');
  await page.waitForLoadState('networkidle');

  await page.locator('#lang-select').selectOption('python3');

  const editor = page.locator('#content-container .CodeMirror');
  await editor.click();
  await page.keyboard.press('Control+A');
  await page.keyboard.insertText('print("hello-e2e")');

  await expect(page.locator('#run-button')).toBeEnabled({ timeout: 15000 });
  await page.locator('#run-button').click();

  const resultCode = page.locator('#result-container .CodeMirror-code');
  await expect(resultCode).toContainText('hello-e2e', { timeout: 15000 });

  await expect(page.locator('#clean-result-button')).toBeEnabled({ timeout: 5000 });
  await page.locator('#clean-result-button').click();

  await expect(page.locator('#file-result')).toBeHidden({ timeout: 5000 });
});
