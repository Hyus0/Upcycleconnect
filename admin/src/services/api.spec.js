import { adminApi } from "./api";
import { MissingEndpointError } from "./http";

describe("adminApi", () => {
  it("exposes capabilities for missing resources", () => {
    const capabilities = adminApi.getCapabilities();
    expect(capabilities.categories.list).toBe(false);
    expect(capabilities.prestations.create).toBe(true);
  });

  it("throws a MissingEndpointError for category listing", async () => {
    await expect(adminApi.listCategories()).rejects.toBeInstanceOf(MissingEndpointError);
  });
});
