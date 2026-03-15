<template>
  <section class="stack">
    <header class="page-header">
      <div>
        <div class="eyebrow">Gestion des utilisateurs</div>
        <h2 class="page-title">Superviser sans masquer les manques API</h2>
        <p class="page-description">
          La liste, le filtrage et la pagination sont operationnels. Les actions CRUD sont
          preparees cote interface, mais le backend n'expose pas encore les endpoints necessaires.
        </p>
      </div>
      <div class="toolbar">
        <button class="button button-primary" disabled>Creer un utilisateur</button>
        <button class="button button-secondary" @click="loadUsers">Actualiser</button>
      </div>
    </header>

    <ResourceNotice
      title="Etat d'integration des utilisateurs"
      message="L'endpoint de liste est disponible, mais les operations detail, creation, modification, activation et suppression restent a implementer cote backend."
      :items="missingUserEndpoints"
    />

    <article class="surface-card section-card">
      <div class="filters-grid">
        <FormField label="Recherche">
          <input v-model="filters.search" placeholder="Nom, email, ville" />
        </FormField>
        <FormField label="Role">
          <BaseSelect v-model="filters.role" :options="roleOptions" />
        </FormField>
        <FormField label="Statut">
          <BaseSelect v-model="filters.status" :options="statusOptions" />
        </FormField>
      </div>
    </article>

    <LoadingState v-if="loading" />
    <ErrorState v-else-if="error" :message="error" retry-label="Recharger" @retry="loadUsers" />
    <EmptyState
      v-else-if="rows.length === 0"
      title="Aucun utilisateur"
      message="Aucun utilisateur ne correspond aux filtres selectionnes."
    />
    <DataTable v-else :columns="columns" :rows="rows" :pagination="pagination" @page-change="changePage">
      <template #cell-fullName="{ row }">
        <div class="identity">
          <strong>{{ row.fullName || "Nom indisponible" }}</strong>
          <span>{{ row.email || "Email indisponible" }}</span>
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
          <button class="button button-ghost" @click="viewUser(row)">Voir</button>
          <button class="button button-secondary" disabled>Modifier</button>
          <button class="button button-warning" disabled>Activer / desactiver</button>
          <button class="button button-danger" disabled>Supprimer</button>
        </div>
      </template>
    </DataTable>

    <article v-if="selectedUser" class="surface-card section-card">
      <div class="page-header compact">
        <div>
          <div class="eyebrow">Fiche lecture seule</div>
          <h3 class="detail-title">{{ selectedUser.fullName }}</h3>
          <p class="page-description">{{ selectedUser.email || "Email indisponible" }}</p>
        </div>
        <button class="button button-secondary" @click="selectedUser = null">Fermer</button>
      </div>
      <div class="two-up">
        <div>
          <strong>Role</strong>
          <p>{{ selectedUser.role }}</p>
        </div>
        <div>
          <strong>Ville</strong>
          <p>{{ selectedUser.city || "-" }}</p>
        </div>
        <div>
          <strong>Code postal</strong>
          <p>{{ selectedUser.postalCode || "-" }}</p>
        </div>
        <div>
          <strong>Statut</strong>
          <p>{{ selectedUser.status }}</p>
        </div>
      </div>
    </article>
  </section>
</template>

<script setup>
import { computed, onMounted, reactive, ref, watch } from "vue";
import BaseSelect from "../components/BaseSelect.vue";
import DataTable from "../components/DataTable.vue";
import EmptyState from "../components/EmptyState.vue";
import ErrorState from "../components/ErrorState.vue";
import FormField from "../components/FormField.vue";
import LoadingState from "../components/LoadingState.vue";
import ResourceNotice from "../components/ResourceNotice.vue";
import StatusBadge from "../components/StatusBadge.vue";
import { adminApi } from "../services/api";
import { useToastStore } from "../store/toastStore";

const { pushToast } = useToastStore();

const loading = ref(true);
const error = ref("");
const rows = ref([]);
const pagination = ref(null);
const selectedUser = ref(null);
const filters = reactive({
  search: "",
  role: "",
  status: "",
  page: 1,
  pageSize: 8
});

const columns = [
  { key: "fullName", label: "Utilisateur" },
  { key: "role", label: "Role" },
  { key: "status", label: "Statut" },
  { key: "city", label: "Ville" }
];

const roleOptions = [
  { label: "Tous les roles", value: "" },
  { label: "Particulier", value: "Particulier" },
  { label: "Prestataire", value: "Prestataire" },
  { label: "Admin", value: "Admin" }
];

const statusOptions = [
  { label: "Tous les statuts", value: "" },
  { label: "Actif", value: "active" },
  { label: "Inactif", value: "inactive" }
];

const missingUserEndpoints = computed(() => [
  "GET /api/admin/users/:id",
  "POST /api/admin/users",
  "PUT /api/admin/users/:id",
  "PATCH /api/admin/users/:id/status",
  "DELETE /api/admin/users/:id"
]);

function roleTone(role) {
  const key = role.toLowerCase();
  if (key.includes("admin")) return "coral";
  if (key.includes("prest")) return "teal";
  return "green";
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

function changePage(page) {
  filters.page = page;
}

function viewUser(row) {
  selectedUser.value = row;
  pushToast({
    title: "Lecture utilisateur",
    message: "La fiche detail s'appuie pour l'instant sur les donnees de liste deja chargees.",
    tone: "amber"
  });
}

watch(
  () => [filters.search, filters.role, filters.status, filters.page],
  () => {
    loadUsers();
  }
);

onMounted(loadUsers);
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

.compact {
  margin-bottom: 18px;
}

.detail-title {
  margin: 10px 0 6px;
  font-family: "Syne", sans-serif;
  font-size: 1.8rem;
}
</style>
