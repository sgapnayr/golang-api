import axios from "axios";

const BASE_URL = "http://localhost:8080";

async function makeRequests() {
  // Simulate a GET request
  await axios.get(`${BASE_URL}/orders`);

  // Simulate a POST request
  await axios.post(`${BASE_URL}/orders`, {
    id: "3",
    item: "Item 3",
    amount: 30,
  });
}

async function stressTest(numUsers: number): Promise<boolean> {
  const promises = [];
  const startTime = Date.now();

  for (let i = 0; i < numUsers; i++) {
    promises.push(makeRequests().catch(() => false));
  }

  const results = await Promise.all(promises);
  const duration = Date.now() - startTime;
  console.log(`Stress test with ${numUsers} users took ${duration}ms`);

  return results.every((result) => result !== false);
}

describe("Stress Test", () => {
  const tiers = [1, 10, 100, 1000];

  tiers.forEach((numUsers) => {
    it(`should handle ${numUsers} users`, async () => {
      const success = await stressTest(numUsers);
      console.log(`Test with ${numUsers} users: ${success ? "PASS" : "FAIL"}`);
      expect(success).toBe(true);
    });
  });
});
