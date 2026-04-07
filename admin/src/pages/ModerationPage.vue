<template>
  <section class="stack">
    <header class="page-header">
      <div>
        <div class="eyebrow">Moderation</div>
        <h2 class="page-title">File de validation</h2>
        <p class="page-description">Pilotage des contenus a publier ou archiver.</p>
      </div>
      <div class="toolbar">
        <button class="button button-secondary" @click="loadQueue">Actualiser</button>
      </div>
    </header>

    <article class="surface-card section-card stack">
      <div class="filters-grid">
        <FormField label="Recherche">
          <input v-model="filters.search" placeholder="Titre, source ou description" />
        </FormField>
        <FormField label="Type">
          <BaseSelect v-model="filters.type" :options="typeOptions" />
        </FormField>
        <FormField label="Statut">
          <BaseSelect v-model="filters.status" :options="statusOptions" />
        </FormField>
      </div>
    </article>

    <LoadingState v-if="loading" />
    <ErrorState v-else-if="error" :message="error" retry-label="Recharger" @retry="loadQueue" />
    <EmptyState
      v-else-if="rows.length === 0"
      title="Aucune validation en attente"
      message="Tous les contenus publies sont deja traites."
    />
    <DataTable v-else :columns="columns" :rows="rows" :pagination="pagination" @page-change="changePage">
      <template #cell-title="{ row }">
        <div class="identity">
          <strong>{{ row.title }}</strong>
          <span>{{ row.description }}</span>
        </div>
      </template>
      <template #cell-type="{ row }">
        <StatusBadge :label="row.type" :tone="row.type === 'event' ? 'teal' : 'green'" />
      </template>
      <template #cell-status="{ row }">
        <StatusBadge :label="row.status" tone="amber" />
      </template>
      <template #actions="{ row }">
        <div class="toolbar actions">
          <button class="button button-primary" @click="publishItem(row)">Publier</button>
          <button class="button button-secondary" @click="archiveItem(row)">Archiver</button>
        </div>
      </template>
    </DataTable>
  </section>
</template>

<script setup>
import { onMounted, reactive, ref, watch } from "vue";
import BaseSelect from "../components/BaseSelect.vue";
import DataTable from "../components/DataTable.vue";
import EmptyState from "../components/EmptyState.vue";
import ErrorState from "../components/ErrorState.vue";
import FormField from "../components/FormField.vue";
import LoadingState from "../components/LoadingState.vue";
import StatusBadge from "../components/StatusBadge.vue";
import { adminApi } from "../services/api";
import { useToastStore } from "../store/toastStore";

const { pushToast } = useToastStore();

const loading = ref(true);
const error = ref("");
const rows = ref([]);
const pagination = ref(null);
const filters = reactive({ search: "", type: "", status: "", page: 1, pageSize: 8 });

const columns = [
  { key: "title", label: "Contenu" },
  { key: "type", label: "Type" },
  { key: "owner", label: "Responsable" },
  { key: "status", label: "Statut" }
];

const typeOptions = [
  { label: "Tous", value: "" },
  { label: "Prestations", value: "prestation" },
  { label: "Evenements", value: "event" }
];

const statusOptions = [
  { label: "Tous", value: "" },
  { label: "Brouillon", value: "draft" },
  { label: "Planifie", value: "planned" }
];

async function loadQueue() {
  loading.value = true;
  error.value = "";
  try {
    const response = await adminApi.listModerationQueue(filters);
    rows.value = response.items;
    pagination.value = response.pagination;
  } catch (err) {
    error.value = err.message ?? "Impossible de charger la moderation.";
  } finally {
    loading.value = false;
  }
}

async function publishItem(row) {
  try {
    await adminApi.publishModerationItem(row.type === "event" ? "events" : "prestations", row.id);
    pushToast({ title: "Contenu publie", message: `${row.title} est maintenant visible.`, tone: "green" });
    await loadQueue();
  } catch (err) {
    pushToast({ title: "Publication impossible", message: err.message ?? "Operation impossible.", tone: "coral" });
  }
}

async function archiveItem(row) {
  try {
    await adminApi.archiveModerationItem(row.type === "event" ? "events" : "prestations", row.id);
    pushToast({ title: "Contenu archive", message: `${row.title} a ete archive.`, tone: "amber" });
    await loadQueue();
  } catch (err) {
    pushToast({ title: "Archivage impossible", message: err.message ?? "Operation impossible.", tone: "coral" });
  }
}

function changePage(page) {
  filters.page = page;
}

watch(() => [filters.search, filters.type, filters.status, filters.page], loadQueue);
onMounted(loadQueue);
</script>

<style scoped>
.identity {
  display: grid;
  gap: 4px;
}

.identity span {
  color: var(--text-secondary);
}

.actions {
  justify-content: flex-end;
}
</style>
