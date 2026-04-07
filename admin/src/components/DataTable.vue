<template>
  <div class="surface-card table-shell">
    <div class="table-wrap">
      <table>
        <thead>
          <tr>
            <th
              v-for="column in columns"
              :key="column.key"
              :class="column.align ? `align-${column.align}` : ''"
            >
              {{ column.label }}
            </th>
            <th v-if="$slots.actions" class="align-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in rows" :key="row[rowKey]">
            <td
              v-for="column in columns"
              :key="column.key"
              :class="column.align ? `align-${column.align}` : ''"
            >
              <slot :name="`cell-${column.key}`" :row="row">
                {{ row[column.key] }}
              </slot>
            </td>
            <td v-if="$slots.actions" class="align-right actions-cell">
              <slot name="actions" :row="row" />
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="pagination" class="table-footer">
      <span>Page {{ pagination.page }} / {{ pagination.totalPages }} - {{ pagination.total }} elements</span>
      <div class="toolbar">
        <button
          class="button pagination-button"
          :disabled="pagination.page <= 1"
          @click="$emit('page-change', pagination.page - 1)"
        >
          Precedent
        </button>
        <button
          class="button pagination-button"
          :disabled="pagination.page >= pagination.totalPages"
          @click="$emit('page-change', pagination.page + 1)"
        >
          Suivant
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  columns: {
    type: Array,
    required: true
  },
  rows: {
    type: Array,
    required: true
  },
  rowKey: {
    type: String,
    default: "id"
  },
  pagination: {
    type: Object,
    default: null
  }
});

defineEmits(["page-change"]);
</script>

<style scoped>
.table-shell {
  overflow: hidden;
  border-radius: 26px;
}

.table-wrap {
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th,
td {
  padding: 16px 18px;
  border-bottom: 1px solid rgba(118, 148, 132, 0.18);
  vertical-align: middle;
}

th {
  text-align: left;
  font-size: 0.82rem;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--text-muted);
  font-family: "Space Mono", monospace;
}

tbody tr:hover {
  background: linear-gradient(90deg, rgba(45, 122, 79, 0.045), rgba(45, 122, 79, 0.02));
}

.align-right {
  text-align: right;
}

.align-center {
  text-align: center;
}

.actions-cell :deep(.button) {
  padding: 9px 12px;
}

.table-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  padding: 16px 18px;
  background:
    linear-gradient(180deg, rgba(18, 29, 24, 0.96), rgba(11, 20, 16, 0.98));
  border-top: 1px solid rgba(118, 148, 132, 0.18);
  color: var(--text-secondary);
}

.pagination-button {
  background:
    linear-gradient(180deg, rgba(43, 113, 74, 0.3), rgba(26, 70, 46, 0.32));
  color: #f3fbf6;
  border: 1px solid rgba(98, 196, 136, 0.28);
}

.pagination-button:disabled {
  background: rgba(255, 255, 255, 0.04);
  border-color: rgba(255, 255, 255, 0.06);
  color: rgba(199, 214, 205, 0.45);
}

@media (max-width: 700px) {
  .table-footer {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
