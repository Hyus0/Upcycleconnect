<template>
  <section class="stack">
    <header class="page-header">
      <div>
        <div class="eyebrow">Categories</div>
        <h2 class="page-title">Taxonomie</h2>
        <p class="page-description">CRUD admin complet pour piloter la taxonomie du catalogue.</p>
      </div>
      <div class="toolbar">
        <button class="button button-secondary" @click="loadCategories">Actualiser</button>
      </div>
    </header>

    <div class="split-grid admin-workspace">
      <article class="surface-card section-card stack admin-panel admin-panel--form">
        <div class="panel-head">
          <div>
            <span class="panel-kicker">{{ editingId ? "Edition" : "Creation" }}</span>
            <h3>{{ editingId ? "Modifier" : "Nouvelle categorie" }}</h3>
          </div>
          <button class="button button-ghost" @click="resetForm">Vider</button>
        </div>
        <div class="mode-note">
          <strong>Taxonomie catalogue</strong>
          <span>Organise les annonces, prestations et futures recherches publiques.</span>
        </div>
        <FormField label="Nom" :error="formError">
          <input v-model="form.name" />
        </FormField>
        <FormField label="Categorie parente">
          <BaseSelect v-model="form.parentId" :options="parentOptions" />
        </FormField>
        <FormField label="Description">
          <textarea v-model="form.description"></textarea>
        </FormField>
        <FormField label="Statut">
          <BaseSelect v-model="form.status" :options="statusOptions" />
        </FormField>
        <div class="toolbar">
          <button class="button button-primary" @click="submitForm">
            {{ editingId ? "Enregistrer" : "Creer" }}
          </button>
        </div>
      </article>

      <article class="surface-card section-card stack admin-panel admin-panel--filters">
        <div class="panel-head panel-head--compact">
          <div>
            <span class="panel-kicker">Catalogue</span>
            <h3>Recherche</h3>
          </div>
          <span class="mini-note">{{ pagination?.total ?? 0 }} categories</span>
        </div>
        <div class="admin-stats">
          <div class="admin-stat">
            <strong>{{ pagination?.total ?? rows.length }}</strong>
            <span>Total</span>
          </div>
          <div class="admin-stat">
            <strong>{{ activeCount }}</strong>
            <span>Actives</span>
          </div>
          <div class="admin-stat">
            <strong>{{ rootCount }}</strong>
            <span>Racines</span>
          </div>
        </div>
        <FormField label="Filtre">
          <input v-model="filters.search" placeholder="Nom ou description" />
        </FormField>
        <button class="button button-secondary filters-reset" type="button" @click="resetFilters">
          Reinitialiser
        </button>
      </article>
    </div>

    <LoadingState v-if="loading" />
    <ErrorState v-else-if="error" :message="error" retry-label="Recharger" @retry="loadCategories" />
    <EmptyState v-else-if="rows.length === 0" title="Aucune categorie" message="Aucun resultat." />
    <DataTable v-else :columns="columns" :rows="rows" :pagination="pagination" @page-change="changePage">
      <template #cell-name="{ row }">
        <div class="identity">
          <strong>{{ row.name }}</strong>
          <span>{{ row.description }}</span>
        </div>
      </template>
      <template #cell-parentId="{ row }">
        <span>{{ resolveParent(row.parentId) }}</span>
      </template>
      <template #cell-status="{ row }">
        <StatusBadge :label="row.status" tone="green" />
      </template>
      <template #actions="{ row }">
        <div class="toolbar actions">
          <button class="button button-secondary" @click="startEdit(row)">Editer</button>
          <button class="button button-danger" @click="confirmDelete(row)">Supprimer</button>
        </div>
      </template>
    </DataTable>

    <ConfirmModal
      :open="Boolean(rowToDelete)"
      title="Supprimer cette categorie ?"
      :message="rowToDelete ? rowToDelete.name : ''"
      confirm-label="Supprimer"
      @cancel="rowToDelete = null"
      @confirm="deleteCurrent"
    />
  </section>
</template>

<script setup>
import { computed, onMounted, reactive, ref, watch } from "vue";
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

const { pushToast } = useToastStore();

const loading = ref(true);
const error = ref("");
const rows = ref([]);
const pagination = ref(null);
const editingId = ref("");
const rowToDelete = ref(null);
const filters = reactive({ search: "", page: 1, pageSize: 7 });
const form = reactive({ name: "", parentId: "", description: "", status: "active" });
const formError = ref("");

const columns = [
  { key: "name", label: "Categorie" },
  { key: "parentId", label: "Parente" },
  { key: "status", label: "Statut" }
];

const statusOptions = [
  { label: "Actif", value: "active" },
  { label: "Archive", value: "archived" }
];

const parentOptions = computed(() => [
  { label: "Aucune", value: "" },
  ...rows.value.map((item) => ({ label: item.name, value: item.id }))
]);
const activeCount = computed(() => rows.value.filter((item) => item.status === "active").length);
const rootCount = computed(() => rows.value.filter((item) => !item.parentId).length);

function resolveParent(parentId) {
  if (!parentId) return "Racine";
  return rows.value.find((item) => item.id === parentId)?.name ?? "Inconnue";
}

function resetForm() {
  editingId.value = "";
  form.name = "";
  form.parentId = "";
  form.description = "";
  form.status = "active";
  formError.value = "";
}

function startEdit(row) {
  editingId.value = row.id;
  form.name = row.name;
  form.parentId = row.parentId;
  form.description = row.description;
  form.status = row.status;
}

async function submitForm() {
  formError.value = form.name.trim().length < 2 ? "Le nom doit contenir au moins 2 caracteres." : "";
  if (formError.value) {
    pushToast({ title: "Categorie invalide", message: formError.value, tone: "coral" });
    return;
  }

  try {
    if (editingId.value) {
      await adminApi.updateCategory(editingId.value, { ...form });
      pushToast({ title: "Categorie mise a jour", message: "Modification enregistree.", tone: "green" });
    } else {
      await adminApi.createCategory({ ...form });
      pushToast({ title: "Categorie creee", message: "Nouvelle categorie ajoutee.", tone: "green" });
    }
    resetForm();
    await loadCategories();
  } catch (err) {
    pushToast({ title: "Echec de l'enregistrement", message: err.message ?? "Operation impossible.", tone: "coral" });
  }
}

async function loadCategories() {
  loading.value = true;
  error.value = "";
  try {
    const response = await adminApi.listCategories(filters);
    rows.value = response.items;
    pagination.value = response.pagination;
  } catch (err) {
    error.value = err.message ?? "Impossible de charger les categories.";
  } finally {
    loading.value = false;
  }
}

function confirmDelete(row) {
  rowToDelete.value = row;
}

async function deleteCurrent() {
  if (!rowToDelete.value) return;
  try {
    await adminApi.deleteCategory(rowToDelete.value.id);
    pushToast({ title: "Categorie supprimee", message: "Suppression effectuee.", tone: "coral" });
    rowToDelete.value = null;
    await loadCategories();
  } catch (err) {
    pushToast({ title: "Echec de suppression", message: err.message ?? "Operation impossible.", tone: "coral" });
  }
}

function changePage(page) {
  filters.page = page;
}

function resetFilters() {
  filters.search = "";
  filters.page = 1;
}

watch(() => [filters.search, filters.page], loadCategories);
onMounted(loadCategories);
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
}

.panel-head h3 {
  margin: 0;
  font-family: "Syne", sans-serif;
}

.actions {
  justify-content: flex-end;
}
</style>
