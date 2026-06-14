import { expect, test } from "@playwright/test";

test("scan thetillhoff.de and show successful keyword-based output", async ({ page }) => {
  await page.goto("/scan?q=thetillhoff.de");

  await expect(page.locator("#resultsSection")).toBeVisible({ timeout: 120_000 });

  const resultsText = await page.locator("#scanResults").innerText();
  expect(resultsText).toContain("# webscan results for thetillhoff.de");
  expect(resultsText).toContain("## DNS scan results");
  expect(resultsText).not.toContain("Error:");
});
