<template>
  <section class="stack">
    <header class="page-header">
      <div>
        <div class="eyebrow">Utilisateurs</div>
        <h2 class="page-title">Gestion des comptes</h2>
        <p class="page-description">Recherche, edition rapide et gestion complete via l'API admin.</p>
      </div>
      <div class="toolbar">
        <button class="button button-primary" @click="openCreate">Nouvel utilisateur</button>
        <button class="button button-secondary" @click="loadUsers">Actualiser</button>
      </div>
    </header>

    <div class="split-grid admin-workspace">
      <article class="surface-card section-card stack admin-panel admin-panel--filters">
        <div class="panel-head panel-head--compact">
          <div>
            <span class="panel-kicker">Comptes</span>
            <h3>Pilotage utilisateurs</h3>
          </div>
          <span class="mini-note">{{ pagination?.total ?? 0 }} comptes</span>
        </div>

        <div class="admin-stats">
          <div class="admin-stat">
            <strong>{{ pagination?.total ?? rows.length }}</strong>
            <span>Total</span>
          </div>
          <div class="admin-stat">
            <strong>{{ activeCount }}</strong>
            <span>Actifs</span>
          </div>
          <div class="admin-stat">
            <strong>{{ adminCount }}</strong>
            <span>Admins</span>
          </div>
        </div>

        <div class="quick-chips">
          <button
            v-for="option in roleOptions"
            :key="option.value || 'all'"
            class="quick-chip"
            :class="{ active: filters.role === option.value }"
            type="button"
            @click="filters.role = option.value"
          >
            {{ option.label }}
          </button>
        </div>

        <div class="filters-compact">
          <FormField label="Recherche">
            <input v-model="filters.search" placeholder="Nom, email, ville" />
          </FormField>
          <FormField label="Statut">
            <BaseSelect v-model="filters.status" :options="statusOptions" />
          </FormField>
        </div>
        <button class="button button-secondary filters-reset" type="button" @click="resetFilters">
          Reinitialiser
        </button>
      </article>

      <article class="surface-card section-card stack admin-panel admin-panel--form">
        <div class="panel-head">
          <div>
            <span class="panel-kicker">{{ editingId ? "Edition" : "Creation" }}</span>
            <h3>{{ editingId ? "Modifier" : "Nouveau compte" }}</h3>
          </div>
          <button class="button button-ghost" @click="resetForm">Vider</button>
        </div>
        <div class="mode-note">
          <strong>{{ editingId ? "Compte existant" : "Compte admin" }}</strong>
          <span>Les champs essentiels sont controles avant enregistrement dans l'API.</span>
        </div>
        <div class="filters-grid form-grid">
          <FormField label="Prenom" :error="formErrors.firstName">
            <input v-model="form.firstName" placeholder="Alice" />
          </FormField>
          <FormField label="Nom" :error="formErrors.lastName">
            <input v-model="form.lastName" placeholder="Martin" />
          </FormField>
          <FormField label="Email" :error="formErrors.email">
            <input v-model="form.email" type="email" placeholder="alice@mail.com" />
          </FormField>
          <FormField label="Ville">
            <input v-model="form.city" placeholder="Paris" />
          </FormField>
          <FormField label="Code postal">
            <input v-model="form.postalCode" placeholder="75011" />
          </FormField>
          <FormField label="Role">
            <BaseSelect v-model="form.role" :options="roleOptions.slice(1)" />
          </FormField>
        </div>
        <div class="toolbar">
          <button class="button button-primary" @click="submitForm">
            {{ editingId ? "Enregistrer" : "Creer" }}
          </button>
          <button v-if="editingId" class="button button-secondary" @click="resetForm">Annuler</button>
        </div>
      </article>
    </div>

    <LoadingState v-if="loading" />
    <ErrorState v-else-if="error" :message="error" retry-label="Recharger" @retry="loadUsers" />
    <EmptyState v-else-if="rows.length === 0" title="Aucun utilisateur" message="Aucun resultat." />
    <DataTable v-else :columns="columns" :rows="rows" :pagination="pagination" @page-change="changePage">
      <template #cell-fullName="{ row }">
        <div class="identity">
          <strong>{{ row.fullName }}</strong>
          <span>{{ row.email }}</span>
        </div>
      </template>
      <template #cell-role="{ row }">
        <StatusBadge :label="row.role" :tone="roleTone(row.role)" />
      </template>
      <template #cell-status="{ row }">
        <StatusBadge :label="row.status" :tone="row.status === 'active' ? 'green' : 'amber'" />
      </template>
      <template #actions="{ row }">
        <div class="toolbar actions">
          <button class="button button-secondary" @click="startEdit(row)">Editer</button>
          <button class="button button-warning" @click="toggleStatus(row)">Statut</button>
          <button class="button button-danger" @click="confirmDelete(row)">Supprimer</button>
        </div>
      </template>
    </DataTable>

    <ConfirmModal
      :open="Boolean(rowToDelete)"
      title="Supprimer ce compte ?"
      :message="rowToDelete ? `${rowToDelete.fullName} sera retire de l'administration.` : ''"
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

const filters = reactive({ search: "", role: "", status: "", page: 1, pageSize: 7 });
const form = reactive({
  firstName: "",
  lastName: "",
  email: "",
  city: "",
  postalCode: "",
  role: "Particulier",
  status: "active"
});
const formErrors = reactive({ firstName: "", lastName: "", email: "" });

const columns = [
  { key: "fullName", label: "Compte" },
  { key: "role", label: "Role" },
  { key: "status", label: "Statut" },
  { key: "city", label: "Ville" }
];

const roleOptions = [
  { label: "Tous", value: "" },
  { label: "Particulier", value: "Particulier" },
  { label: "Prestataire", value: "Prestataire" },
  { label: "Admin", value: "Admin" }
];

const statusOptions = [
  { label: "Tous", value: "" },
  { label: "Actif", value: "active" },
  { label: "Inactif", value: "inactive" }
];

const activeCount = computed(() => rows.value.filter((row) => row.status === "active").length);
const adminCount = computed(() => rows.value.filter((row) => row.role?.toLowerCase().includes("admin")).length);

function roleTone(role) {
  const key = role.toLowerCase();
  if (key.includes("admin")) return "coral";
  if (key.includes("prest")) return "teal";
  return "green";
}

function validateForm() {
  formErrors.firstName = form.firstName.trim().length < 2 ? "Minimum 2 caracteres." : "";
  formErrors.lastName = form.lastName.trim().length < 2 ? "Minimum 2 caracteres." : "";
  formErrors.email = /\S+@\S+\.\S+/.test(form.email) ? "" : "Email invalide.";
  return !formErrors.firstName && !formErrors.lastName && !formErrors.email;
}

function resetForm() {
  editingId.value = "";
  form.firstName = "";
  form.lastName = "";
  form.email = "";
  form.city = "";
  form.postalCode = "";
  form.role = "Particulier";
  form.status = "active";
  formErrors.firstName = "";
  formErrors.lastName = "";
  formErrors.email = "";
}

function openCreate() {
  resetForm();
  window.scrollTo({ top: 0, behavior: "smooth" });
}

function startEdit(row) {
  editingId.value = row.id;
  form.firstName = row.firstName;
  form.lastName = row.lastName;
  form.email = row.email;
  form.city = row.city;
  form.postalCode = row.postalCode;
  form.role = row.role;
  form.status = row.status;
  window.scrollTo({ top: 0, behavior: "smooth" });
}

async function submitForm() {
  if (!validateForm()) return;
  try {
    const payload = { ...form };
    if (editingId.value) {
      await adminApi.updateUser(editingId.value, payload);
      pushToast({ title: "Compte mis a jour", message: "Modification enregistree.", tone: "green" });
    } else {
      await adminApi.createUser(payload);
      pushToast({ title: "Compte cree", message: "Nouvel utilisateur ajoute.", tone: "green" });
    }
    resetForm();
    await loadUsers();
  } catch (err) {
    pushToast({ title: "Echec de l'enregistrement", message: err.message ?? "Operation impossible.", tone: "coral" });
  }
}

async function loadUsers() {
  loading.value = true;
  error.value = "";
  try {
    const response = await adminApi.listUsers(filters);
    rows.value = response.items;
    pagination.value = response.pagination;
  } catch (err) {
    error.value = err.message ?? "Impossible de charger les utilisateurs.";
  } finally {
    loading.value = false;
  }
}

async function toggleStatus(row) {
  try {
    await adminApi.toggleUserStatus(row.id);
    pushToast({ title: "Statut mis a jour", message: `${row.fullName} a change d'etat.`, tone: "amber" });
    await loadUsers();
  } catch (err) {
    pushToast({ title: "Echec de mise a jour", message: err.message ?? "Operation impossible.", tone: "coral" });
  }
}

function confirmDelete(row) {
  rowToDelete.value = row;
}

async function deleteCurrent() {
  if (!rowToDelete.value) return;
  try {
    await adminApi.deleteUser(rowToDelete.value.id);
    pushToast({ title: "Compte supprime", message: "Suppression effectuee.", tone: "coral" });
    rowToDelete.value = null;
    await loadUsers();
  } catch (err) {
    pushToast({ title: "Echec de suppression", message: err.message ?? "Operation impossible.", tone: "coral" });
  }
}

function changePage(page) {
  filters.page = page;
}

function resetFilters() {
  filters.search = "";
  filters.role = "";
  filters.status = "";
  filters.page = 1;
}

watch(() => [filters.search, filters.role, filters.status, filters.page], loadUsers);

onMounted(loadUsers);
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
