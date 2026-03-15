import { mount } from "@vue/test-utils";
import DataTable from "./DataTable.vue";

describe("DataTable", () => {
  it("renders rows and emits page changes", async () => {
    const wrapper = mount(DataTable, {
      props: {
        columns: [{ key: "name", label: "Nom" }],
        rows: [{ id: "1", name: "Alice" }],
        pagination: {
          page: 1,
          totalPages: 2,
          total: 2
        }
      }
    });

    expect(wrapper.text()).toContain("Alice");
    await wrapper.findAll("button")[1].trigger("click");
    expect(wrapper.emitted("page-change")[0]).toEqual([2]);
  });
});
