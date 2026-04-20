<template>
  <section class="stack">
    <header class="page-header">
      <div>
        <div class="eyebrow">Evenements</div>
        <h2 class="page-title">Planning</h2>
        <p class="page-description">Gestion admin des sessions, ateliers et collectes via l'API.</p>
      </div>
      <div class="toolbar">
        <button class="button button-secondary" @click="loadEvents">Actualiser</button>
      </div>
    </header>

    <div class="split-grid admin-workspace">
      <article class="surface-card section-card stack admin-panel admin-panel--form">
        <div class="panel-head">
          <div>
            <span class="panel-kicker">{{ editingId ? "Edition" : "Creation" }}</span>
            <h3>{{ editingId ? "Modifier" : "Nouvel evenement" }}</h3>
          </div>
          <button class="button button-ghost" @click="resetForm">Vider</button>
        </div>
        <div class="mode-note">
          <strong>Regle metier</strong>
          <span>Un evenement planifie ou publie ne peut pas etre date dans le passe.</span>
        </div>
        <div class="filters-grid form-grid">
          <FormField label="Titre" :error="formErrors.title"><input v-model="form.title" /></FormField>
          <FormField label="Lieu" :error="formErrors.location"><input v-model="form.location" /></FormField>
          <FormField label="Date" :error="formErrors.date" :hint="dateHint"><DatePickerField v-model="form.date" :min="eventDateMin" /></FormField>
          <FormField label="Capacite" :error="formErrors.capacity"><input v-model="form.capacity" type="number" min="0" /></FormField>
          <FormField label="Statut"><BaseSelect v-model="form.status" :options="statusOptions.slice(1)" /></FormField>
        </div>
        <FormField label="Description"><textarea v-model="form.description"></textarea></FormField>
        <div class="toolbar">
          <button class="button button-primary" @click="submitForm">
            {{ editingId ? "Enregistrer" : "Creer" }}
          </button>
        </div>
      </article>

      <article class="surface-card section-card stack admin-panel admin-panel--filters">
        <div class="panel-head panel-head--compact">
          <div>
            <span class="panel-kicker">Planning</span>
            <h3>Evenements actifs</h3>
          </div>
          <span class="mini-note">{{ pagination?.total ?? 0 }} evenements</span>
        </div>
        <div class="admin-stats">
          <div class="admin-stat">
            <strong>{{ pagination?.total ?? rows.length }}</strong>
            <span>Total</span>
          </div>
          <div class="admin-stat">
            <strong>{{ publishedCount }}</strong>
            <span>Publies</span>
          </div>
          <div class="admin-stat">
            <strong>{{ plannedCount }}</strong>
            <span>Planifies</span>
          </div>
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
        <div class="filters-compact">
          <FormField label="Recherche"><input v-model="filters.search" placeholder="Titre ou lieu" /></FormField>
          <FormField label="Date"><DatePickerField v-model="filters.date" /></FormField>
        </div>
        <button class="button button-secondary filters-reset" type="button" @click="resetFilters">
          Reinitialiser
        </button>
      </article>
    </div>

    <LoadingState v-if="loading" />
    <ErrorState v-else-if="error" :message="error" retry-label="Recharger" @retry="loadEvents" />
    <EmptyState v-else-if="rows.length === 0" title="Aucun evenement" message="Aucun resultat." />
    <DataTable v-else :columns="columns" :rows="rows" :pagination="pagination" @page-change="changePage">
      <template #cell-title="{ row }">
        <div class="identity">
          <strong>{{ row.title }}</strong>
          <span>{{ row.location }}</span>
        </div>
      </template>
      <template #cell-date="{ row }"><strong>{{ row.date }}</strong></template>
      <template #cell-status="{ row }"><StatusBadge :label="row.status" tone="amber" /></template>
      <template #actions="{ row }">
        <div class="toolbar actions">
          <button class="button button-secondary" @click="startEdit(row)">Editer</button>
          <button class="button button-danger" @click="confirmDelete(row)">Supprimer</button>
        </div>
      </template>
    </DataTable>

    <ConfirmModal
      :open="Boolean(rowToDelete)"
      title="Supprimer cet evenement ?"
      :message="rowToDelete ? rowToDelete.title : ''"
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
import DatePickerField from "../components/DatePickerField.vue";
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
const filters = reactive({ search: "", date: "", status: "", page: 1, pageSize: 7 });
const form = reactive({ title: "", location: "", date: "", status: "planned", capacity: 0, description: "" });
const formErrors = reactive({ title: "", location: "", date: "", capacity: "" });

const statusOptions = [
  { label: "Tous", value: "" },
  { label: "Planifie", value: "planned" },
  { label: "Publie", value: "published" },
  { label: "Archive", value: "archived" }
];

const columns = [
  { key: "title", label: "Evenement" },
  { key: "date", label: "Date" },
  { key: "status", label: "Statut" },
  { key: "capacity", label: "Places", align: "right" }
];

const today = new Date().toISOString().slice(0, 10);

const eventDateMin = computed(() => (form.status === "archived" ? "" : today));
const dateHint = computed(() =>
  form.status === "archived"
    ? "Un evenement archive peut conserver une date passee."
    : "Les evenements planifies ou publies doivent etre dates a partir d'aujourd'hui."
);
const publishedCount = computed(() => rows.value.filter((row) => row.status === "published").length);
const plannedCount = computed(() => rows.value.filter((row) => row.status === "planned").length);

function resetForm() {
  editingId.value = "";
  form.title = "";
  form.location = "";
  form.date = "";
  form.status = "planned";
  form.capacity = 0;
  form.description = "";
  formErrors.title = "";
  formErrors.location = "";
  formErrors.date = "";
  formErrors.capacity = "";
}

function startEdit(row) {
  editingId.value = row.id;
  form.title = row.title;
  form.location = row.location;
  form.date = row.date;
  form.status = row.status;
  form.capacity = row.capacity;
  form.description = row.description;
}

async function submitForm() {
  formErrors.title = form.title.trim().length < 3 ? "Titre trop court." : "";
  formErrors.location =
    (form.status === "planned" || form.status === "published") && form.location.trim().length < 2
      ? "Le lieu est obligatoire pour un evenement visible."
      : "";
  formErrors.date = form.date ? "" : "La date est obligatoire.";
  formErrors.capacity = Number(form.capacity) < 0 ? "Capacite invalide." : "";
  if (form.status === "published" && Number(form.capacity) <= 0) {
    formErrors.capacity = "Un evenement publie doit avoir une capacite superieure a 0.";
  }
  if (form.date && form.status !== "archived" && form.date < today) {
    formErrors.date = "Un evenement planifie ou publie ne peut pas etre dans le passe.";
  }
  if (formErrors.title || formErrors.location || formErrors.date || formErrors.capacity) {
    pushToast({ title: "Evenement invalide", message: "Corrige les champs avant enregistrement.", tone: "coral" });
    return;
  }

  const payload = { ...form, capacity: Number(form.capacity) };
  try {
    if (editingId.value) {
      await adminApi.updateEvent(editingId.value, payload);
      pushToast({ title: "Evenement mis a jour", message: "Modification enregistree.", tone: "green" });
    } else {
      await adminApi.createEvent(payload);
      pushToast({ title: "Evenement cree", message: "Nouvelle entree ajoutee.", tone: "green" });
    }
    resetForm();
    await loadEvents();
  } catch (err) {
    pushToast({ title: "Echec de l'enregistrement", message: err.message ?? "Operation impossible.", tone: "coral" });
  }
}

async function loadEvents() {
  loading.value = true;
  error.value = "";
  try {
    const response = await adminApi.listEvents(filters);
    rows.value = response.items;
    pagination.value = response.pagination;
  } catch (err) {
    error.value = err.message ?? "Impossible de charger les evenements.";
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
    await adminApi.deleteEvent(rowToDelete.value.id);
    pushToast({ title: "Evenement supprime", message: "Suppression effectuee.", tone: "coral" });
    rowToDelete.value = null;
    await loadEvents();
  } catch (err) {
    pushToast({ title: "Echec de suppression", message: err.message ?? "Operation impossible.", tone: "coral" });
  }
}

function changePage(page) {
  filters.page = page;
}

function resetFilters() {
  filters.search = "";
  filters.date = "";
  filters.status = "";
  filters.page = 1;
}

watch(() => [filters.search, filters.date, filters.status, filters.page], loadEvents);
onMounted(loadEvents);
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
