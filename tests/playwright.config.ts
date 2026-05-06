import { defineConfig } from "@playwright/test";

export default defineConfig({
  testDir: ".",
  workers: 1,
  timeout: 120_000,
  expect: {
    timeout: 30_000,
  },
  use: {
    headless: true,
    baseURL: "http://127.0.0.1:4173",
  },
  webServer: {
    command: "bash ./start-test-stack.sh",
    url: "http://127.0.0.1:4173/api/health",
    reuseExistingServer: true,
    timeout: 120_000,
  },
});
