import { expect, test } from "@playwright/test";

test("health endpoint returns healthy status", async ({ request }) => {
  const response = await request.get("/api/health");
  expect(response.status()).toBe(200);

  const body = await response.json();
  expect(body).toEqual({ status: "healthy" });
});

test("scan endpoint rejects empty target", async ({ request }) => {
  const response = await request.post("/api/scan", {
    data: { target: "" },
  });
  expect(response.status()).toBe(400);

  const body = await response.json();
  expect(body.error).toMatch(/invalid request|target cannot be empty/i);
});

test("scan endpoint succeeds for thetillhoff.de", async ({ request }) => {
  const enqueue = await request.post("/api/scan", {
    data: { target: "thetillhoff.de", follow: false },
  });
  expect(enqueue.status()).toBe(202);

  const queued = await enqueue.json();
  expect(queued.status).toBe("queued");
  expect(typeof queued.job_id).toBe("string");
  expect(queued.job_id.length).toBeGreaterThan(0);

  let final: any;
  await expect
    .poll(
      async () => {
        const statusResponse = await request.get(`/api/scan/${queued.job_id}`, {
          timeout: 120_000,
        });
        expect(statusResponse.status()).toBe(200);
        final = await statusResponse.json();
        return final.status;
      },
      {
        timeout: 180_000,
      },
    )
    .toBe("completed");

  expect(final.job_id).toBe(queued.job_id);
  expect(final.target).toBe("thetillhoff.de");

  expect(final.results).toContain("# webscan results for thetillhoff.de");
  expect(final.results).toContain("## DNS scan results");
  expect(final.results).not.toContain("Error:");
  expect(typeof final.stderr).toBe("string");
  expect(final.stderr.length).toBeGreaterThan(0);
});
