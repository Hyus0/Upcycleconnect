<template>
  <section class="stack">
    <header class="page-header">
      <div>
        <div class="eyebrow">Finances</div>
        <h2 class="page-title">Pilotage financier</h2>
        <p class="page-description">Suivi des flux, abonnements et montants en attente.</p>
      </div>
      <button class="button button-secondary" @click="loadFinance">Actualiser</button>
    </header>

    <div class="stats-grid finance-grid">
      <article class="surface-card section-card stack finance-stat" v-for="item in financeStats" :key="item.label">
        <div class="eyebrow">{{ item.label }}</div>
        <strong class="stat-value">{{ item.value }}</strong>
      </article>
    </div>

    <article class="surface-card section-card stack command-panel">
      <div class="panel-head panel-head--compact">
        <div>
          <span class="panel-kicker">Flux</span>
          <h3>Transactions</h3>
        </div>
        <span class="mini-note">{{ pagination?.total ?? 0 }} lignes</span>
      </div>
      <div class="quick-chips">
        <button
          v-for="option in statusOptions"
          :key="option.value || 'all'"
          class="quick-chip"
          :class="{ active: filters.status === option.value }"
          type="button"
          @click="filters.status = option.value"
        >
          {{ option.label }}
        </button>
      </div>
      <div class="filters-grid command-filters">
        <FormField label="Categorie">
          <BaseSelect v-model="filters.category" :options="categoryOptions" />
        </FormField>
      </div>
      <button class="button button-secondary filters-reset" type="button" @click="resetFilters">
        Reinitialiser
      </button>
    </article>

    <LoadingState v-if="loading" />
    <ErrorState v-else-if="error" :message="error" retry-label="Recharger" @retry="loadFinance" />
    <DataTable v-else :columns="columns" :rows="rows" :pagination="pagination" @page-change="changePage">
      <template #cell-amount="{ row }">
        <strong>{{ formatCurrency(row.amount) }}</strong>
      </template>
      <template #cell-status="{ row }">
        <StatusBadge :label="row.status" :tone="row.status === 'paid' ? 'green' : 'amber'" />
      </template>
    </DataTable>
  </section>
</template>

<script setup>
import { computed, onMounted, reactive, ref, watch } from "vue";
import BaseSelect from "../components/BaseSelect.vue";
import DataTable from "../components/DataTable.vue";
import ErrorState from "../components/ErrorState.vue";
import FormField from "../components/FormField.vue";
import LoadingState from "../components/LoadingState.vue";
import StatusBadge from "../components/StatusBadge.vue";
import { adminApi } from "../services/api";
import { formatCurrency } from "../utils/format";

const loading = ref(true);
const error = ref("");
const rows = ref([]);
const pagination = ref(null);
const summary = ref({});
const filters = reactive({ category: "", status: "", page: 1, pageSize: 8 });

const columns = [
  { key: "label", label: "Libelle" },
  { key: "category", label: "Categorie" },
  { key: "amount", label: "Montant", align: "right" },
  { key: "status", label: "Statut" },
  { key: "dueDate", label: "Echeance" }
];

const financeStats = computed(() => [
  { label: "Chiffre encaisse", value: formatCurrency(summary.value.paidTotal ?? 0) },
  { label: "En attente", value: formatCurrency(summary.value.pendingTotal ?? 0) },
  { label: "Abonnements actifs", value: String(summary.value.activeSubscriptions ?? 0) }
]);

const categoryOptions = [
  { label: "Toutes", value: "" },
  { label: "Abonnement", value: "subscription" },
  { label: "Prestation", value: "service" },
  { label: "Evenement", value: "event" }
];

const statusOptions = [
  { label: "Tous", value: "" },
  { label: "En attente", value: "pending" },
  { label: "Paye", value: "paid" },
  { label: "Echoue", value: "failed" }
];

async function loadFinance() {
  loading.value = true;
  error.value = "";
  try {
    const response = await adminApi.getFinanceOverview(filters);
    summary.value = response.summary;
    rows.value = response.items;
    pagination.value = response.pagination;
  } catch (err) {
    error.value = err.message ?? "Impossible de charger les finances.";
  } finally {
    loading.value = false;
  }
}

function changePage(page) {
  filters.page = page;
}

function resetFilters() {
  filters.category = "";
  filters.status = "";
  filters.page = 1;
}

watch(() => [filters.category, filters.status, filters.page], loadFinance);
onMounted(loadFinance);
</script>

<style scoped>
.stat-value {
  font-size: 1.8rem;
  font-family: "Syne", sans-serif;
}

.finance-grid {
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
}

.finance-stat {
  background:
    radial-gradient(circle at 88% 0%, rgba(98, 196, 136, 0.18), transparent 36%),
    rgba(18, 28, 24, 0.86);
}
</style>
