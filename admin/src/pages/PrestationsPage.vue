<template>
  <section class="stack">
    <header class="page-header">
      <div>
        <div class="eyebrow">Prestations</div>
        <h2 class="page-title">Catalogue</h2>
        <p class="page-description">Edition rapide avec fallback local si l'API ne repond pas.</p>
      </div>
      <div class="toolbar">
        <button class="button button-primary" @click="showCreate = !showCreate">
          {{ showCreate ? "Fermer" : "Nouvelle prestation" }}
        </button>
        <button class="button button-secondary" @click="loadPrestations">Actualiser</button>
      </div>
    </header>

    <div class="split-grid">
      <article class="surface-card section-card stack">
        <div class="panel-head">
          <h3>Filtres</h3>
          <span class="mini-note">{{ pagination?.total ?? 0 }} lignes</span>
        </div>
        <div class="filters-grid">
          <FormField label="Recherche">
            <input v-model="filters.search" placeholder="Titre ou description" />
          </FormField>
          <FormField label="Type">
            <BaseSelect v-model="filters.type" :options="typeFilterOptions" />
          </FormField>
          <FormField label="Tri">
            <BaseSelect v-model="filters.sortBy" :options="sortOptions" />
          </FormField>
        </div>
      </article>

      <article v-if="showCreate" class="surface-card section-card stack">
        <div class="panel-head">
          <h3>{{ editingId ? "Modifier" : "Ajouter" }}</h3>
          <button class="button button-ghost" @click="resetForm">Vider</button>
        </div>
        <div class="filters-grid">
          <FormField label="Titre" :error="formErrors.title">
            <input v-model="form.title" />
          </FormField>
          <FormField label="Type">
            <BaseSelect v-model="form.type" :options="typeOptions" />
          </FormField>
          <FormField label="Prix" :error="formErrors.price">
            <input v-model="form.price" type="number" min="0" step="0.01" />
          </FormField>
          <FormField label="Statut">
            <BaseSelect v-model="form.status" :options="statusOptions.slice(1)" />
          </FormField>
        </div>
        <FormField label="Description" :error="formErrors.description">
          <textarea v-model="form.description"></textarea>
        </FormField>
        <div class="toolbar">
          <button class="button button-primary" @click="submitForm">
            {{ editingId ? "Enregistrer" : "Creer" }}
          </button>
        </div>
      </article>
    </div>

    <LoadingState v-if="loading" />
    <ErrorState v-else-if="error" :message="error" retry-label="Recharger" @retry="loadPrestations" />
    <EmptyState v-else-if="rows.length === 0" title="Aucune prestation" message="Aucun resultat." />
    <DataTable v-else :columns="columns" :rows="rows" :pagination="pagination" @page-change="changePage">
      <template #cell-title="{ row }">
        <div class="identity">
          <strong>{{ row.title }}</strong>
          <span>{{ row.description }}</span>
        </div>
      </template>
      <template #cell-price="{ row }">{{ formatCurrency(row.price) }}</template>
      <template #cell-status="{ row }"><StatusBadge :label="row.status" tone="teal" /></template>
      <template #actions="{ row }">
        <div class="toolbar actions">
          <button class="button button-secondary" @click="startEdit(row)">Editer</button>
          <button class="button button-danger" @click="confirmDelete(row)">Supprimer</button>
        </div>
      </template>
    </DataTable>

    <ConfirmModal
      :open="Boolean(rowToDelete)"
      title="Supprimer cette prestation ?"
      :message="rowToDelete ? rowToDelete.title : ''"
      confirm-label="Supprimer"
      @cancel="rowToDelete = null"
      @confirm="deleteCurrent"
    />
  </section>
</template>

<script setup>
import { onMounted, reactive, ref, watch } from "vue";
import BaseSelect from "../components/BaseSelect.vue";
import ConfirmModal from "../components/ConfirmModal.vue";
import DataTable from "../components/DataTable.vue";
import EmptyState from "../components/EmptyState.vue";
import ErrorState from "../components/ErrorState.vue";
import FormField from "../components/FormField.vue";
import LoadingState from "../components/LoadingState.vue";
import StatusBadge from "../components/StatusBadge.vue";
import { adminApi } from "../services/api";
import { useToastStore } from "../store/toastStore";
import { formatCurrency } from "../utils/format";

const { pushToast } = useToastStore();

const loading = ref(true);
const error = ref("");
const rows = ref([]);
const pagination = ref(null);
const showCreate = ref(true);
const editingId = ref("");
const rowToDelete = ref(null);

const filters = reactive({ search: "", type: "", sortBy: "title", page: 1, pageSize: 7 });
const form = reactive({ title: "", description: "", type: "service", price: 0, status: "draft" });
const formErrors = reactive({ title: "", description: "", price: "" });

const columns = [
  { key: "title", label: "Prestation" },
  { key: "type", label: "Type" },
  { key: "price", label: "Prix", align: "right" },
  { key: "status", label: "Statut" }
];

const typeOptions = [
  { label: "Service", value: "service" },
  { label: "Vente", value: "vente" },
  { label: "Don", value: "don" }
];
const typeFilterOptions = [{ label: "Tous", value: "" }, ...typeOptions];
const sortOptions = [
  { label: "Titre", value: "title" },
  { label: "Prix croissant", value: "price-asc" },
  { label: "Prix decroissant", value: "price-desc" }
];
const statusOptions = [
  { label: "Tous", value: "" },
  { label: "Brouillon", value: "draft" },
  { label: "Publie", value: "published" }
];

function validateForm() {
  formErrors.title = form.title.trim().length < 3 ? "Titre trop court." : "";
  formErrors.description = form.description.trim().length < 10 ? "Description trop courte." : "";
  formErrors.price = Number(form.price) < 0 ? "Prix invalide." : "";
  return !formErrors.title && !formErrors.description && !formErrors.price;
}

function resetForm() {
  editingId.value = "";
  form.title = "";
  form.description = "";
  form.type = "service";
  form.price = 0;
  form.status = "draft";
  formErrors.title = "";
  formErrors.description = "";
  formErrors.price = "";
}

function startEdit(row) {
  showCreate.value = true;
  editingId.value = row.id;
  form.title = row.title;
  form.description = row.description;
  form.type = row.type;
  form.price = row.price;
  form.status = row.status;
  window.scrollTo({ top: 0, behavior: "smooth" });
}

async function submitForm() {
  if (!validateForm()) return;
  if (editingId.value) {
    await adminApi.updatePrestation(editingId.value, { ...form });
    pushToast({ title: "Prestation mise a jour", message: "Modification enregistree.", tone: "green" });
  } else {
    await adminApi.createPrestation({ ...form });
    pushToast({ title: "Prestation creee", message: "Nouvelle ligne ajoutee.", tone: "green" });
  }
  resetForm();
  await loadPrestations();
}

async function loadPrestations() {
  loading.value = true;
  error.value = "";
  try {
    const response = await adminApi.listPrestations(filters);
    rows.value = response.items;
    pagination.value = response.pagination;
  } catch (err) {
    error.value = err.message ?? "Impossible de charger les prestations.";
  } finally {
    loading.value = false;
  }
}

function confirmDelete(row) {
  rowToDelete.value = row;
}

async function deleteCurrent() {
  if (!rowToDelete.value) return;
  await adminApi.deletePrestation(rowToDelete.value.id);
  pushToast({ title: "Prestation supprimee", message: "Suppression locale effectuee.", tone: "coral" });
  rowToDelete.value = null;
  await loadPrestations();
}

function changePage(page) {
  filters.page = page;
}

watch(() => [filters.search, filters.type, filters.sortBy, filters.page], loadPrestations);

onMounted(loadPrestations);
</script>

<style scoped>
.identity {
  display: grid;
  gap: 4px;
}

.identity span,
.mini-note {
  color: var(--text-secondary);
}

.panel-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.panel-head h3 {
  margin: 0;
  font-family: "Syne", sans-serif;
}

.actions {
  justify-content: flex-end;
}
</style>
