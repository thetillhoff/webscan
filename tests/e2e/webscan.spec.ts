import { expect, test } from "@playwright/test";

test("scan thetillhoff.de and show successful keyword-based output", async ({ page }) => {
  await page.goto("/");

  await page.locator("#targetInput").fill("thetillhoff.de");
  await page.locator("#scanButton").click();

  await expect
    .poll(
      async () => (await page.locator("#scanStatus").innerText()).trim(),
      {
        timeout: 120_000,
      },
    )
    .toMatch(/Scan completed successfully|Scan failed|Scan error/i);

  const statusText = (await page.locator("#scanStatus").innerText()).trim();
  expect(statusText).toMatch(/Scan completed successfully/i);

  const resultsText = await page.locator("#scanResults").innerText();
  expect(resultsText).toContain("# webscan results for thetillhoff.de");
  expect(resultsText).toContain("## DNS scan results");
  expect(resultsText).not.toContain("Error:");
});
