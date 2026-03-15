<template>
  <section class="stack">
    <header class="page-header">
      <div>
        <div class="eyebrow">Gestion des prestations</div>
        <h2 class="page-title">Catalogue admin branche sur les annonces</h2>
        <p class="page-description">
          L'API actuelle permet la lecture et la creation via `/api/annonces`. La mise a jour et
          la suppression sont signalees comme manquantes.
        </p>
      </div>
      <div class="toolbar">
        <button class="button button-primary" @click="showCreate = !showCreate">
          {{ showCreate ? "Masquer le formulaire" : "Nouvelle prestation" }}
        </button>
        <button class="button button-secondary" @click="loadPrestations">Actualiser</button>
      </div>
    </header>

    <ResourceNotice
      title="Mapping backend actuel"
      message="Les prestations sont temporairement mappees sur la ressource annonces du backoffice."
      :items="[
        'GET /api/annonces',
        'POST /api/annonces',
        'PUT /api/annonces/:id a implementer',
        'DELETE /api/annonces/:id a implementer'
      ]"
      tone="green"
      badge="partiel"
    />

    <article v-if="showCreate" class="surface-card section-card stack">
      <div class="filters-grid">
        <FormField label="Titre" :error="formErrors.title">
          <input v-model="form.title" placeholder="Ex: Atelier textile premium" />
        </FormField>
        <FormField label="Type">
          <BaseSelect v-model="form.type" :options="typeOptions" />
        </FormField>
        <FormField label="Prix" :error="formErrors.price">
          <input v-model="form.price" type="number" min="0" step="0.01" placeholder="0" />
        </FormField>
      </div>
      <FormField label="Description" :error="formErrors.description">
        <textarea
          v-model="form.description"
          placeholder="Description de la prestation et valeur apporte"
        ></textarea>
      </FormField>
      <div class="toolbar">
        <button class="button button-primary" @click="submitCreate">Enregistrer</button>
        <button class="button button-secondary" @click="resetForm">Reinitialiser</button>
      </div>
    </article>

    <article class="surface-card section-card">
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

    <LoadingState v-if="loading" />
    <ErrorState
      v-else-if="error"
      :message="error"
      retry-label="Recharger"
      @retry="loadPrestations"
    />
    <EmptyState
      v-else-if="rows.length === 0"
      title="Aucune prestation"
      message="Aucune prestation n'a ete retournee ou ne correspond aux filtres."
    />
    <DataTable v-else :columns="columns" :rows="rows" :pagination="pagination" @page-change="changePage">
      <template #cell-title="{ row }">
        <div class="identity">
          <strong>{{ row.title }}</strong>
          <span>{{ row.description }}</span>
        </div>
      </template>

      <template #cell-price="{ row }">
        {{ formatCurrency(row.price) }}
      </template>

      <template #cell-status="{ row }">
        <StatusBadge :label="row.status" tone="teal" />
      </template>

      <template #actions="{ row }">
        <div class="toolbar actions">
          <button class="button button-secondary" @click="editUnavailable(row)">Modifier</button>
          <button class="button button-danger" @click="confirmDelete(row)">Supprimer</button>
        </div>
      </template>
    </DataTable>

    <ConfirmModal
      :open="Boolean(rowToDelete)"
      title="Suppression impossible pour le moment"
      message="Le backend actuel n'expose pas encore DELETE /api/annonces/:id. L'action reste bloquee pour eviter toute promesse trompeuse."
      confirm-label="Compris"
      @cancel="rowToDelete = null"
      @confirm="rowToDelete = null"
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
import ResourceNotice from "../components/ResourceNotice.vue";
import StatusBadge from "../components/StatusBadge.vue";
import { adminApi } from "../services/api";
import { useToastStore } from "../store/toastStore";
import { formatCurrency } from "../utils/format";

const { pushToast } = useToastStore();

const loading = ref(true);
const error = ref("");
const rows = ref([]);
const pagination = ref(null);
const showCreate = ref(false);
const rowToDelete = ref(null);

const filters = reactive({
  search: "",
  type: "",
  sortBy: "title",
  page: 1,
  pageSize: 8
});

const form = reactive({
  title: "",
  description: "",
  type: "service",
  price: 0
});

const formErrors = reactive({
  title: "",
  description: "",
  price: ""
});

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

const typeFilterOptions = [{ label: "Tous les types", value: "" }, ...typeOptions];

const sortOptions = [
  { label: "Titre", value: "title" },
  { label: "Prix croissant", value: "price-asc" },
  { label: "Prix decroissant", value: "price-desc" }
];

function validateForm() {
  formErrors.title = form.title.trim().length < 3 ? "Le titre doit contenir au moins 3 caracteres." : "";
  formErrors.description =
    form.description.trim().length < 10 ? "La description doit contenir au moins 10 caracteres." : "";
  formErrors.price = Number(form.price) < 0 ? "Le prix doit etre positif." : "";
  return !formErrors.title && !formErrors.description && !formErrors.price;
}

function resetForm() {
  form.title = "";
  form.description = "";
  form.type = "service";
  form.price = 0;
  formErrors.title = "";
  formErrors.description = "";
  formErrors.price = "";
}

async function submitCreate() {
  if (!validateForm()) {
    pushToast({
      title: "Formulaire incomplet",
      message: "Corrige les champs en erreur avant de continuer.",
      tone: "coral"
    });
    return;
  }

  try {
    await adminApi.createPrestation(form);
    pushToast({
      title: "Prestation creee",
      message: "La prestation a ete enregistree via le backoffice.",
      tone: "green"
    });
    resetForm();
    showCreate.value = false;
    await loadPrestations();
  } catch (err) {
    pushToast({
      title: "Creation impossible",
      message: err.message ?? "Impossible de creer la prestation.",
      tone: "coral"
    });
  }
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

function changePage(page) {
  filters.page = page;
}

function editUnavailable(row) {
  pushToast({
    title: "Edition non disponible",
    message: `L'edition de "${row.title}" attend encore l'endpoint PUT /api/annonces/:id.`,
    tone: "amber"
  });
}

function confirmDelete(row) {
  rowToDelete.value = row;
}

watch(
  () => [filters.search, filters.type, filters.sortBy, filters.page],
  () => {
    loadPrestations();
  }
);

onMounted(loadPrestations);
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
