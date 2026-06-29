import { test, expect } from '@playwright/test';

// Tests requiring a live app at APP_URL (default: http://localhost:52674).
// Two browser contexts open the same file URL and verify shared state.

test('two clients open the same file and see the same lang', async ({ browser }) => {
  const ctxA = await browser.newContext();
  const ctxB = await browser.newContext();

  const pageA = await ctxA.newPage();
  await pageA.goto('/');
  await pageA.waitForLoadState('networkidle');

  const fileUrl = pageA.url();

  const pageB = await ctxB.newPage();
  await pageB.goto(fileUrl);
  await pageB.waitForLoadState('networkidle');

  const langA = await pageA.locator('#lang-select').inputValue();
  const langB = await pageB.locator('#lang-select').inputValue();

  expect(langA).toBe(langB);

  await ctxA.close();
  await ctxB.close();
});

test('content typed in one client appears in the other', async ({ browser }) => {
  const ctxA = await browser.newContext();
  const ctxB = await browser.newContext();

  const pageA = await ctxA.newPage();
  await pageA.goto('/');
  await pageA.waitForLoadState('networkidle');
  const fileUrl = pageA.url();

  const pageB = await ctxB.newPage();
  await pageB.goto(fileUrl);
  await pageB.waitForLoadState('networkidle');

  // Give both clients time to subscribe
  await pageA.waitForTimeout(500);

  // Type in A's CodeMirror
  const editorA = pageA.locator('.CodeMirror').first();
  await editorA.click();
  await pageA.keyboard.type('collab sync test');

  // Wait for WS propagation (throttle 500ms + network)
  await pageB.waitForTimeout(2000);

  const codeB = await pageB.locator('.CodeMirror-code').textContent();
  expect(codeB).toContain('collab sync test');

  await ctxA.close();
  await ctxB.close();
});

test('file name is the same across two clients on the same file', async ({ browser }) => {
  const ctxA = await browser.newContext();
  const ctxB = await browser.newContext();

  const pageA = await ctxA.newPage();
  await pageA.goto('/');
  await pageA.waitForLoadState('networkidle');
  const fileUrl = pageA.url();

  const pageB = await ctxB.newPage();
  await pageB.goto(fileUrl);
  await pageB.waitForLoadState('networkidle');

  const nameA = await pageA.locator('#file-name').textContent();
  const nameB = await pageB.locator('#file-name').textContent();

  expect(nameA).toBe(nameB);

  await ctxA.close();
  await ctxB.close();
});
