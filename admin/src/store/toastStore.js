import { reactive, readonly } from "vue";

const state = reactive({
  items: []
});

let nextToastId = 1;

export function useToastStore() {
  function pushToast({ title, message, tone = "green", timeout = 3200 }) {
    const id = nextToastId++;
    state.items.push({ id, title, message, tone });

    if (timeout > 0) {
      window.setTimeout(() => dismissToast(id), timeout);
    }
  }

  function dismissToast(id) {
    const index = state.items.findIndex((item) => item.id === id);
    if (index >= 0) {
      state.items.splice(index, 1);
    }
  }

  return {
    toasts: readonly(state.items),
    pushToast,
    dismissToast
  };
}
