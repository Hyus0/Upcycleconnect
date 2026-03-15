import { adminApi } from "./api";

describe("adminApi", () => {
  it("exposes local CRUD capabilities", () => {
    const capabilities = adminApi.getCapabilities();
    expect(capabilities.categories.list).toBe(true);
    expect(capabilities.prestations.create).toBe(true);
    expect(capabilities.users.delete).toBe(true);
  });

  it("returns a portal snapshot structure", async () => {
    global.fetch = vi.fn().mockResolvedValue({
      json: async () => ({
        users: [],
        prestations: [],
        categories: [],
        events: []
      })
    });

    const snapshot = await adminApi.getPortalSnapshot();
    expect(snapshot).toHaveProperty("particulier");
    expect(snapshot).toHaveProperty("prestataire");
    expect(snapshot).toHaveProperty("salarie");
  });
});
