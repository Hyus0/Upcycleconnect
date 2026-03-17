import { adminApi } from "./api";

describe("adminApi", () => {
  it("exposes local CRUD capabilities", () => {
    const capabilities = adminApi.getCapabilities();
    expect(capabilities.categories.list).toBe(true);
    expect(capabilities.prestations.create).toBe(true);
    expect(capabilities.users.delete).toBe(true);
  });
});
